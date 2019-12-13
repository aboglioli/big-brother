package config

import (
	"encoding/json"
	"os"
	"sync"
)

type service struct {
	Port int16 `json:"port"`
}

type Configuration struct {
	Discovery service `json:"discovery"`
	User      service `json:"user"`

	MongoURL        string `json:"mongoUrl"`
	MongoAuthSource string `json:"mongoAuthSource"`
	MongoUsername   string `json:"mongoUsername"`
	MongoPassword   string `json:"mongoPassword"`

	RabbitURL string `json:"rabbitUrl"`

	RedisURL      string `json:"redisUrl"`
	RedisPassword string `json:"redisPassword"`
	RedisDB       int    `json:"redisDb"`

	AuthEnabled bool   `json:"authEnabled"`
	AuthURL     string `json:"authUrl"`
}

var once sync.Once
var config *Configuration

func Get() Configuration {
	once.Do(func() {
		config = &Configuration{
			Discovery: service{Port: 1492},
			User:      service{Port: 3344},

			MongoURL:        "mongodb://localhost:27017",
			MongoAuthSource: "admin",
			MongoUsername:   "admin",
			MongoPassword:   "admin",

			RabbitURL: "amqp://guest:guest@localhost:5672",

			RedisURL:      "localhost:6379",
			RedisPassword: "",
			RedisDB:       0,

			AuthEnabled: false,
		}

		file, err := os.Open("config.json")
		if err == nil && file != nil {
			json.NewDecoder(file).Decode(config)
		}
		defer file.Close()
	})

	return *config
}
