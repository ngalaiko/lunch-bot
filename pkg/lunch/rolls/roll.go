package rolls

import (
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

type Roll struct {
	ID        string      `dynamodbav:"id"`
	UserID    string      `dynamodbav:"user_id"`
	PlaceName places.Name `dynamodbav:"place_name"`
	Time      time.Time   `dynamodbav:"time,unixtime"`
}

func NewRoll(user *users.User, placeName places.Name, now time.Time) *Roll {
	return &Roll{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		PlaceName: placeName,
		Time:      now,
	}
}
