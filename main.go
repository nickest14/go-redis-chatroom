package main

import (
	"context"
	"fmt"
	rediswrap "go-redis-chatroom/redis"
)

func main() {
	ctx := context.Background()
	rdb := rediswrap.Client
	fmt.Println(rdb.Get(ctx, "test"))
}
