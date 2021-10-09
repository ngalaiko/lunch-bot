package rolls

import (
	"fmt"
	"strings"
	"time"

	"lunch/pkg/users"
)

type Roll struct {
	UserID    string
	PlaceName string
	Time      time.Time
}

func NewRoll(user *users.User, placeID string, now time.Time) *Roll {
	return &Roll{
		UserID:    user.ID,
		PlaceName: placeID,
		Time:      now,
	}
}

func rollFromKey(key string) (*Roll, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("unexpected key format")
	}
	t, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse key time: %w", err)
	}
	return &Roll{
		UserID:    parts[3],
		PlaceName: parts[4],
		Time:      t,
	}, nil
}

func (r *Roll) key() string {
	year, week := r.Time.ISOWeek()
	return fmt.Sprintf("%d/%d/%s/%s/%s", year, week, r.Time.Format(time.RFC3339), r.UserID, r.PlaceName)
}

func (r *Roll) value() string {
	return ""
}
