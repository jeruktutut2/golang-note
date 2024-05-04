package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

func ConsumeEmailKafka(ctx context.Context, brokers []string, topic goka.Stream, consumerGroup goka.Group) {
	cb := func(ctx goka.Context, msg interface{}) {
		fmt.Println("key:", ctx.Key(), "message:", msg)
	}
	newProcessor(ctx, brokers, topic, consumerGroup, cb)
}

func ConsumeTextMessageKafka(ctx context.Context, brokers []string, topic goka.Stream, consumerGroup goka.Group) {
	cb := func(ctx goka.Context, msg interface{}) {
		fmt.Println("key:", ctx.Key(), "message:", msg)
	}
	newProcessor(ctx, brokers, topic, consumerGroup, cb)
}

func newProcessor(ctx context.Context, brokers []string, topic goka.Stream, consumerGroup goka.Group, cb func(ctx goka.Context, msg interface{})) {
	g := goka.DefineGroup(consumerGroup,
		goka.Input(topic, new(codec.String), cb),
		goka.Persist(new(codec.Int64)),
	)

	tmc := goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1

	config := goka.DefaultConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	goka.ReplaceGlobalConfig(config)

	tm, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		log.Fatalf("Error creating topic manager: %v", err)
	}
	defer tm.Close()
	err = tm.EnsureStreamExists(string(topic), 8)
	if err != nil {
		log.Printf("Error creating kafka topic %s: %v", topic, err)
	}

	p, err := goka.NewProcessor(brokers,
		g,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder))
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}
	go func() {
		if err = p.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()
}
