package main

import (
	"github.com/goodsign/monday"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Page struct {
	Title string
	Body  []byte
}

const version = "2019.4.2.27"
const deleteLogsAfter = 240 * time.Hour

func main() {
	LogDirectoryFileCheck("MAIN")
	CreateConfigIfNotExists()
	LoadSettingsFromConfigFile()
	router := httprouter.New()
	timeStreamer := sse.New()
	workplaces := sse.New()
	overview := sse.New()

	router.GET("/lcd_rompa", LcdRompa)
	router.GET("/css/darcula.css", darcula)
	router.GET("/js/metro.min.js", metrojs)
	router.GET("/css/metro-all.css", metrocss)

	router.Handler("GET", "/time", timeStreamer)
	router.Handler("GET", "/workplaces", workplaces)
	router.Handler("GET", "/overview", overview)
	go StreamTime(timeStreamer)
	go StreamWorkplaceData(workplaces)
	go StreamOverview(overview)
	LogInfo("MAIN", "Server running")
	_ = http.ListenAndServe(":80", router)
}

func StreamOverview(streamer *sse.Streamer) {
	var workplaces []Workplace
	for {
		LogInfo("MAIN", "Streaming overview")
		workplaces = nil
		connectionString, dialect := CheckDatabaseType()
		db, err := gorm.Open(dialect, connectionString)
		defer db.Close()
		if err != nil {
			LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
			return
		}
		db.Where("WorkplaceDivisionID = ?", 1).Find(&workplaces)
		production := 0
		downtime := 0
		offline := 0
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
		}
		sum := production + offline + downtime
		if sum == 0 {
			streamer.SendString("", "overview", "Produkce 0%;Prostoj 0%;Vypnuto 0%;Porucha 0%")
			time.Sleep(10 * time.Second)
			continue
		}
		productionPercent := production / sum
		downtimePercent := downtime / sum
		offlinePercent := sum - productionPercent - downtimePercent
		breakdownPercent := 0
		streamer.SendString("", "overview", "Produkce "+strconv.Itoa(productionPercent*100)+"%;Prostoj "+strconv.Itoa(downtimePercent*100)+"%;Vypnuto "+strconv.Itoa(offlinePercent*100)+"%;Porucha "+strconv.Itoa(breakdownPercent*100)+"%")
		time.Sleep(10 * time.Second)
	}
}

func StreamWorkplaceData(streamer *sse.Streamer) {
	var workplaces []Workplace
	for {
		workplaces = nil
		connectionString, dialect := CheckDatabaseType()
		db, err := gorm.Open(dialect, connectionString)

		if err != nil {
			LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
			return
		}
		defer db.Close()
		db.Where("WorkplaceDivisionID = ?", 1).Find(&workplaces)

		for _, workplace := range workplaces {
			color := "red"
			workplaceState := WorkplaceState{}
			db.Where("WorkplaceID = ?", workplace.OID).Where("DTE is null").Find(&workplaceState)
			switch workplaceState.StateID {
			case 1:
				color = "green"
			case 2:
				color = "orange"
			}
			duration := time.Now().Sub(workplaceState.DTS)
			streamer.SendString("", "workplaces", workplace.Name+";"+workplace.Name+"<br>User<br>InforData <span class=\"badge-bottom\">"+duration.String()+"</span>;"+color)
		}
		time.Sleep(10 * time.Second)
	}
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("html/" + tmpl + ".html")
	_ = t.Execute(w, p)
}
