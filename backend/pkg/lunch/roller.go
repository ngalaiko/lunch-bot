package lunch

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"lunch/pkg/lunch/boosts"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/events"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	"lunch/pkg/lunch/rolls"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/lunch/rooms"
	storage_rooms "lunch/pkg/lunch/rooms/storage"
	"lunch/pkg/users"
	storage_users "lunch/pkg/users/storage"
)

var (
	ErrNoPoints = fmt.Errorf("no points left")
	ErrNoPlaces = fmt.Errorf("no places to choose from")
)

type Roller struct {
	*registry

	placesStore *storage_places.Storage
	rollsStore  *storage_rolls.Storage
	boostsStore *storage_boosts.Storage
	usersStore  storage_users.Storage
	roomsStore  *storage_rooms.Storage

	rand *rand.Rand
}

func New(eventsStorage events.Storage, usersStore storage_users.Storage) *Roller {
	return &Roller{
		registry:    newEventsRegistry(),
		placesStore: storage_places.New(eventsStorage),
		rollsStore:  storage_rolls.New(eventsStorage),
		boostsStore: storage_boosts.New(eventsStorage),
		roomsStore:  storage_rooms.New(eventsStorage),
		usersStore:  usersStore,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Roller) CreateRoom(ctx context.Context, name string) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	room := rooms.New(user.ID, name)
	if err := r.roomsStore.Create(ctx, room); err != nil {
		return fmt.Errorf("failed to store place: %w", err)
	}

	r.RoomCreated(&Room{
		Room: room,
		User: user,
		Members: []*users.User{
			user,
		},
	})

	return nil
}

func (r *Roller) LeaveRoom(ctx context.Context, roomID rooms.ID) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	room, err := r.roomsStore.Room(ctx, roomID)
	if errors.Is(err, storage_rooms.ErrNotFound) {
		return fmt.Errorf("room not found")
	} else if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if err := r.roomsStore.Leave(ctx, user, roomID); err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	roomView := &Room{
		Room: room,
		User: allUsers[room.UserID],
	}

	for uid := range room.MemberIDs {
		roomView.Members = append(roomView.Members, allUsers[uid])
	}

	r.RoomUpdated(roomView)

	return nil
}

func (r *Roller) JoinRoom(ctx context.Context, roomID rooms.ID) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	room, err := r.roomsStore.Room(ctx, roomID)
	if errors.Is(err, storage_rooms.ErrNotFound) {
		return fmt.Errorf("room not found")
	} else if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if err := r.roomsStore.Join(ctx, user, roomID); err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	roomView := &Room{
		Room: room,
		User: allUsers[room.UserID],
	}

	for uid := range room.MemberIDs {
		roomView.Members = append(roomView.Members, allUsers[uid])
	}

	r.RoomUpdated(roomView)

	return nil
}

func (r *Roller) ListRooms(ctx context.Context) ([]*Room, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}
	roomIDs, err := r.roomsStore.Rooms(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}
	rooms := make([]*rooms.Room, len(roomIDs))
	for id := range roomIDs {
		room, err := r.roomsStore.Room(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get room: %w", err)
		}
		rooms = append(rooms, room)
	}
	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	result := make([]*Room, 0, len(roomIDs))
	for _, room := range rooms {
		roomView := &Room{
			Room: room,
			User: allUsers[room.UserID],
		}
		for uid := range room.MemberIDs {
			roomView.Members = append(roomView.Members, allUsers[uid])
		}
		result = append(result, roomView)
	}
	return result, nil
}

func (r *Roller) CreatePlace(ctx context.Context, roomID rooms.ID, name string) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	place := places.NewPlace(roomID, user.ID, name)
	if err := r.placesStore.Create(ctx, place); err != nil {
		return fmt.Errorf("failed to store place: %w", err)
	}

	r.PlaceCreated(&Place{
		Place: place,
		User:  user,
	})

	return nil
}

