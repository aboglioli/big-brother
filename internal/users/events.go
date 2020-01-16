package users

import (
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/models"
)

type UserEvent struct {
	events.Event
	User *models.User `json:"user"`
}

func NewUserEvent(u *models.User, eventType string) *UserEvent {
	return &UserEvent{
		Event: events.Event{
			Type: eventType,
		},
		User: u,
	}
}
