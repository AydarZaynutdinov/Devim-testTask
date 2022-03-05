package config

import (
	"encoding/json"
	"math"
	"os"
)

const (
	defaultConfigFilePath = "./config.json"
	configFilePathKey     = "CONFIG_FILE_PATH"
)

type Config struct {
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Radius float64 `json:"radius"`
}

func NewConfig() (*Config, error) {
	configFilePath := os.Getenv(configFilePathKey)
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}
	config := &Config{}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

// Check returns true if {x, y} coordinate inside circle with center in {c.X, c.Y} and radius c.Radius.
//  Returns else in other cases.
func (c Config) Check(x, y int) bool {
	difX := x - c.X
	difY := y - c.Y
	distance := math.Sqrt(math.Pow(float64(difX), 2) + math.Pow(float64(difY), 2))
	return distance <= c.Radius
}
