package consumer

import (
	"dapurcoklat/domain"
	"dapurcoklat/service_b/kafka/producer"
	"encoding/json"
	"fmt"
	"os"

	redis_conn "dapurcoklat/service_b/redis"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// KafkaConsumer hold sarama consumer
type KafkaConsumer struct {
	Consumer  sarama.Consumer
	Producter *producer.KafkaProducer
	Rdb       *redis_conn.RedisConn
}

// Consume function to consume message from apache kafka
func (c *KafkaConsumer) Consume(topics []string, signals chan os.Signal) {
	chanMessage := make(chan *sarama.ConsumerMessage, 256)

	for _, topic := range topics {
		partitionList, err := c.Consumer.Partitions(topic)
		if err != nil {
			logrus.Errorf("Unable to get partition got error %v", err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(c.Consumer, topic, partition, chanMessage)
		}
	}
	logrus.Infof("Kafka is consuming....")

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			logrus.Infof("New Message from kafka, message: %v", string(msg.Value))

			req := domain.WebRequest{}
			res := domain.WebResponse{}
			err := json.Unmarshal(msg.Value, &req)
			if err != nil {
				logrus.Errorf("Unable to unmarshal message got error %v", err)

				res.Message = "Unable to unmarshal message"
				res.Success = false

				data, err := json.Marshal(res)
				if err != nil {
					logrus.Errorf("Unable to marshal response got error %v", err)
				}

				c.Producter.SendMessage("service_a", string(data))
				continue
			}

			// check if redis is available
			val, err := c.Rdb.Get(fmt.Sprintf("user:%v", req.ID))
			if err == redis.Nil {
				res.Success = false
				res.Message = "Data Not Found"

				data, err := json.Marshal(res)
				if err != nil {
					logrus.Errorf("Unable to marshal response got error %v", err)
				}

				c.Producter.SendMessage("service_a", string(data))
				continue
			} else if err != nil {
				res.Success = false
				res.Message = "Unable to connect to redis"

				data, err := json.Marshal(res)
				if err != nil {
					logrus.Errorf("Unable to marshal response got error %v", err)
				}

				c.Producter.SendMessage("service_a", string(data))
				continue
			}

			user := domain.User{}

			err = json.Unmarshal([]byte(val), &user)
			if err != nil {
				logrus.Errorf("Unable to unmarshal user got error %v", err)

				res.Success = false
				res.Message = "Unable to unmarshal user"

				data, err := json.Marshal(res)
				if err != nil {
					logrus.Errorf("Unable to marshal response got error %v", err)
				}

				c.Producter.SendMessage("service_a", string(data))
				continue
			}

			res.User = user
			res.Success = true
			res.Message = "Success"

			data, err := json.Marshal(res)
			if err != nil {
				logrus.Errorf("Unable to marshal response got error %v", err)
				continue
			}

			c.Producter.SendMessage("service_a", string(data))
		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		logrus.Errorf("Unable to consume partition %v got error %v", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			logrus.Errorf("Unable to close partition %v: %v", partition, err)
		}

		// check if panic
		if r := recover(); r != nil {
			logrus.Errorf("Recovered from panic: %v", r)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}

}
