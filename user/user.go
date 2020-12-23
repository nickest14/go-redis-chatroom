package user

import (
	"context"
	rediswrap "go-redis-chatroom/redis"

	"github.com/go-redis/redis/v8"
)

const (
	usersKey       = "users"
	userChannelFmt = "user:%s:channels"
	ChannelsKey    = "channels"
)

// User struct
type User struct {
	name             string
	channelsHandler  *redis.PubSub
	stopListenerChan chan struct{}
	listening        bool
	MessageChan      chan redis.Message
}

//Connect connect user to user channels on redis
func Connect(name string) (*User, error) {
	ctx := context.Background()
	rdb := rediswrap.Client
	if _, err := rdb.SAdd(ctx, usersKey, name).Result(); err != nil {
		return nil, err
	}
	u := &User{
		name:             name,
		stopListenerChan: make(chan struct{}),
		MessageChan:      make(chan redis.Message),
	}

	return u, nil
}
