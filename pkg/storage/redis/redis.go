package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"atomys.codes/webhooked/internal/valuable"
)

type storage struct {
	client *redis.Client
	config *config
	ctx    context.Context
}

type config struct {
	Host     string
	Port     string
	Database int
	Username valuable.Valuable
	Password valuable.Valuable
	Key      string
}

func NewStorage(configRaw map[string]interface{}) (*storage, error) {

	newClient := storage{
		config: &config{},
		ctx:    context.Background(),
	}

	if err := valuable.Decode(configRaw, &newClient.config); err != nil {
		return nil, err
	}

	newClient.client = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", newClient.config.Host, newClient.config.Port),
			Username: newClient.config.Username.First(),
			Password: newClient.config.Password.First(),
			DB:       newClient.config.Database,
		},
	)

	// Ping Redis for testing config
	if err := newClient.client.Ping(newClient.ctx).Err(); err != nil {
		return nil, err
	}

	return &newClient, nil
}

// Name is the function for identified if the storage config is define in the webhooks
// @return name of the storage
func (c storage) Name() string {
	return "redis"
}

// Push is the function for push data in the storage
// A run is made from external caller
// @param value that will be pushed
// @return an error if the push failed
func (c storage) Push(value interface{}) error {
	if err := c.client.RPush(c.ctx, c.config.Key, value).Err(); err != nil {
		return err
	}

	return nil
}
