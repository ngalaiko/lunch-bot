package lunch

import (
	"context"
	"log"
	"sync"
)

type handler func(context.Context, *event) error

type registry struct {
	handlersGuard *sync.RWMutex
	handlers      map[Type][]handler
}

func newEventsRegistry() *registry {
	return &registry{
		handlersGuard: &sync.RWMutex{},
		handlers:      make(map[Type][]handler),
	}
}

func (r *registry) PlaceCreated(place *Place) {
	r.pub(&event{
		Type:  TypePlaceCreated,
		Place: place,
	})
}

func (r *registry) OnPlaceCreated(fn func(context.Context, *Place) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Place)
	}, TypePlaceCreated)
}

func (r *registry) RoomUpdated(room *Room) {
	r.pub(&event{
		Type: TypeRoomUpdated,
		Room: room,
	})
}

func (r *registry) OnRoomUpdated(fn func(context.Context, *Room) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Room)
	}, TypeRoomUpdated)
}

func (r *registry) RoomCreated(room *Room) {
	r.pub(&event{
		Type: TypeRoomCreated,
		Room: room,
	})
}

func (r *registry) OnRoomCreated(fn func(context.Context, *Room) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Room)
	}, TypeRoomCreated)
}

func (r *registry) RollCreated(roll *Roll) {
	r.pub(&event{
		Type: TypeRollCreated,
		Roll: roll,
	})
}

func (r *registry) OnRollCreated(fn func(context.Context, *Roll) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Roll)
	}, TypeRollCreated)
}

func (r *registry) BoostCreated(boost *Boost) {
	r.pub(&event{
		Type:  TypeBoostCreated,
		Boost: boost,
	})
}

func (r *registry) OnBoostCreated(fn func(context.Context, *Boost) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Boost)
	}, TypeBoostCreated)
}

func (r *registry) pub(evt *event) {
	r.handlersGuard.RLock()
	handlers := r.handlers[evt.Type]
	r.handlersGuard.RUnlock()

	log.Printf("[INFO] event: '%s'", evt.Type.String())

	ctx := context.Background()
	for _, fn := range handlers {
		fn := fn
		go func() {
			if err := fn(ctx, evt); err != nil {
				log.Printf("[ERROR] error handling %s: %s", evt.Type.String(), err)
			}
		}()
	}
}

func (r *registry) sub(fn handler, tt ...Type) {
	r.handlersGuard.Lock()
	for _, t := range tt {
		r.handlers[t] = append(r.handlers[t], fn)
	}
	r.handlersGuard.Unlock()
}
