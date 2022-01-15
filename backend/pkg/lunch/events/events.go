package events

import (
	"context"
	"log"
	"sync"

	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
)

type EventHandler func(context.Context, *Event) error

type Registry struct {
	handlersGuard *sync.RWMutex
	handlers      map[Type][]EventHandler
}

func NewRegistry() *Registry {
	return &Registry{
		handlersGuard: &sync.RWMutex{},
		handlers:      make(map[Type][]EventHandler),
	}
}

func (e *Registry) PlaceCreated(place *places.Place) {
	e.emit(&Event{
		Type:  TypePlaceCreated,
		Place: place,
	})
}

func (e *Registry) RollCreated(roll *rolls.Roll) {
	e.emit(&Event{
		Type: TypeRollCreated,
		Roll: roll,
	})
}

func (e *Registry) BoostCreated(boost *boosts.Boost) {
	e.emit(&Event{
		Type:  TypeBoostCreated,
		Boost: boost,
	})
}

func (e *Registry) emit(event *Event) {
	e.handlersGuard.RLock()
	handlers := e.handlers[event.Type]
	e.handlersGuard.RUnlock()

	log.Printf("[INFO] event: '%s'", event.Type.String())

	ctx := context.Background()
	for _, fn := range handlers {
		fn := fn
		go func() {
			if err := fn(ctx, event); err != nil {
				log.Printf("[ERROR] %s", err)
			}
		}()
	}
}

func (e *Registry) Subscribe(fn EventHandler, tt ...Type) {
	e.handlersGuard.Lock()
	for _, t := range tt {
		e.handlers[t] = append(e.handlers[t], fn)
	}
	e.handlersGuard.Unlock()
}
