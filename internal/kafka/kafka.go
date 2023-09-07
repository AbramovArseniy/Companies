package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	"github.com/AbramovArseniy/Companies/internal/storage/postgres/db"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Kafka struct {
	Consumer sarama.Consumer
	Storage  db.Querier
}

func New(dbPool *pgxpool.Pool, cfg cfg.Config) (*Kafka, error) {
	dbConn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while acquiring database connection: %w", err)
	}
	storage := db.New(dbConn)
	consumer, err := sarama.NewConsumer(cfg.Brokers, nil)
	if err != nil {
		return nil, err
	}
	return &Kafka{
		Consumer: consumer,
		Storage:  storage,
	}, nil
}

func (k *Kafka) ListenTagChanges(ChangesTopic string, gr *sync.WaitGroup) error {
	partitionList, err := k.Consumer.Partitions(ChangesTopic) //get all partitions
	if err != nil {
		return err
	}
	initialOffset := sarama.OffsetOldest //offset to start reading message from
	for _, partition := range partitionList {
		pc, err := k.Consumer.ConsumePartition(ChangesTopic, partition, initialOffset)
		if err != nil {
			return err
		}
		gr.Add(1)
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				// here is what to do with kafka info
				log.Println(string(msg.Value))
			}
			gr.Done()
		}(pc)
	}
	return nil
}

func (k *Kafka) Close() error {
	err := k.Consumer.Close()
	if err != nil {
		return err
	}
	return nil
}
