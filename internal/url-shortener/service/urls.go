package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Alzoww/url-shortener/config"
	"github.com/Alzoww/url-shortener/internal/url-shortener/storage"
	"github.com/segmentio/kafka-go"
	"time"
)

type UrlService struct {
	storage storage.Interface
	kafkaW  *kafka.Writer
}

type msgToSend struct {
	savedAT time.Time
	alias   string
	url     string
}

func NewUrlService(storage storage.Interface, kafkaCfg config.KafkaConfig) *UrlService {
	return &UrlService{
		storage: storage,
		kafkaW: &kafka.Writer{
			Addr:                   kafka.TCP("localhost:29092"),
			Topic:                  "ush.events.local",
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (s *UrlService) URLSave(urlToSave, alias string) error {

	_, err := json.Marshal(msgToSend{
		savedAT: time.Now(),
		alias:   alias,
		url:     urlToSave,
	})
	if err != nil {
		return fmt.Errorf("failed marshal: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(alias),
		Value: []byte(urlToSave),
	}

	err = s.kafkaW.WriteMessages(context.Background(), kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to write kafka message: %w", err)
	}

	return s.storage.SaveURL(urlToSave, alias)
}

func (s *UrlService) URLGet(alias string) (string, error) {
	return s.storage.GetURL(alias)
}
