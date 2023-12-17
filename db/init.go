package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	idemsEnv "api-medese/config"
	models "api-medese/models"
)

var (
	// DBConn is DBConn
	DBConn *gorm.DB
)

// InitDatabase is InitDatabase
func InitDatabase() {
	var err error
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=require TimeZone=Asia/Bangkok", idemsEnv.GetCfg().DbHost, idemsEnv.GetCfg().DbUsername, idemsEnv.GetCfg().DbPass, idemsEnv.GetCfg().DbName, idemsEnv.GetCfg().DbPort)
	fmt.Println(dsn)
	DBConn, err = gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection successfully")
	DBConn.AutoMigrate(
		models.Users{},
		models.LogProduct{},
		models.LogStock{},
		models.LogManualOrder{},
	)
}
