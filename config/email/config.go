package email

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SMTPConfig struct {
	From     string
	To       []string
	Password string
	SMTPHost string
	SMTPPort string
	Subject  string
	Message  string
}

func ExtractSMTPConfig(configFile string) (*SMTPConfig, error) {
	if v := os.Getenv("SMTP_CONFIG_FILE_PATH"); v != "" {
		configFile = v
	}

	p, _ := filepath.Abs("./src/emailservice/")
	configFilePath := filepath.Join(p, configFile)

	fmt.Println("abs\n", configFilePath)
	f, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := &SMTPConfig{}
	if err = json.NewDecoder(f).Decode(sc); err != nil {
		return nil, err
	}

	return sc, err
}
