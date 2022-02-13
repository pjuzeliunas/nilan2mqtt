package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/pjuzeliunas/nilan"
	"github.com/pjuzeliunas/nilan2mqtt/internal"
	"github.com/pjuzeliunas/nilan2mqtt/internal/config"
)

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

func sendSimpleConfig(client mqtt.Client, topic string, config config.SimpleConfig) {
	d, _ := json.Marshal(config)
	t := client.Publish(topic, 0, false, d)
	t.Wait()
}

func sendFanConfig(client mqtt.Client, topic string, config config.Fan) {
	d, _ := json.Marshal(config)
	t := client.Publish(topic, 0, false, d)
	t.Wait()
}

func setUpConfig(client mqtt.Client) {
	sendSimpleConfig(client, "homeassistant/sensor/nilan/1/config", config.RoomTemperature())
	sendSimpleConfig(client, "homeassistant/sensor/nilan/2/config", config.OutdoorTemperature())
	sendSimpleConfig(client, "homeassistant/sensor/nilan/3/config", config.HumidityAvg())
	sendSimpleConfig(client, "homeassistant/sensor/nilan/4/config", config.Humidity())
	sendSimpleConfig(client, "homeassistant/sensor/nilan/5/config", config.DHWTemperatureTop())
	sendSimpleConfig(client, "homeassistant/sensor/nilan/6/config", config.DHWTemperatureBottom())
	sendSimpleConfig(client, "homeassistant/sensor/nilan/7/config", config.SupplyFlowTemperature())
	sendFanConfig(client, "homeassistant/fan/nilan/config", config.NilanVentilation())
}

func publishReadings(client mqtt.Client, readings internal.ReadingsDTO) {
	d, _ := json.Marshal(readings)
	t := client.Publish("homeassistant/sensor/nilan/state", 0, false, d)
	t.Wait()
}

func publishVentilationState(client mqtt.Client, ventilationState internal.VentilationDTO) {
	d, _ := json.Marshal(ventilationState)
	t := client.Publish("homeassistant/fan/nilan/state", 0, false, d)
	t.Wait()
}

func main() {
	c := nilan.Controller{Config: nilan.Config{NilanAddress: "192.168.1.31:502"}}

	mqttC := mqttClient("192.168.1.18", 1883, "", "")
	defer mqttC.Disconnect(0)

	setUpConfig(mqttC)

	for {
		readings := c.FetchReadings()
		readingsDTO := internal.CreateReadingsDTO(readings)
		publishReadings(mqttC, readingsDTO)

		settings := c.FetchSettings()
		ventilationDTO := internal.CreateVentilationDTO(settings)
		publishVentilationState(mqttC, ventilationDTO)

		time.Sleep(time.Minute)
	}
}
