package admin

import (
	"encoding/json"
	"os"
)

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminConfig struct {
	Admins []AdminUser `json:"admins"`
}

func LoadConfig(path string) (*AdminConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg AdminConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *AdminConfig) Authenticate(username, password string) bool {
	for _, a := range c.Admins {
		if a.Username == username && a.Password == password {
			return true
		}
	}
	return false
}
