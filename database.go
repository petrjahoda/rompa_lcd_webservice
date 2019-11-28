package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type WorkplacePort struct {
	OID          int `gorm:"column:OID"`
	DevicePortId int `gorm:"column:DevicePortID"`
}

func (WorkplacePort) TableName() string {
	return "workplace_port"
}

type DeviceInputAnalog struct {
	OID          int       `gorm:"column:OID"`
	DevicePortId int       `gorm:"column:DevicePortID"`
	DT           time.Time `gorm:"column:DT"`
	Data         float32   `gorm:"column:Data"`
}

func (DeviceInputAnalog) TableName() string {
	return "device_input_analog"
}

type TerminalInputOrder struct {
	OID      int       `gorm:"column:OID"`
	DTS      time.Time `gorm:"column:DTS"`
	DTE      time.Time `gorm:"column:DTE; default:null"`
	Interval float32   `gorm:"column:Interval"`
	OrderID  int       `gorm:"column:OrderID"`
	UserID   int       `gorm:"column:UserID"`
	DeviceID int       `gorm:"column:DeviceID"`
}

func (TerminalInputOrder) TableName() string {
	return "terminal_input_order"
}

type User struct {
	OID   int    `gorm:"column:OID"`
	Login string `gorm:"column:Login"`
	Name  string `gorm:"column:Name"`
}

func (User) TableName() string {
	return "user"
}

type Order struct {
	OID     int    `gorm:"column:OID"`
	Name    string `gorm:"column:Name"`
	Barcode string `gorm:"column:Barcode"`
}

func (Order) TableName() string {
	return "order"
}

type Workplace struct {
	OID                 int    `gorm:"column:OID"`
	Name                string `gorm:"column:Name"`
	WorkplaceDivisionId int    `gorm:"column:WorkplaceDivisionID"`
	DeviceID            int    `gorm:"column:DeviceID"`
}

func (Workplace) TableName() string {
	return "workplace"
}

type WorkplaceState struct {
	OID         int       `gorm:"column:OID"`
	StateID     int       `gorm:"column:StateID"`
	WorkplaceID int       `gorm:"column:WorkplaceID"`
	DTS         time.Time `gorm:"column:DTS"`
}

func (WorkplaceState) TableName() string {
	return "workplace_state"
}

func CheckDatabaseType() (string, string) {
	var connectionString string
	var dialect string
	if DatabaseType == "postgres" {
		connectionString = "host=" + DatabaseIpAddress + " sslmode=disable port=" + DatabasePort + " user=" + DatabaseLogin + " dbname=" + DatabaseName + " password=" + DatabasePassword
		dialect = "postgres"
	} else if DatabaseType == "mysql" {
		connectionString = DatabaseLogin + ":" + DatabasePassword + "@tcp(" + DatabaseIpAddress + ":" + DatabasePort + ")/" + DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
		dialect = "mysql"
	}
	return connectionString, dialect
}
