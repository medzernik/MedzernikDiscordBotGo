package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	ServerInfo struct {
		Prefix        string `yaml:"prefix"`
		GuildIDNumber string `yaml:"guildIDNumber"`
		WeatherAPIKey string `yaml:"weatherAPIKey"`
		ServerToken   string `yaml:"serverToken"`
	} `yaml:"serverInfo"`
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
		ChannelLogID string `yaml:"channelLogID"`
	} `yaml:"channelLog"`
	MuteFunction struct {
		MuteRoleID           string `yaml:"MuteRoleID"`
		TrustedMutingEnabled bool   `yaml:"trustedMutingEnabled"`
	} `yaml:"muteFunction"`
	AutoLocker struct {
		Enabled          bool         `yaml:"enabled"`
		TimeDayUnlock    time.Weekday `yaml:"timeDayUnlock"`
		TimeHourUnlock   int          `yaml:"timeHourUnlock"`
		TimeMinuteUnlock int          `yaml:"timeMinuteUnlock"`
		TimeDayLock      time.Weekday `yaml:"timeDayLock"`
		TimeHourLock     int          `yaml:"timeHourLock"`
		TimeMinuteLock   int          `yaml:"timeMinuteLock"`
	} `yaml:"autoLocker"`
}

var Cfg Config

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

	fmt.Println(Cfg)
}
