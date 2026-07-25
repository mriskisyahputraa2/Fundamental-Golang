package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location string `json:"location" binding:"required"`
	UserId int `json: "userId" binding:"required"`
}

// Penulisan Array dalam Variabel
var events []Event = []Event{}

// Fungsi untuk simpan event
func (e Event) Save(){
	events = append(events, e)
}


// Fungsi Untuk Menampilkan Semua Event
func GetAllEvent() []Event{
	return events
}