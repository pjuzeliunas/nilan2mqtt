package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func clear(client mqtt.Client, topic string) {
	t := client.Publish(topic, 0, false, "")
	t.Wait()
}

func clearAll(client mqtt.Client) {
	clear(client, "homeassistant/sensor/nilan/1/config")
	clear(client, "homeassistant/sensor/nilan/2/config")
	clear(client, "homeassistant/sensor/nilan/3/config")
	clear(client, "homeassistant/sensor/nilan/4/config")
	clear(client, "homeassistant/sensor/nilan/5/config")
	clear(client, "homeassistant/sensor/nilan/6/config")
	clear(client, "homeassistant/sensor/nilan/7/config")
}

func mqttClient(brokerAddress string, port int, username string, password string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", brokerAddress, port))
	opts.SetClientID(uuid.New().String())
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		panic(err)
	}
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}

func main() {
	mqttC := mqttClient("192.168.1.18", 1883, "", "")
	defer mqttC.Disconnect(0)
	clearAll(mqttC)
}
