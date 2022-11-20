package main

import (
	"dapurcoklat/config"
	"dapurcoklat/domain"
	"dapurcoklat/service_a/kafka/consumer"
	"dapurcoklat/service_a/kafka/producer"
	"dapurcoklat/service_a/ws"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}

func main() {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	hub := ws.NewHub()

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
		Consumer: consumers,
		WsHub:    hub,
	}

	signals := make(chan os.Signal, 1)
	go kafkaConsumer.Consume([]string{"service_a"}, signals)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}

		hub.Join(conn)

		defer func() {

			hub.Leave(conn)

			conn.Close()
		}()

		// Continuosly read and write message
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {

				log.Println("read failed:", err)
				break
			}

			if err != nil {
				log.Println("json failed:", err)
				break
			}

			req := domain.WebRequest{}
			res := domain.WebResponse{}
			fmt.Println(string(msg))
			err = json.Unmarshal(msg, &req)
			if err != nil {

				log.Println("json failed:", err)

				res.Message = "Unable to parse request"
				res.Success = false
				res.User = domain.User{}

				data, err := json.Marshal(res)
				if err != nil {
					log.Println("json failed:", err)
				}

				conn.WriteMessage(websocket.TextMessage, data)
				continue
			}

			kafka.SendMessage("service_b", string(msg))

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
