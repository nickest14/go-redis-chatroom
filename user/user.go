package user

import (
	"context"
	"encoding/json"
	"fmt"
	"go-redis-chatroom/constants"
	rediswrap "go-redis-chatroom/redis"

	"github.com/go-redis/redis/v8"
)

type receiveMsg struct {
	Content string `json:"content,omitempty"`
	Channel string `json:"channel,omitempty"`
	Sender  string `json:"sender,omitempty"`
}

// User struct
type User struct {
	name             string
	channelsHandler  *redis.PubSub
	stopListenerChan chan struct{}
	listening        bool
	MessageChan      chan receiveMsg
	isVisitor        bool
}

//Connect connect user to user channels on redis
func Connect(name string, isVisitor bool) (*User, error) {
	u := &User{
		name:             name,
		stopListenerChan: make(chan struct{}),
		MessageChan:      make(chan receiveMsg),
		isVisitor:        isVisitor,
	}
	if err := u.channelConnect(); err != nil {
		return nil, err
	}

	return u, nil
}

// Disconnect the user channels
func (u *User) Disconnect() error {
	if u.channelsHandler != nil {
		// if err := u.channelsHandler.Unsubscribe(ctx); err != nil {
		// 	return err
		// }
		if err := u.channelsHandler.Close(); err != nil {
			return err
		}
	}
	if u.listening {
		u.stopListenerChan <- struct{}{}
	}

	close(u.MessageChan)

	return nil
}

func (u *User) channelConnect() error {
	ctx := context.Background()
	rdb := rediswrap.Client

	if !u.isVisitor {
		if _, err := rdb.SAdd(ctx, constants.UsersKey, u.name).Result(); err != nil {
			return err
		}
	}
	userChannels, err := GetChannels(u.name)
	if err != nil {
		return err
	}

	if len(userChannels) == 0 {
		fmt.Println("no channels to connect to for user: ", u.name)
		return nil
	}

	return u.doConnect(ctx, rdb, userChannels...)
}

// GetChannels return the user channels
func GetChannels(username string) ([]string, error) {
	ctx := context.Background()
	rdb := rediswrap.Client

	var c []string
	c1, err := rdb.SMembers(ctx, constants.ChannelsKey).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c1...)

	// get all user channels from userchannels key
	c2, err := rdb.SMembers(ctx, fmt.Sprintf(constants.UserChannels, username)).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c2...)

	return c, nil
}

func (u *User) doConnect(ctx context.Context, rdb *redis.Client, channels ...string) error {
	// subscribe all channels in one request
	pubSub := rdb.Subscribe(ctx, channels...)
	// keep channel handler to be used in unsubscribe
	u.channelsHandler = pubSub
	// The Listener
	go func() {
		u.listening = true
		fmt.Println("starting the listener for user:", u.name, "on channels:", channels)
		for {
			select {
			case msg, ok := <-pubSub.Channel():
				if !ok {
					fmt.Println(u.name, "has disconnected")
					return
				}
				var msgMap map[string]string
				if err := json.Unmarshal([]byte(msg.Payload), &msgMap); err == nil {
					msgDetail := receiveMsg{
						Sender:  msgMap["sender"],
						Channel: msgMap["channel"],
						Content: msgMap["content"],
					}
					u.MessageChan <- msgDetail
				}
			case <-u.stopListenerChan:
				fmt.Println("stopping the listener for user:", u.name)
				return
			}
		}
	}()
	return nil
}

// HandleSubscribe handle subscribe and unsubscribe the specific channel
func (u *User) HandleSubscribe(subscrube bool, channel string) error {
	ctx := context.Background()
	rdb := rediswrap.Client
	userChannelsKey := fmt.Sprintf(constants.UserChannels, u.name)

	if subscrube == true {
		// handle subscribe
		if rdb.SIsMember(ctx, userChannelsKey, channel).Val() {
			return nil
		}
		if err := rdb.SAdd(ctx, userChannelsKey, channel).Err(); err != nil {
			return err
		}
		u.channelsHandler.Subscribe(ctx, channel)
	} else {
		// handle unsubscribe
		if !rdb.SIsMember(ctx, userChannelsKey, channel).Val() {
			return nil
		}
		if err := rdb.SRem(ctx, userChannelsKey, channel).Err(); err != nil {
			return err
		}
		u.channelsHandler.Unsubscribe(ctx, channel)
	}
	return nil
}

// Chat function send message to users
func Chat(channel string, content string, u *User) error {
	ctx := context.Background()
	rdb := rediswrap.Client
	var sender string
	if u == nil {
		sender = "system"
	} else {
		sender = u.name
	}
	sendMessage, _ := json.Marshal(map[string]string{
		"channel": channel,
		"sender":  sender,
		"content": content,
	})
	return rdb.Publish(ctx, channel, sendMessage).Err()
}
