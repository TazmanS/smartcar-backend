package mqtt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/TazmanS/smartcar-backend/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	MQTTTopicActions   = "smartcar/actions_topic"
	MQTTTopicSessionId = "smartcar/session_id_topic"
	MQTTTopicHeartbeat = "smartcar/heartbeat_topic"

	MQTTSubSessionId = "smartcar/session_id_sub"
)

type Client struct {
	client mqtt.Client
}

func New(cfg *config.Config) (*Client, error) {
	opts := mqtt.NewClientOptions()

	broker := fmt.Sprintf(
		"tcp://%s:%s",
		cfg.MQTTHost,
		cfg.MQTTPort,
	)

	idBytes := make([]byte, 8)

	if _, err := rand.Read(idBytes); err != nil {
		return nil, err
	}

	clientID := "smartcar-backend-" + hex.EncodeToString(idBytes)

	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected to MQTT Broker")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)

	token := client.Connect()
	token.Wait()

	if err := token.Error(); err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Close() {
	c.client.Disconnect(250)
}

func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, 0, handler)
	token.Wait()

	return token.Error()
}

func (c *Client) Publish(topic string, payload string) error {
	token := c.client.Publish(topic, 0, false, payload)
	token.Wait()

	return token.Error()
}
