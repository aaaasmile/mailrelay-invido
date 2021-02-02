package conf

import (
	"bytes"
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

func ReadConfig(configfile string, rawsecret *[]byte) (*Config, error) {
	_, err := os.Stat(configfile)
	if err != nil {
		return nil, err
	}
	if _, err := toml.DecodeFile(configfile, &Current); err != nil {
		return nil, err
	}
	info := SecretConfig{}

	err = json.NewDecoder(bytes.NewReader(*rawsecret)).Decode(&info)
	if err != nil {
		return nil, err
	}
	Current.SecretConfig = &info
	log.Println("Relay sistem is ", info.RemoteSendHost)

	return Current, nil
}
