package main

import (
	"dapurcoklat/config"
	"dapurcoklat/domain"
	"dapurcoklat/service_b/kafka/consumer"
	"dapurcoklat/service_b/kafka/producer"
	"encoding/json"
	"os"

	redis_conn "dapurcoklat/service_b/redis"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

func main() {

	redisConn := redis_conn.NewRedisConn()

	// ==================== set users to redis =====================
	user1 := domain.User{
		ID:   1,
		Name: "ihsan",
	}
	user1Json, _ := json.Marshal(user1)

	user2 := domain.User{
		ID:   2,
		Name: "Tono",
	}
	user2Json, _ := json.Marshal(user2)

	user3 := domain.User{
		ID:   3,
		Name: "Yadi",
	}
	user3Json, _ := json.Marshal(user3)

	err := redisConn.Set("user:1", string(user1Json), 0)
	if err != nil {
		panic(err)
	}
	err = redisConn.Set("user:2", string(user2Json), 0)
	if err != nil {
		panic(err)
	}
	err = redisConn.Set("user:3", string(user3Json), 0)
	if err != nil {
		panic(err)
	}
	// =============================================================

	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	kafkaConfig := config.GetKafkaConfig("", "")
	producers, err := sarama.NewSyncProducer([]string{"kafka:9092"}, kafkaConfig)
	if err != nil {
		logrus.Errorf("Unable to create kafka producer got error %v", err)
		return
	}
	defer func() {
		if err := producers.Close(); err != nil {
			logrus.Errorf("Unable to stop kafka producer: %v", err)
			return
		}
	}()

	logrus.Infof("Success create kafka sync-producer")

	kafka := &producer.KafkaProducer{
		Producer: producers,
	}

	consumers, err := sarama.NewConsumer([]string{"kafka:9092"}, kafkaConfig)
	if err != nil {
		logrus.Errorf("Error create kakfa consumer got error %v", err)
	}
	defer func() {
		if err := consumers.Close(); err != nil {
			logrus.Fatal(err)
			return
		}
	}()

	kafkaConsumer := &consumer.KafkaConsumer{
		Consumer:  consumers,
		Producter: kafka,
		Rdb:       redisConn,
	}

	signals := make(chan os.Signal, 1)
	kafkaConsumer.Consume([]string{"service_b"}, signals)

}
