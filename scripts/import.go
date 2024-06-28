package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"myapp/models"
)

func main() {
	// Подключение к базе данных
	db, err := gorm.Open(sqlite.Open("mac_addresses.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Выключаем вывод логов GORM для чистоты вывода
	})
	if err != nil {
		log.Fatalf("failed to connect database, got error: %v", err)
	}
	log.Println("Connected to database")

	// Миграция схемы
	err = db.AutoMigrate(&models.Device{})
	if err != nil {
		log.Fatalf("failed to migrate schema, got error: %v", err)
	}
	log.Println("Migrated schema")

	// Открытие текстового файла
	file, err := os.Open("D:/Education/Projects/Проект ВГУ 2024/MAC_address_lockup/oui.txt")
	if err != nil {
		log.Fatalf("failed to open file, got error: %v", err)
	}
	log.Println("Opened file")
	defer file.Close()

	// Создаем транзакцию для импорта данных
	tx := db.Begin()

	scanner := bufio.NewScanner(file)
	var macAddress, manufacturer, address, city, country string
	var readingAddress bool

	for scanner.Scan() {
		line := scanner.Text()

		// Проверяем строку с MAC-адресом и производителем
		if strings.Contains(line, "(base 16)") || strings.Contains(line, "(hex)") {
			// Получаем MAC-адрес
			macAddress = strings.Fields(line)[0]

			// Получаем производителя
			manufacturer = strings.TrimSpace(strings.Join(strings.Fields(line)[2:], "  "))

			readingAddress = true
			address, city, country = "", "", ""
		} else if readingAddress && strings.TrimSpace(line) != "" {
			// Читаем адресные данные
			if address == "" {
				address = strings.TrimSpace(line)
			} else if city == "" {
				city = strings.TrimSpace(line)
			} else {
				country = strings.TrimSpace(line)
				readingAddress = false

				var existingDevice models.Device
				result := db.Where("mac = ?", macAddress).First(&existingDevice)
				if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					log.Fatalf("failed to check existing device, got error: %v", result.Error)
				}

				if result.RowsAffected > 0 {
					log.Printf("Device with MAC '%s' already exists, skipping insertion", macAddress)
					continue // Пропускаем вставку, если запись уже существует
				}

				device := models.Device{
					MAC:          macAddress,
					Manufacturer: manufacturer,
					Address:      address,
					City:         city,
					Country:      country,
				}

				result = db.Create(&device)
				if result.Error != nil {
					log.Fatalf("failed to create device, got error: %v", result.Error)
				}

			}
		}
	}

	if err := scanner.Err(); err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки при сканировании файла
		log.Fatalf("failed to scan file, got error: %v", err)
	}

	// Фиксируем изменения в базе данных при успешном завершении транзакции
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("failed to commit transaction, got error: %v", err)
	}

	log.Println("Data imported successfully")
}
