package config

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
)

type service struct {
	Port int16 `json:"port"`
}

type Configuration struct {
	User service `json:"user"`

	MongoURL        string `json:"mongoUrl"`
	MongoAuthSource string `json:"mongoAuthSource"`
	MongoUsername   string `json:"mongoUsername"`
	MongoPassword   string `json:"mongoPassword"`

	PostgresURL      string `json:"postgresUrl"`
	PostgresUsername string `json:"postgresUsername"`
	PostgresPassword string `json:"postgresPassword"`

	RabbitURL string `json:"rabbitUrl"`

	RedisURL      string `json:"redisUrl"`
	RedisPassword string `json:"redisPassword"`
	RedisDB       int    `json:"redisDb"`

	AuthEnabled bool   `json:"authEnabled"`
	AuthURL     string `json:"authUrl"`
	JWTSecret   []byte `json:"jwtSecret"`
	BcryptCost  int    `json:"bcryptCost"`
}

var once sync.Once
var config *Configuration

func Get() Configuration {
	once.Do(func() {
		config = &Configuration{
			User: service{Port: 3344},

			MongoURL:        "mongodb://localhost:27017",
			MongoAuthSource: "admin",
			MongoUsername:   "admin",
			MongoPassword:   "admin",

			PostgresURL:      "localhost:5432",
			PostgresUsername: "admin",
			PostgresPassword: "admin",

			RabbitURL: "amqp://guest:guest@localhost:5672",

			RedisURL:      "localhost:6379",
			RedisPassword: "",
			RedisDB:       0,

			AuthEnabled: false,
			JWTSecret:   []byte("my_secret_key"),
			BcryptCost:  bcrypt.DefaultCost,
		}

		file, err := os.Open("config.json")
		if err == nil && file != nil {
			json.NewDecoder(file).Decode(config)
		}
		defer file.Close()
	})

	return *config
}
