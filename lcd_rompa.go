package main

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

type LcdWorkplaces struct {
	LcdWorkplaces []LcdWorkplace
	Version       string
}
type LcdWorkplace struct {
	Name       string
	User       string
	StateColor string
	Duration   time.Duration
	InforData  string
}

func LcdRompa(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Displaying LCD for Rompa")
	tmpl := template.Must(template.ParseFiles("html/lcd_rompa.html"))
	var workplaces []Workplace
	lcdWorkplaces := LcdWorkplaces{}

	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("WorkplaceDivisionID = ?", 1).Find(&workplaces)

	for _, workplace := range workplaces {
		LogInfo("MAIN", "Adding worklace: " + workplace.Name)
		lcdWorkplace := LcdWorkplace{Name: workplace.Name, User: "loading", StateColor: "", Duration: time.Duration(0 * time.Hour), InforData: "loading..."}
		lcdWorkplaces.LcdWorkplaces = append(lcdWorkplaces.LcdWorkplaces, lcdWorkplace)
	}
	lcdWorkplaces.Version = "version: " + version
	_ = tmpl.Execute(writer, lcdWorkplaces)
}