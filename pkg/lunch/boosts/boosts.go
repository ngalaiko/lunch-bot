package boosts

import (
	"fmt"
	"strings"
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/users"
)

type Boost struct {
	UserID    string
	PlaceName places.Name
	Time      time.Time
}

func NewBoost(user *users.User, placeName places.Name, now time.Time) *Boost {
	return &Boost{
		UserID:    user.ID,
		PlaceName: placeName,
		Time:      now,
	}
}

func boostFromKey(key string) (*Boost, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("unexpected key format")
	}
	t, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse key time: %w", err)
	}
	return &Boost{
		UserID:    parts[3],
		PlaceName: places.Name(parts[4]),
		Time:      t,
	}, nil
}

func (r *Boost) key() string {
	year, week := r.Time.ISOWeek()
	return fmt.Sprintf("%d/%d/%s/%s/%s", year, week, r.Time.Format(time.RFC3339), r.UserID, r.PlaceName)
}

func (r *Boost) value() string {
	return ""
}
