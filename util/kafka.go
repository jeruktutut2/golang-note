package util

import (
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	Brokers             = []string{"localhost:9092"}
	Topic   goka.Stream = "example-stream"
	Group   goka.Group  = "example-group"
)

func NewEmmiter(brokers []string, topic goka.Stream) *goka.Emitter {
	emitter, err := goka.NewEmitter(brokers, topic, new(codec.String))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	return emitter
}
