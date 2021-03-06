package main

import (
	"github.com/davidscholberg/go-durationfmt"
	"github.com/goodsign/monday"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const version = "2020.4.3.2"
const deleteLogsAfter = 240 * time.Hour

func main() {
	LogDirectoryFileCheck("MAIN")
	CreateConfigIfNotExists()
	LoadSettingsFromConfigFile()
	router := httprouter.New()
	time := sse.New()
	workplaces := sse.New()
	overview := sse.New()

	router.GET("/lcd_rompa", LcdRompa)
	router.GET("/mobile_rompa", MobileRompa)
	router.GET("/css/darcula.css", darcula)
	router.GET("/js/metro.min.js", metrojs)
	router.GET("/css/metro-all.css", metrocss)

	router.Handler("GET", "/time", time)
	router.Handler("GET", "/workplaces", workplaces)
	router.Handler("GET", "/overview", overview)
	go StreamTime(time)
	go StreamWorkplaces(workplaces)
	go StreamOverview(overview)
	LogInfo("MAIN", "Server running")
	RestartAllLCDs()
	_ = http.ListenAndServe(":80", router)
}

func RestartAllLCDs() {
	var televisions []Lcd
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	db.Find(&televisions)

	for _, television := range televisions {
		LogInfo("MAIN", "Restarting television ["+television.Name+"] with ip address ["+television.IPAddress+"]")
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		response, err := client.Get("http://" + television.IPAddress + "/restart.php")
		if err != nil {
			LogWarning("MAIN", "Television is offline")
			continue
		} else {
			response.Body.Close()
			LogInfo("MAIN", "Television ["+television.Name+"] with ip address ["+television.IPAddress+"] restarted")
		}
	}
}

func StreamOverview(streamer *sse.Streamer) {
	var workplaces []Workplace
	for {
		LogInfo("MAIN", "Streaming overview")
		workplaces = nil
		connectionString, dialect := CheckDatabaseType()
		db, err := gorm.Open(dialect, connectionString)
		if err != nil {
			LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		db.Where("WorkplaceDivisionID = ?", 1).Find(&workplaces)
		production := 0
		downtime := 0
		offline := 0
		repair := 0
		for _, workplace := range workplaces {
			workplaceState := WorkplaceState{}
			db.Where("WorkplaceID = ?", workplace.OID).Where("DTE is null").Find(&workplaceState)
			switch workplaceState.StateID {
			case 1:
				production++
			case 2:
				downtime++
			case 3:
				offline++
			}
			terminalInputIdle := TerminalInputIdle{}
			db.Where("DeviceID = ?", workplace.DeviceID).Where("DTE is null").Where("IdleId=136").Find(&terminalInputIdle)
			if terminalInputIdle.OID > 0 {
				repair++
			}
		}
		downtime = downtime - repair
		sum := production + offline + downtime + repair
		if sum == 0 {
			streamer.SendString("", "overview", "Produkce 0%;Prostoj 0%;Vypnuto 0%;Porucha 0%")
			time.Sleep(10 * time.Second)
			continue
		}
		LogInfo("MAIN", "Production: "+strconv.Itoa(production)+", Downtime: "+strconv.Itoa(downtime)+", Offline: "+strconv.Itoa(offline))
		productionPercent := production * 100 / sum
		downtimePercent := downtime * 100 / sum
		repairPercent := repair * 100 / sum
		offlinePercent := 100 - productionPercent - downtimePercent - repair
		floatPointMiscalculation := offline == 0 && offlinePercent > 0
		if floatPointMiscalculation {
			offlinePercent = 0
			if downtimePercent > productionPercent {
				downtimePercent++
			} else {
				productionPercent++
			}
		}
		LogInfo("MAIN", "Production: "+strconv.Itoa(productionPercent)+", Downtime: "+strconv.Itoa(downtimePercent)+", Offline: "+strconv.Itoa(offlinePercent))
		streamer.SendString("", "overview", "Produkce "+strconv.Itoa(productionPercent)+"%;Prostoj "+strconv.Itoa(downtimePercent)+"%;Vypnuto "+strconv.Itoa(offlinePercent)+"%;Porucha "+strconv.Itoa(repairPercent*100)+"%")
		db.Close()
		time.Sleep(10 * time.Second)
	}
}

func StreamWorkplaces(streamer *sse.Streamer) {
	var workplaces []Workplace
	for {
		LogInfo("MAIN", "Streaming workplaces")
		workplaces = nil
		connectionString, dialect := CheckDatabaseType()
		db, err := gorm.Open(dialect, connectionString)
		if err != nil {
			LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		db.Where("WorkplaceDivisionID = ?", 1).Find(&workplaces)
		LogInfo("MAIN", "Workplaces count: "+strconv.Itoa(len(workplaces)))
		for _, workplace := range workplaces {
			workplaceState := WorkplaceState{}
			db.Where("WorkplaceID = ?", workplace.OID).Where("DTE is null").Find(&workplaceState)
			terminalInputOrder := TerminalInputOrder{}
			db.Where("DeviceID = ?", workplace.DeviceID).Where("DTE is null").Find(&terminalInputOrder)
			workplaceHasOpenOrder := terminalInputOrder.OID > 0
			var tools string
			var products string
			var userName string
			if workplaceHasOpenOrder {
				user := User{}
				db.Where("OID = ?", terminalInputOrder.UserID).Find(&user)
				userName = user.Name
				order := Order{}
				db.Where("OID = ?", terminalInputOrder.OrderID).Find(&order)
				tools, products = GetToolsAndProductsForWorkplace(order, workplace)
				LogInfo(workplace.Name, "Data for order downloaded")
			} else {
				LogInfo(workplace.Name, "No open order, not downloading Infor data")
			}
			terminalInputIdle := TerminalInputIdle{}
			db.Where("DeviceID = ?", workplace.DeviceID).Where("DTE is null").Where("IdleId=136").Find(&terminalInputIdle)
			color := GetColorForWorkplace(workplaceState, terminalInputIdle.OID)
			LogInfo(workplace.Name, "Workplace color: "+color)
			duration, err := durationfmt.Format(time.Now().Sub(workplaceState.DTS), "%dd %hh %mm")
			if err != nil {
				LogError(workplace.Name, "Problem parsing datetime: "+err.Error())
			}
			LogInfo(workplace.Name, "Streaming data to LCD")
			streamer.SendString("", "workplaces", workplace.Name+";"+workplace.Name+"<br>"+userName+"<br>"+tools+"<br>"+products+"<span class=\"badge-bottom\">"+duration+"</span>;"+color)
			LogInfo(workplace.Name, "Data streamed")
		}
		LogInfo("MAIN", "Workplaces streamed, waiting 10 second for another run")
		db.Close()
		time.Sleep(10 * time.Second)
	}
}

func GetToolsAndProductsForWorkplace(order Order, workplace Workplace) (string, string) {
	switch order.Name {
	case "Internal":
		{
			LogInfo(workplace.Name, "Internal order opened, not downloading Infor data")
			return "", ""
		}
	default:
		{
			LogInfo(workplace.Name, "Getting Infor data for open order: "+order.Name)
			tools, products := GetInforData(order)
			return tools, products
		}
	}
}

func GetColorForWorkplace(workplaceState WorkplaceState, terminalInputIdleId int) string {
	repairTerminalIdleOpened := terminalInputIdleId > 0
	if repairTerminalIdleOpened {
		return "orange"
	}
	switch workplaceState.StateID {
	case 1:
		return "green"
	case 2:
		return "red"
	default:
		return "yellow"
	}
}

func GetInforData(order Order) (string, string) {
	db, err := gorm.Open("mssql", "sqlserver://DataReader:RompaCZ@10.60.1.5/ERPLN105?database=rompaln")
	defer db.Close()
	if err != nil {
		println("Error opening db: " + err.Error())
		return "", ""
	}
	var tools string
	var items string
	command := "select	stuff((select ', '+ltrim((t_tool)) from ttirpt2401000 T where T.t_item=PS.t_mitm and T.t_prmd=(PS.t_prmd) and T.t_prmv=(PS.t_pmrv) FOR XML path('')),1,1,'') Tools,		stuff((select ', '+ltrim((t_pitm)) from ttirpt2301000 Pr where Pr.t_prmd=(PS.t_prmd) and Pr.t_prmv=(PS.t_pmrv) FOR XML path('')),1,1,'') Items from ttirpt4011000 PS		where PS.t_prsh='" + order.Name + "'"
	row := db.Raw(command).Row()
	err = row.Scan(&tools, &items)
	if err != nil {
		LogError(order.Name, err.Error())
	}
	return strings.Replace(tools, " ", "", 1), strings.Replace(items, " ", "", 1)
}

func StreamTime(streamer *sse.Streamer) {
	for {
		streamer.SendString("", "time", monday.Format(time.Now(), "Monday, 2. January 2006 15:04:05", monday.LocaleCsCZ))
		time.Sleep(1 * time.Second)
	}
}

func darcula(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "css/darcula.css")
}

func metrojs(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "js/metro.min.js")
}

func metrocss(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "css/metro-all.css")
}
