package events

import (
	"context"
	"log"
	"sync"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
)

type handler func(context.Context, *event) error

type Registry struct {
	handlersGuard *sync.RWMutex
	handlers      map[Type][]handler
}

func NewRegistry() *Registry {
	return &Registry{
		handlersGuard: &sync.RWMutex{},
		handlers:      make(map[Type][]handler),
	}
}

func (r *Registry) PlaceCreated(place *places.Place) {
	r.pub(&event{
		Type:  TypePlaceCreated,
		Place: place,
	})
}

func (r *Registry) OnPlaceCreated(fn func(context.Context, *places.Place) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Place)
	}, TypePlaceCreated)
}

func (r *Registry) RollCreated(roll *rolls.Roll) {
	r.pub(&event{
		Type: TypeRollCreated,
		Roll: roll,
	})
}

func (r *Registry) OnRollCreated(fn func(context.Context, *rolls.Roll) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Roll)
	}, TypeRollCreated)
}

func (r *Registry) BoostCreated(boost *boosts.Boost) {
	r.pub(&event{
		Type:  TypeBoostCreated,
		Boost: boost,
	})
}

func (r *Registry) OnBoostCreated(fn func(context.Context, *boosts.Boost) error) {
	r.sub(func(ctx context.Context, e *event) error {
		return fn(ctx, e.Boost)
	}, TypeBoostCreated)
}

func (r *Registry) pub(evt *event) {
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

func (r *Registry) sub(fn handler, tt ...Type) {
	r.handlersGuard.Lock()
	for _, t := range tt {
		r.handlers[t] = append(r.handlers[t], fn)
	}
	r.handlersGuard.Unlock()
}
