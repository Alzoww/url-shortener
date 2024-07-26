package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type AnalyticService struct {
	kafkaR *kafka.Reader
	urls   []string
}

func newAnalyticService() *AnalyticService {
	return &AnalyticService{
		kafkaR: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"localhost:29092"},
			Topic:    "ush.events.local",
			GroupID:  "test-consumer-group",
			MinBytes: 10e1,
			MaxBytes: 10e6,
		}),
		urls: make([]string, 0, 100),
	}
}

func main() {
	an := newAnalyticService()

	ch := make(chan string)

	go an.consume(ch)

	go an.analyze(ch)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	log.Println("analytic service gracefully shutdown...")

	if err := an.kafkaR.Close(); err != nil {
		log.Println("failed to close kafka conn")
		return
	}
}

func (s *AnalyticService) consume(ch chan string) {
	for {
		msg, err := s.kafkaR.FetchMessage(context.Background())
		if err != nil {
			log.Println(err)
		}

		err = s.kafkaR.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Println(err)
		}

		ch <- string(msg.Value)
	}
}

func (s *AnalyticService) analyze(ch chan string) {
	for v := range ch {
		s.urls = append(s.urls, v)

		fmt.Printf("Всего ссылок сохранено: %d\n", len(s.urls))
		fmt.Printf("Средняя длина ссылок: %d\n", avgLenUrls(s.urls))

	}
}

func avgLenUrls(urls []string) int {
	total := 0
	for _, url := range urls {
		total += len(url)
	}

	return total / len(urls)
}

func analyze() {}
