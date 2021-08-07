package main

import (
	"github.com/go-fed/activity/pub"
	"github.com/go-redis/redis/v8"
)

func main() {
	s := &service{}
	db := &database{
		client: redis.NewClient(&redis.Options{}),
	}
	actor := pub.NewFederatingActor( /* CommonBehavior */ s /*FederatingProtocol*/, s /*Database */, db /*Clock*/, s)
	return
}
