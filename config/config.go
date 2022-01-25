// Package config Configuration module that holds the configuration logic
package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

// Config Configuration type structure to use in memory.
type Config struct {
	ServerInfo struct {
		GuildIDNumber string `yaml:"guildIDNumber"`
		WeatherAPIKey string `yaml:"weatherAPIKey"`
		ServerToken   string `yaml:"serverToken"`
		BotStatus     string `yaml:"botStatus"`
		BotLogo       string `yaml:"botLogo"`
		BotName       string `yaml:"botName"`
		LogLevel      string `yaml:"logLevel"`
	} `yaml:"serverInfo"`
	Modules struct {
		Administration     bool `yaml:"administration"`
		Logging            bool `yaml:"logging"`
		Lottery            bool `yaml:"lottery"`
		Planning           bool `yaml:"planning"`
		Weather            bool `yaml:"weather"`
		Purge              bool `yaml:"purge"`
		COVIDSlovakInfo    bool `yaml:"COVIDSlovakInfo"`
		TimedChannelUnlock bool `yaml:"timedChannelUnlock"`
	}
	RoleAdmin struct {
		RoleAdminID string `yaml:"roleAdminID"`
		RoleModID   string `yaml:"roleModID"`
	} `yaml:"roleAdmin"`
	RoleTrusted struct {
		RoleTrustedID1   string `yaml:"roleTrustedID1"`
		RoleTrustedID2   string `yaml:"roleTrustedID2"`
		RoleTrustedID3   string `yaml:"roleTrustedID3"`
		RoleTrustedID4   string `yaml:"roleTrustedID4"`
		ChannelTrustedID string `yaml:"channelTrustedID"`
	} `yaml:"roleTrusted"`
	ChannelLog struct {
		ChannelLogID   string `yaml:"channelLogID"`
		GamePlannedLog string `yaml:"gamePlannedLog"`
	} `yaml:"channelLog"`
	AutoLocker struct {
		Enabled              bool         `yaml:"enabled"`
		AutoUnlockTrustedID1 bool         `yaml:"autoUnlockTrusted1"`
		TimeDayUnlock        time.Weekday `yaml:"timeDayUnlock"`
		TimeHourUnlock       int          `yaml:"timeHourUnlock"`
		TimeMinuteUnlock     int          `yaml:"timeMinuteUnlock"`
		TimeDayLock          time.Weekday `yaml:"timeDayLock"`
		TimeHourLock         int          `yaml:"timeHourLock"`
		TimeMinuteLock       int          `yaml:"timeMinuteLock"`
	} `yaml:"autoLocker"`
	LotteryChecker struct {
		Enabled         bool         `yaml:"enabled"`
		TimeDayStart    time.Weekday `yaml:"timeDayStart"`
		TimeHourStart   int          `yaml:"timeHourStart"`
		TimeMinuteStart int          `yaml:"timeMinuteStart"`
		TimeDayEnd      time.Weekday `yaml:"timeDayEnd"`
		TimeHourEnd     int          `yaml:"timeHourEnd"`
		TimeMinuteEnd   int          `yaml:"timeMinuteEnd"`
	}
}

var Cfg Config

// LoadConfig Loads the config file. It must be in the root of the directory, next to the main executable.
func LoadConfig() {
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Config loaded.")

}

// SaveConfig This function allows you to save the current config to the file
func SaveConfig() {
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(&Cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Configure This function allows you to set values in the config straight from discord.
func Configure() {

}
