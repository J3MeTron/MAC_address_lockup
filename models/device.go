package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	MAC          string
	Manufacturer string
	Address      string
	City         string
	Country      string
}
