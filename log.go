package main

import (
	"github.com/gookit/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LogInfo(reference, data string) {
	color.Green.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data)
	AppendDataToLog("INF", reference, data)
}

func LogError(reference, data string) {
	color.Red.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --ERR-- " + data)
	AppendDataToLog("ERR", reference, data)
	AppendDataToErrLog("ERR", reference, data)
}

func LogWarning(reference, data string) {
	color.Yellow.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --WAR-- " + data)
	AppendDataToLog("WRN", reference, data)
}

func LogDebug(reference, data string) {
	color.Blue.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --DEB-- " + data)
}

func LogDirectoryFileCheck(reference string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	logDirectory := filepath.Join(".", "log")
	_, checkPathError := os.Stat(logDirectory)
	logDirectoryExists := checkPathError == nil
	if logDirectoryExists {
		color.Blue.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --DEB-- " + "Log directory already exists ")
		return
	}
	color.Yellow.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --WRN-- " + "Log directory does not exist, creating")
	mkdirError := os.MkdirAll(logDirectory, 0777)
	if mkdirError != nil {
		color.Red.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR--" + "Unable to create directory for log file: " + mkdirError.Error())
		return
	}
}

func AppendDataToLog(logLevel string, reference string, data string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	logNameDateTimeFormat := "2006-01-02"
	logDirectory := filepath.Join(".", "log")
	logFileName := reference + " " + time.Now().Format(logNameDateTimeFormat) + ".log"
	logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
	f, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --WAR-- " + "Log file not present: " + err.Error())
		return
	}
	defer f.Close()
	logData := time.Now().Format("2006-01-02 15:04:05.000   ") + reference + "   " + logLevel + "   " + data
	if _, err := f.WriteString(logData + "\r\n"); err != nil {
		color.Red.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- " + "Cannot write to file: " + err.Error())
	}
}

func AppendDataToErrLog(logLevel string, reference string, data string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	logNameDateTimeFormat := "2006-01-02"
	logDirectory := filepath.Join(".", "log")
	logFileName := reference + " " + time.Now().Format(logNameDateTimeFormat) + ".err"
	logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
	f, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --WAR-- " + "Log file not present: " + err.Error())
		return
	}
	defer f.Close()
	logData := time.Now().Format("2006-01-02 15:04:05.000   ") + reference + "   " + logLevel + "   " + data
	if _, err := f.WriteString(logData + "\r\n"); err != nil {
		color.Red.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- " + "Cannot write to file: " + err.Error())
	}
}

func DeleteOldLogFiles() {
	directory, err := ioutil.ReadDir("log")
	if err != nil {
		LogError("MAIN", "Problem opening log directory")
		return
	}
	now := time.Now()
	logDirectory := filepath.Join(".", "log")
	for _, file := range directory {
		if fileAge := now.Sub(file.ModTime()); fileAge > deleteLogsAfter {
			LogInfo("MAIN", "Deleting old log file "+file.Name()+" with age of "+fileAge.String())
			logFullPath := strings.Join([]string{logDirectory, file.Name()}, "/")
			var err = os.Remove(logFullPath)
			if err != nil {
				LogError("MAIN", "Problem deleting file "+file.Name()+", "+err.Error())
			}
		}
	}
}