func (r *Roller) ListRolls(ctx context.Context, roomID rooms.ID) ([]*Roll, error) {
	allRolls, err := r.rollsStore.Rolls(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}
	if len(allRolls) == 0 {
		return nil, nil
	}

	allPlaces, err := r.placesStore.Places(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	rolls := make([]*Roll, 0, len(allRolls))
	for _, roll := range allRolls {
		rolls = append(rolls, &Roll{
			Roll:  roll,
			User:  allUsers[roll.UserID],
			Place: allPlaces[roll.PlaceID],
		})
	}

	return rolls, nil
}

func (r *Roller) ListBoosts(ctx context.Context, roomID rooms.ID) ([]*Boost, error) {
	allBoosts, err := r.boostsStore.Boosts(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	allPlaces, err := r.placesStore.Places(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list places: %w", err)
	}

	boosts := make([]*Boost, 0, len(allBoosts))
	for _, b := range allBoosts {
		boosts = append(boosts, &Boost{
			Boost: b,
			User:  allUsers[b.UserID],
			Place: allPlaces[b.PlaceID],
		})
	}

	return boosts, nil
}

func filterNonDeletedPlaces(pp map[places.ID]*places.Place) map[places.ID]*places.Place {
	result := make(map[places.ID]*places.Place, len(pp))
	for id, place := range pp {
		if !place.IsDeleted {
			result[id] = place
		}
	}
	return result
}

func (r *Roller) ListPlaces(ctx context.Context, roomID rooms.ID, now time.Time) ([]*Place, error) {
	allPlaces, err := r.placesStore.Places(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}

	allPlaces = filterNonDeletedPlaces(allPlaces)

	if len(allPlaces) == 0 {
		return nil, ErrNoPlaces
	}

	allBoosts, err := r.boostsStore.Boosts(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	allRolls, err := r.rollsStore.Rolls(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}

	allUsers, err := r.usersStore.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	history := buildHistory(allRolls, allBoosts, now)
	weights := history.getWeights(allPlaces, now)
	weightsSum := 0.0
	for _, weight := range weights {
		weightsSum += weight
	}

	views := make([]*Place, 0, len(allPlaces))
	for i, weight := range weights {
		if _, ok := allPlaces[i]; !ok {
			continue
		}

		chance := weight / weightsSum
		views = append(views, &Place{
			Place:  allPlaces[i],
			User:   allUsers[allPlaces[i].UserID],
			Chance: chance,
		})
	}

	return views, nil
}

func (r *Roller) CreateBoost(ctx context.Context, roomID rooms.ID, placeID places.ID, now time.Time) error {
	user, ok := users.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected to find who in the context")
	}

	place, err := r.placesStore.Place(ctx, roomID, placeID)
	if err != nil {
		return fmt.Errorf("failed to get place: %w", err)
	}

	allRolls, err := r.rollsStore.Rolls(ctx, roomID)
	if err != nil {
		return fmt.Errorf("failed to list rolls: %w", err)
	}

	allBoosts, err := r.boostsStore.Boosts(ctx, roomID)
	if err != nil {
		return fmt.Errorf("failed to list boosts: %w", err)
	}

	history := buildHistory(allRolls, allBoosts, now)
	if err := history.CanBoost(user.ID, now); err != nil {
		return fmt.Errorf("can't boost any more: %w", err)
	}

	boost := boosts.NewBoost(user.ID, roomID, placeID, now)
	if err := r.boostsStore.Create(ctx, boost); err != nil {
		return fmt.Errorf("failed to store boost: %w", err)
	}

	r.BoostCreated(&Boost{
		Boost: boost,
		User:  user,
		Place: place,
	})

	return nil
}

func (r *Roller) CreateRoll(ctx context.Context, roomID rooms.ID, now time.Time) (*Roll, error) {
	user, ok := users.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("expected to find who in the context")
	}
	allPlaces, err := r.placesStore.Places(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list names: %w", err)
	}
	allPlaces = filterNonDeletedPlaces(allPlaces)

	if len(allPlaces) == 0 {
		return nil, ErrNoPlaces
	}

	allRolls, err := r.rollsStore.Rolls(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %w", err)
	}

	allBoosts, err := r.boostsStore.Boosts(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to list boosts: %w", err)
	}

	history := buildHistory(allRolls, allBoosts, now)
	if err := history.CanRoll(user.ID, now); err != nil {
		return nil, fmt.Errorf("failed to validate rules: %w", err)
	}

	weights := history.getWeights(allPlaces, now)
	randomIndex := weightedRandom(r.rand, weights)
	randomPlace := allPlaces[randomIndex]

	roll := rolls.NewRoll(user.ID, roomID, randomPlace.ID, now)
	if err := r.rollsStore.Create(ctx, roll); err != nil {
		return nil, fmt.Errorf("failed to store roll result: %w", err)
	}

	rollView := &Roll{
		Roll:  roll,
		User:  user,
		Place: randomPlace,
	}

	r.RollCreated(rollView)

	return rollView, nil
}

// weightedRandom returns a random index i from the slice of weights, proportional to the weights[i] value.
func weightedRandom(rand *rand.Rand, weights map[places.ID]float64) places.ID {
	weightsSum := 0.0
	for _, weight := range weights {
		weightsSum += weight
	}

	remainingDistance := rand.Float64() * weightsSum
	for i, weight := range weights {
		remainingDistance -= weight
		if remainingDistance < 0 {
			return i
		}
	}

	panic("should never happen")
}
