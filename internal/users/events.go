package users

import (
	"github.com/aboglioli/big-brother/pkg/events"
)

type UserEvent struct {
	events.Event
	User *User `json:"user"`
}

func NewUserEvent(u *User, eventType string) *UserEvent {
	return &UserEvent{
		Event: events.Event{
			Type: eventType,
		},
		User: u,
	}
}
