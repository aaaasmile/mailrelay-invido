package conf

import (
	"encoding/json"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceURL    string
	DebugVerbose  bool
	TokenFilename string
	SecretConfig  *SecretConfig
}

type SecretConfig struct {
	EMailRelay      string `json:"email-relay"`
	EMailLogin      string `json:"email-login"`
	EmailPassword   string `json:"email-password"`
	ServiceUser     string `json:"service-user"`
	ServicePassword string `json:"service-password"`
	RemoteSendHost  string `json:"remote-host"`
	HostName        string `json:"host"`
}

var Current = &Config{}

func ReadConfig(configfile string) (*Config, error) {
	_, err := os.Stat(configfile)
	if err != nil {
		return nil, err
	}
	if _, err := toml.DecodeFile(configfile, &Current); err != nil {
		return nil, err
	}

	Current.SecretConfig, err = readSecretFromJSONFile(Current.TokenFilename)
	if err != nil {
		return nil, err
	}
	return Current, nil
}

func readSecretFromJSONFile(cfgFile string) (*SecretConfig, error) {
	log.Println("Read Secret configuration file ", cfgFile)
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	info := SecretConfig{}

	err = json.NewDecoder(f).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
