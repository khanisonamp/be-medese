package main

import (
	"api-medese/db"
	"api-medese/router"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	app "api-medese/framework/fiber"

	idemsEnv "api-medese/config"

	apiName "api-medese/api/v1"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func main() {

	initTimeZone()

	vp := viper.New()
	setEnv := idemsEnv.NewViperContext(vp)
	setEnv.InitConfig()

	if err := idemsEnv.GetCfg().Load(setEnv); err != nil {
		log.Fatal("Load configuration error")
	}

	db.InitDatabase()

	// CrestedLogStockAuto CronJob
	cronJob := cron.New()
	cronJob.AddFunc("30 12 * * *", func() {
		apiName.CrestedLogStockAuto()
	})
	cronJob.Start()
	defer cronJob.Stop()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r := app.NewFiberApp()
	router.SetRouter(r)
	if err := r.Listen(getPort()); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Sprint("listen:", err.Error()))
	}
}

func initConfig() {
	fmt.Println(os.Getenv("ENV"))
	switch os.Getenv("ENV") {
	case "dev":
		os.Setenv("ENV", "dev")
		viper.SetConfigName("config_dev")
	case "uat":
		os.Setenv("ENV", "uat")
		viper.SetConfigName("config")
	default:
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8081"
		fmt.Println("Port " + port)
	}
	return ":" + port
}
