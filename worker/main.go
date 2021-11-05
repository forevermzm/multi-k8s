package main

import (
	"context"
	"fmt"
	"strconv"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func fib(num int64) int64 {
	if num <= 0 {
		return 0
	} else if num == 1 || num == 2 {
		return 1
	} else {
		return fib(num-1) + fib(num-2)
	}
}

func main() {
	fmt.Println("Starting app")
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPort = os.Getenv("REDIS_PORT")
	fmt.Println("Using redis at host: " + redisHost + " and port: " + redisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var channelName = "insert"
	pubsub := rdb.Subscribe(ctx, channelName)
	defer pubsub.Close()

	for {
		// ReceiveTimeout is a low level API. Use ReceiveMessage instead.
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			break
		}

		fmt.Println("received", msg.Payload, "from", msg.Channel)
		intVal, err := strconv.ParseInt(msg.Payload, 10, 0)
		var fibValue = fib(intVal)
		fmt.Println("Calculated fib value", fibValue)
		rdb.HSet(ctx, "values", msg.Payload, fibValue)
	}

	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
	// Output: key value
	// key2 does not exist
}
