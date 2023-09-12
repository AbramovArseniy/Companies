package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	"github.com/AbramovArseniy/Companies/internal/storage/postgres/generated/db"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Kafka struct {
	AlertsDBConn   *pgxpool.Conn
	ChangesDBConn  *pgxpool.Conn
	Consumer       sarama.Consumer
	AlertsStorage  db.Querier
	ChangesStorage db.Querier
}

func New(dbPool *pgxpool.Pool, cfg cfg.Config) (*Kafka, error) {
	alertsDBConn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while acquiring database connection: %w", err)
	}
	changesDBConn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while acquiring database connection: %w", err)
	}
	changesStorage := db.New(changesDBConn)
	consumer, err := sarama.NewConsumer(cfg.Brokers, nil)
	if err != nil {
		return nil, err
	}
	alertsStorage := db.New(alertsDBConn)
	return &Kafka{
		Consumer:       consumer,
		AlertsDBConn:   alertsDBConn,
		ChangesDBConn:  changesDBConn,
		AlertsStorage:  alertsStorage,
		ChangesStorage: changesStorage,
	}, nil
}

func (k *Kafka) ListenAlerts(AlertsTopic string) error {
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
				err = k.HandleAlert(msg.Value)
				if err != nil {
					log.Println(err)
				}
			}
		}(pc)
	}
	return nil
}

func (k *Kafka) HandleAlert(jsonInfo []byte) error {
	var alert Alert
	err := json.Unmarshal(jsonInfo, &alert)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json: %v", err)
	}
	err = k.AlertsStorage.SaveAlert(context.Background(),
		db.SaveAlertParams{
			Type:     alert.Type,
			Uuid:     alert.TagID,
			Column3:  alert.TimeStamp,
			Severity: alert.Severity,
			State:    alert.State,
		})
	if err != nil {
		return fmt.Errorf("cannot update data in database: %v", err)
	}
	return nil
}

func (k *Kafka) ListenTagChanges(ChangesTopic string) error {
	var err error = nil
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
				err = k.SaveTagValue(msg.Value)
				if err != nil {
					log.Println(err)
				}
			}
		}(pc)
	}
	return err
}

func (k *Kafka) SaveTagValue(jsonInfo []byte) error {
	var tagChange TagChange
	err := json.Unmarshal(jsonInfo, &tagChange)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json: %v", err)
	}
	err = k.ChangesStorage.UpdateTag(context.Background(), db.UpdateTagParams{Uuid: tagChange.UUID, Value: tagChange.Value})
	if err != nil {
		return fmt.Errorf("cannot update data in database: %v", err)
	}
	err = k.ChangesStorage.SaveChange(context.Background(), db.SaveChangeParams{Uuid: tagChange.UUID, Column2: tagChange.TimeStamp})
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
	k.AlertsDBConn.Release()
	k.ChangesDBConn.Release()
	return nil
}
