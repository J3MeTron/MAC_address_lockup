package services

import (
	"errors"
	"myapp/database"
	"myapp/models"
)

func GetDeviceData(mac string) (*models.Device, error) {
	var device models.Device
	if err := database.DB.Where("mac = ?", mac).First(&device).Error; err != nil {
		return nil, errors.New("device not found")
	}

	return &device, nil
}
