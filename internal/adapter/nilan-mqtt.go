package adapter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/pjuzeliunas/nilan"
	"github.com/pjuzeliunas/nilan2mqtt/internal/config"
	"github.com/pjuzeliunas/nilan2mqtt/internal/dto"
)

type NilanMQTTAdapter struct {
	nilanController nilan.Controller
	mqttClient      mqtt.Client

	readingsChan chan nilan.Readings
	settingsChan chan nilan.Settings

	running bool
}

func NewNilanMQTTAdapter(nilanAddress string, mqttBrokerAddress string, mqttUsername string, mqttPassword string) NilanMQTTAdapter {
	a := NilanMQTTAdapter{}

	a.setUpController(nilanAddress)
	a.setUpMQTTClient(mqttBrokerAddress, mqttUsername, mqttPassword)
	a.running = false

	return a
}

func (a *NilanMQTTAdapter) Start() {
	a.running = true
	a.tryConnectToMQTT(0)
	a.sendConfig()

	a.readingsChan = make(chan nilan.Readings)
	go a.startFetchingReadings()
	go a.startPublishingReadings()

	a.settingsChan = make(chan nilan.Settings)
	go a.startFetchingSettings()
	go a.startPublishingSettings()

	a.subscribeForTopics()
}

func (a *NilanMQTTAdapter) Stop() {
	a.running = false
	a.mqttClient.Disconnect(5000)
}

func (a *NilanMQTTAdapter) setUpController(address string) {
	a.nilanController = nilan.Controller{Config: nilan.Config{NilanAddress: address}}
}

func (a *NilanMQTTAdapter) setUpMQTTClient(address string, username string, password string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", address))
	opts.SetClientID(uuid.New().String())
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.OnConnectionLost = a.reconnect
	a.mqttClient = mqtt.NewClient(opts)
}

func (a *NilanMQTTAdapter) subscribeForTopics() {
	topics := []string{
		"homeassistant/fan/nilan/set",
		"homeassistant/fan/nilan/speed/set",
		"homeassistant/fan/nilan/mode/set",
	}
	for _, t := range topics {
		token := a.mqttClient.Subscribe(t, 1, a.processMessage)
		token.Wait()
	}
}

func (a *NilanMQTTAdapter) processMessage(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	switch msg.Topic() {
	case "homeassistant/fan/nilan/set":
		settings := nilan.Settings{
			VentilationOnPause: boolAddr(payload == "OFF"),
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "homeassistant/fan/nilan/speed/set":
		speed, _ := strconv.Atoi(payload)
		settings := nilan.Settings{}
		if speed == 0 {
			settings.VentilationOnPause = boolAddr(true)
		} else {
			settings.VentilationOnPause = boolAddr(false)
			settings.FanSpeed = dto.FanSpeed(speed)
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "homeassistant/fan/nilan/mode/set":
		settings := nilan.Settings{
			VentilationMode: dto.Mode(payload),
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	}
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func boolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func (a *NilanMQTTAdapter) reconnect(client mqtt.Client, err error) {
	a.tryConnectToMQTT(3)
}

func (a *NilanMQTTAdapter) tryConnectToMQTT(attempts int) {
	if token := a.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		if attempts > 0 {
			time.Sleep(time.Second * 10)
			a.tryConnectToMQTT(attempts - 1)
		} else {
			panic(token.Error())
		}
	}
}

func (a *NilanMQTTAdapter) sendSimpleConfig(topic string, config config.SimpleConfig) {
	d, _ := json.Marshal(config)
	t := a.mqttClient.Publish(topic, 0, false, d)
	t.Wait()
}

func (a *NilanMQTTAdapter) sendFanConfig(topic string, config config.Fan) {
	d, _ := json.Marshal(config)
	t := a.mqttClient.Publish(topic, 0, false, d)
	t.Wait()
}

func (a *NilanMQTTAdapter) sendConfig() {
	a.sendSimpleConfig("homeassistant/sensor/nilan/1/config", config.RoomTemperature())
	a.sendSimpleConfig("homeassistant/sensor/nilan/2/config", config.OutdoorTemperature())
	a.sendSimpleConfig("homeassistant/sensor/nilan/3/config", config.HumidityAvg())
	a.sendSimpleConfig("homeassistant/sensor/nilan/4/config", config.Humidity())
	a.sendSimpleConfig("homeassistant/sensor/nilan/5/config", config.DHWTemperatureTop())
	a.sendSimpleConfig("homeassistant/sensor/nilan/6/config", config.DHWTemperatureBottom())
	a.sendSimpleConfig("homeassistant/sensor/nilan/7/config", config.SupplyFlowTemperature())
	a.sendFanConfig("homeassistant/fan/nilan/config", config.NilanVentilation())
}

func (a *NilanMQTTAdapter) startFetchingReadings() {
	for a.running {
		a.fetchReadings()
		time.Sleep(time.Second * 5)
	}
	close(a.readingsChan)
}

func (a *NilanMQTTAdapter) fetchReadings() {
	readings, err := a.nilanController.FetchReadings()
	if err != nil {
		fmt.Println(err)
		return
	}
	a.readingsChan <- *readings
}

func (a *NilanMQTTAdapter) startFetchingSettings() {
	for a.running {
		a.fetchSettings()
		time.Sleep(time.Second * 5)
	}
	close(a.settingsChan)
}

func (a *NilanMQTTAdapter) fetchSettings() {
	settings, err := a.nilanController.FetchSettings()
	if err != nil {
		fmt.Println(err)
		return
	}
	a.settingsChan <- *settings
}

func (a *NilanMQTTAdapter) startPublishingReadings() {
	for readings := range a.readingsChan {
		readingsDTO := dto.CreateReadingsDTO(readings)
		a.publishReadings(readingsDTO)
	}
}

func (a *NilanMQTTAdapter) publishReadings(readings dto.Readings) {
	d, _ := json.Marshal(readings)
	t := a.mqttClient.Publish("homeassistant/sensor/nilan/state", 0, false, d)
	t.Wait()
}

func (a *NilanMQTTAdapter) startPublishingSettings() {
	for settings := range a.settingsChan {
		ventilationDTO := dto.CreateVentilationDTO(settings)
		a.publishVentilationState(ventilationDTO)
	}
}

func (a *NilanMQTTAdapter) publishVentilationState(ventilationState dto.Ventilation) {
	d, _ := json.Marshal(ventilationState)
	t := a.mqttClient.Publish("homeassistant/fan/nilan/state", 0, false, d)
	t.Wait()
}
