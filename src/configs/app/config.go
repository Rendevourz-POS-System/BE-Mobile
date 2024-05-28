package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	config *Config
)

type Config struct {
	App         App
	Database    Database
	Email       Email
	Domain      Domain
	Proxy       Proxy
	AccessToken AccessToken
	Image       Image
	Midtrans    Payment
}

type App struct {
	Port         int
	Name         string
	Environment  string
	Locale       string
	Key          string
	Debug        bool
	MigrateKey   string
	UploadFolder string
}

type Database struct {
	Host              []string
	Port              []int
	Database          []string
	Username          []string
	Password          []string
	Schema            []string
	ExternalTableName []string
}

type Email struct {
	File                string
	SenderHost          string
	SenderPort          int
	SenderEmailName     string
	SenderEmailAddress  string
	SenderEmailPassword string
	Attachments         string
}

type Domain struct {
	Name         string
	Protocol     string
	FrontendPath string
	BackendPath  string
	Port         string
}

type AccessToken struct {
	AccessTokenExpireHour       int
	AccessTokenHeaderName       string
	AccessTokenHeaderPrefix     string
	AccessTokenSecret           string
	Key                         string
	RefreshTokenSecret          string
	RefreshTokenExpireHour      int
	VerificationTokenExpireHour int
}

type Proxy struct {
	TrustedProxies []string
}

type Image struct {
	Folder      string
	PetPath     string
	ShelterPath string
	UserPath    string
	ProfilePath string
	TempPath    string
}

type Payment struct {
	ServerKey   string
	Environment string
	Url         string
}

func loadConfig(environment string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(fmt.Sprintf("../config/config-%s", environment))
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func TrustedProxies(servers *gin.Engine) *gin.Engine {
	err := servers.SetTrustedProxies(GetConfig().Proxy.TrustedProxies)
	if err != nil {
		panic("Error setting trusted proxies : " + err.Error())
	}
	return servers
}

func GetConfig() *Config {
	if config == nil {
		v, err := loadConfig(os.Getenv("APP_ENV"))
		if err != nil {
			panic(err)
		}
		err = v.Unmarshal(&config)
		if err != nil {
			log.Printf("unable to decode into struct, %v", err)
			panic(err)
		}
		fmt.Println("Config Environment : ", config.Midtrans)
	}
	return config
}
