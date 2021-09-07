package main

import (
	"context"
	"fmt"
	"github.com/go-chassis/kie-client"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
)

var once sync.Once
var configMsg *Config

type Config struct {
	Endpoint   string
	Project    string
	WatchTime  string
	LabelKey   string
	LabelValue string
}

func readConfig(configName string, configPath string, configType string) {
	once.Do(func() {
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)
		viper.SetConfigType(configType)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s\n", err))
		}
		configMsg = &Config{
			Endpoint:   viper.GetString("kieConfig.endpoint"),
			WatchTime:  viper.GetString("kieConfig.watchTime"),
			Project:    viper.GetString("kieConfig.project"),
			LabelKey:   viper.GetString("kieConfig.labelKey"),
			LabelValue: viper.GetString("kieConfig.labelValue"),
		}
	})
}

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}

}

func execute() error {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	readConfig("config", basePath+"/configs", "yaml")
	c, err := kie.NewClient(kie.Config{Endpoint: configMsg.Endpoint})
	if err != nil {
		log.Fatal(err)
		return err
	}
	for {
		lableMap := make(map[string]string)
		lableMap[configMsg.LabelKey] = configMsg.LabelValue
		resp, revision, err := c.List(context.TODO(),
			kie.WithGetProject(configMsg.Project),
			kie.WithLabels(lableMap),
			kie.WithWait(configMsg.WatchTime))
		if err != nil {
			return err
		}
		if resp != nil && resp.Data != nil {
			fmt.Printf("length %v", len(resp.Data))
			fmt.Printf("data: %v \n revision: %v \n", resp.Data, revision)
		}
	}
	return nil
}
