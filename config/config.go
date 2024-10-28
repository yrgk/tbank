package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigStruct struct {
	DSN             string
	Port            string
	TbankApiUrl     string
	CrmApiUrl       string
	NotificationURL string
}

var Config ConfigStruct

func GetConfig() {
	// if err := godotenv.Load("../../.env"); err != nil {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not found: %s", err)
	}

	DSN := os.Getenv("DSN")
	port := os.Getenv("PORT")
	TbankApiUrl := os.Getenv("TbankApiUrl")
	NotificationURL := os.Getenv("NotificationURL")
	CrmApiUrl := os.Getenv("CrmApiUrl")

	Config = ConfigStruct{
		DSN:             DSN,
		Port:            port,
		TbankApiUrl:     TbankApiUrl,
		NotificationURL: NotificationURL,
		CrmApiUrl: CrmApiUrl,
	}
}
