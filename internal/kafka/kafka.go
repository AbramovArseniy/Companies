package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	"github.com/AbramovArseniy/Companies/internal/storage/postgres/generated/db"
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

/*func (k *Kafka) ListenAlerts(AlertsTopic string) error {
	partitionList, err := k.Consumer.Partitions(AlertsTopic)
	if err != nil {
		return err
	}
	initialOffset := sarama.OffsetOldest
	for _, partition := range partitionList {
		pc, err := k.Consumer.ConsumePartition(AlertsTopic, partition, initialOffset)
		if err != nil {
			return err
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				k.SaveTagValue(msg.Value)
			}
		}(pc)
	}
	return nil
}*/

/*
	func (k *Kafka) HandlerAlert(jsonInfo []byte) error {
		var alert Alert
		err := json.Unmarshal(jsonInfo, &alert)
		if err != nil {
			return fmt.Errorf("cannot unmarshal json: %v", err)
		}
		err = k.Storage.UpdateTag(context.Background(), db.UpdateTagParams{Uuid: TChange.UUID, Value: TChange.Value})
		if err != nil {
			return fmt.Errorf("cannot update data in database: %v", err)
		}
		err = k.Storage.SaveChange(context.Background(), db.SaveChangeParams{Uuid: TChange.UUID, Column2: TChange.TimeStamp})
		if err != nil {
			return fmt.Errorf("cannot update data in database: %v", err)
		}
		return nil
	}
*/
func (k *Kafka) ListenTagChanges(ChangesTopic string) error {
	partitionList, err := k.Consumer.Partitions(ChangesTopic)
	if err != nil {
		return err
	}
	initialOffset := sarama.OffsetOldest
	for _, partition := range partitionList {
		pc, err := k.Consumer.ConsumePartition(ChangesTopic, partition, initialOffset)
		if err != nil {
			return err
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				k.SaveTagValue(msg.Value)
			}
		}(pc)
	}
	return nil
}

func (k *Kafka) SaveTagValue(jsonInfo []byte) error {
	var tagChange TagChange
	err := json.Unmarshal(jsonInfo, &tagChange)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json: %v", err)
	}
	err = k.Storage.UpdateTag(context.Background(), db.UpdateTagParams{Uuid: tagChange.UUID, Value: tagChange.Value})
	if err != nil {
		return fmt.Errorf("cannot update data in database: %v", err)
	}
	err = k.Storage.SaveChange(context.Background(), db.SaveChangeParams{Uuid: tagChange.UUID, Column2: tagChange.TimeStamp})
	if err != nil {
		return fmt.Errorf("cannot update data in database: %v", err)
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
