package lunch

import (
	"testing"
	"time"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
	"lunch/pkg/users"
)

func TestHistory(t *testing.T) {
	t.Parallel()

	today := time.Date(2021, time.September, 6, 9, 0, 0, 0, time.UTC) // Monday

	testCases := []struct {
		name       string
		boosts     []*boosts.Boost
		rolls      []*rolls.Roll
		time       time.Time
		expected   *rollsHistory
		canBoost   map[users.ID]error
		canRoll    map[users.ID]error
		pointsLeft map[users.ID]int
	}{
		{
			name: "boost and roll",
			time: today.Add(time.Hour),
			rolls: []*rolls.Roll{
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today.Add(time.Minute),
				},
			},
			boosts: []*boosts.Boost{
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today,
				},
			},
			canBoost: map[users.ID]error{
				users.ID("1"): ErrNoPoints,
				users.ID("2"): nil,
			},
			canRoll: map[users.ID]error{
				users.ID("1"): ErrNoPoints,
				users.ID("2"): nil,
			},
			pointsLeft: map[users.ID]int{
				users.ID("1"): 0,
				users.ID("2"): 1,
			},
			expected: &rollsHistory{
				ThisWeekBoosts: []*boosts.Boost{
					{
						UserID:  users.ID("1"),
						PlaceID: places.ID("1"),
						Time:    today,
					},
				},
				RollsPerWeekday: map[time.Weekday][]*rolls.Roll{
					time.Monday: {
						{
							UserID:  users.ID("1"),
							PlaceID: places.ID("1"),
							Time:    today.Add(time.Minute),
						},
					},
				},
				LastRolled: map[places.ID]time.Time{
					places.ID("1"): today.Add(time.Minute),
				},
				ActiveBoosts: map[places.ID]int{},
			},
		},
		{
			name: "roll and boost",
			time: today.Add(time.Hour),
			rolls: []*rolls.Roll{
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today,
				},
			},
			boosts: []*boosts.Boost{
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today.Add(time.Minute),
				},
			},
			canBoost: map[users.ID]error{
				users.ID("1"): ErrNoPoints,
				users.ID("2"): nil,
			},
			canRoll: map[users.ID]error{
				users.ID("1"): ErrNoPoints,
				users.ID("2"): nil,
			},
			pointsLeft: map[users.ID]int{
				users.ID("1"): 0,
				users.ID("2"): 1,
			},
			expected: &rollsHistory{
				ThisWeekBoosts: []*boosts.Boost{
					{
						UserID:  users.ID("1"),
						PlaceID: places.ID("1"),
						Time:    today.Add(time.Minute),
					},
				},
				RollsPerWeekday: map[time.Weekday][]*rolls.Roll{
					time.Monday: {
						{
							UserID:  users.ID("1"),
							PlaceID: places.ID("1"),
							Time:    today,
						},
					},
				},
				LastRolled: map[places.ID]time.Time{
					places.ID("1"): today,
				},
				ActiveBoosts: map[places.ID]int{
					places.ID("1"): 1,
				},
			},
		},
		{
			name: "two rolls",
			time: today.Add(2 * time.Hour),
			rolls: []*rolls.Roll{
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today,
				},
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today.Add(time.Hour),
				},
			},
			canBoost: map[users.ID]error{
				users.ID("1"): ErrNoPoints,
				users.ID("2"): nil,
			},
			canRoll: map[users.ID]error{
				users.ID("1"): ErrNoPoints,
				users.ID("2"): nil,
			},
			pointsLeft: map[users.ID]int{
				users.ID("1"): 0,
				users.ID("2"): 1,
			},
			expected: &rollsHistory{
				ThisWeekBoosts: []*boosts.Boost{},
				RollsPerWeekday: map[time.Weekday][]*rolls.Roll{
					time.Monday: {
						{
							UserID:  users.ID("1"),
							PlaceID: places.ID("1"),
							Time:    today,
						},
						{
							UserID:  users.ID("1"),
							PlaceID: places.ID("1"),
							Time:    today.Add(time.Hour),
						},
					},
				},
				LastRolled: map[places.ID]time.Time{
					places.ID("1"): today.Add(time.Hour),
				},
				ActiveBoosts: map[places.ID]int{},
			},
		},
		{
			name: "one boost",
			time: today,
			boosts: []*boosts.Boost{
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today,
				},
			},
			canBoost: map[users.ID]error{
				users.ID("1"): ErrNoPoints,
				users.ID("2"): nil,
			},
			canRoll: map[users.ID]error{
				users.ID("1"): nil,
				users.ID("2"): nil,
			},
			pointsLeft: map[users.ID]int{
				users.ID("1"): 0,
				users.ID("2"): 1,
			},
			expected: &rollsHistory{
				ThisWeekBoosts: []*boosts.Boost{
					{
						UserID:  users.ID("1"),
						PlaceID: places.ID("1"),
						Time:    today,
					},
				},
				RollsPerWeekday: map[time.Weekday][]*rolls.Roll{},
				LastRolled:      map[places.ID]time.Time{},
				ActiveBoosts: map[places.ID]int{
					places.ID("1"): 1,
				},
			},
		},
		{
			name: "one roll",
			time: today,
			rolls: []*rolls.Roll{
				{
					UserID:  users.ID("1"),
					PlaceID: places.ID("1"),
					Time:    today,
				},
			},
			canBoost: map[users.ID]error{
				users.ID("1"): nil,
				users.ID("2"): nil,
			},
			canRoll: map[users.ID]error{
				users.ID("1"): nil,
				users.ID("2"): nil,
			},
			pointsLeft: map[users.ID]int{
				users.ID("1"): 1,
				users.ID("2"): 1,
			},
			expected: &rollsHistory{
				ThisWeekBoosts: []*boosts.Boost{},
				RollsPerWeekday: map[time.Weekday][]*rolls.Roll{
					time.Monday: {
						{
							UserID:  users.ID("1"),
							PlaceID: places.ID("1"),
							Time:    today,
						},
					},
				},
				LastRolled: map[places.ID]time.Time{
					places.ID("1"): today,
				},
				ActiveBoosts: map[places.ID]int{},
			},
		},
		{
			name: "no nothing",
			time: today,
			canBoost: map[users.ID]error{
				users.ID("1"): nil,
			},
			canRoll: map[users.ID]error{
				users.ID("1"): nil,
			},
			pointsLeft: map[users.ID]int{
				users.ID("1"): 1,
			},
			expected: &rollsHistory{
				ThisWeekBoosts:  []*boosts.Boost{},
				RollsPerWeekday: map[time.Weekday][]*rolls.Roll{},
				LastRolled:      map[places.ID]time.Time{},
				ActiveBoosts:    map[places.ID]int{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := buildHistory(tc.rolls, tc.boosts, tc.time)
			assertEqual(t, tc.expected, actual)

			for uID, expected := range tc.canBoost {
				assertEqual(t, expected, actual.CanBoost(uID, tc.time))
			}

			for uID, expected := range tc.canRoll {
				assertEqual(t, expected, actual.CanRoll(uID, tc.time))
			}

			for uID, expected := range tc.pointsLeft {
				assertEqual(t, expected, actual.pointsLeft(uID))
			}
		})
	}
}
