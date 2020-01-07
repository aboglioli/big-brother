package business

import (
	"github.com/aboglioli/big-brother/db/models"
)

type Business struct {
	models.Base
	Name string `json:"name" bson:"name"`
}
