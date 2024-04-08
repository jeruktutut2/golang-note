package util

import (
	"fmt"
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	Brokers             = []string{"localhost:9092"}
	Topic   goka.Stream = "example-stream"
	Group   goka.Group  = "example-group"
)

func NewEmmiter() *goka.Emitter {
	emitter, err := goka.NewEmitter(Brokers, Topic, new(codec.String))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	return emitter
}

func NewProcessor() (p *goka.Processor, err error) {
	cb := func(ctx goka.Context, msg interface{}) {
		fmt.Println("key:", ctx.Key(), "message:", msg)
	}

	g := goka.DefineGroup(Group,
		goka.Input(Topic, new(codec.String), cb),
		goka.Persist(new(codec.Int64)),
	)

	p, err = goka.NewProcessor(Brokers, g)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}
	return
}
