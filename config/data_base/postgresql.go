package data_base

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Task struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Title       string `gorm:"not null" json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

type ConfigDB struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	DBName   string `json:"DBName"`
	Password string `json:"Password"`
	SSLMode  string `json:"SSLMODE"`
}

func OpenConnection(configPath string) (*gorm.DB, error) {
	config, err := readConfigFromJSON(configPath)
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		config.Host, config.Port, config.User, config.DBName, config.SSLMode, config.Password,
	)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func readConfigFromJSON(filePath string) (*ConfigDB, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config ConfigDB
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Task{})
}
