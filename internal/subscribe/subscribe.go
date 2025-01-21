package subscribe

import (
    "encoding/json"
    "fmt"
    "log"
    "os"

    mqtt "github.com/eclipse/paho.mqtt.golang"
		"mqtt-backend-service/internal/types"
		"mqtt-backend-service/internal/database"
)

type Service struct {
	db database.Service
}

func NewService(db database.Service) *Service {
	return &Service{db: db}
}
func (s *Service) messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	if msg.Topic() == os.Getenv("MQTT_TOPIC") {
		var data types.DeviceData
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return
		}

		if data.DeviceID == "" || data.Humidity == 0 || data.Temperature == 0 || data.Timestamp.IsZero() {
			log.Printf("Missing required data fields")
			return
		}

		if err := s.db.InsertDeviceData(data); err != nil {
			log.Printf("Error inserting data into the database: %v", err)
			return
		}

		fmt.Println("Data inserted:", data)
	} else {
		log.Printf("Unknown topic: %s", msg.Topic())
	}
}

func (s *Service) Subscribe(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, s.messagePubHandler)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s: %v", topic, token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT broker")
}

func Start(db database.Service) {
    var broker = os.Getenv("MQTT_BROKER")
    var port = 1883
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", broker, port))
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

		subService := NewService(db)

		topic := os.Getenv("MQTT_TOPIC")
		subService.Subscribe(client, topic)

    select {}
}