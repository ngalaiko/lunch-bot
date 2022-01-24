package storage

import (
	"context"
	"fmt"

	"lunch/pkg/store"
	"lunch/pkg/users"
)

var _ Storage = &dynamoDB{}

type dynamoDB struct {
	storage   *store.DynamoDB
	tableName string
}

func NewDynamoDB(storage *store.DynamoDB, tableName string) *dynamoDB {
	return &dynamoDB{
		storage:   storage,
		tableName: tableName,
	}
}

func (d *dynamoDB) Create(ctx context.Context, user *users.User) error {
	if err := d.storage.Execute(ctx, fmt.Sprintf(`
		INSERT INTO "%s" value {'id': ?, 'name': ?}
	`, d.tableName), user.ID, user.Name); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *dynamoDB) Get(ctx context.Context, id string) (*users.User, error) {
	users := []*users.User{}
	if err := d.storage.Query(ctx, &users, fmt.Sprintf(`SELECT * FROM "%s" WHERE id = ?`, d.tableName), id); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrNotFound
	}
	return users[0], nil
}

func (d *dynamoDB) List(ctx context.Context) (map[string]*users.User, error) {
	uu := []*users.User{}
	if err := d.storage.Query(ctx, &uu, fmt.Sprintf(`SELECT * FROM "%s"`, d.tableName)); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	m := make(map[string]*users.User, len(uu))
	for _, u := range uu {
		m[u.ID] = u
	}
	return m, nil
}
