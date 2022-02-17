package adapter

import (
	"encoding/json"
	"fmt"
	"log"
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
	errorsChan   chan nilan.Errors

	running bool
	// update frequency in seconds, default: 15
	updateFrequency int
}

func NewNilanMQTTAdapter(nilanAddress string, mqttBrokerAddress string, mqttUsername string, mqttPassword string) NilanMQTTAdapter {
	a := NilanMQTTAdapter{}

	a.setUpController(nilanAddress)
	a.setUpMQTTClient(mqttBrokerAddress, mqttUsername, mqttPassword)
	a.running = false
	a.updateFrequency = 15

	return a
}

func (a *NilanMQTTAdapter) Start() {
	a.running = true
	a.tryConnectToMQTT(12)
	log.Default().Println("connection to MQTT broker established")
	log.Default().Println("sending HA configuration via MQTT")
	a.sendAllConfigs()

	a.readingsChan = make(chan nilan.Readings)
	a.settingsChan = make(chan nilan.Settings)
	a.errorsChan = make(chan nilan.Errors)

	go a.startFetchingNilanData()
	go a.startPublishingReadings()
	go a.startPublishingSettings()
	go a.startPublishingErrors()

	a.subscribeForTopics()
	log.Default().Println("nilan2mqtt is running")
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
		"nilan/fan/set",
		"nilan/fan/speed/set",
		"nilan/fan/mode/set",
		"nilan/dhw/set",
		"nilan/heating/set",
		"nilan/room_temp/set",
		"nilan/dhw/temp/set",
		"nilan/supply/set",
	}
	for _, t := range topics {
		token := a.mqttClient.Subscribe(t, 1, a.processMessage)
		token.Wait()
	}
}

func (a *NilanMQTTAdapter) processMessage(client mqtt.Client, msg mqtt.Message) {
	log.Default().Printf("received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	payload := string(msg.Payload())
	switch msg.Topic() {
	case "nilan/fan/set":
		settings := nilan.Settings{
			VentilationOnPause: config.OnOffVal(payload),
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "nilan/fan/speed/set":
		speed, _ := strconv.Atoi(payload)
		settings := nilan.Settings{}
		if speed == 0 {
			settings.VentilationOnPause = config.BoolAddr(true)
		} else {
			settings.VentilationOnPause = config.BoolAddr(false)
			settings.FanSpeed = dto.FanSpeed(speed)
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "nilan/fan/mode/set":
		settings := nilan.Settings{
			VentilationMode: dto.Mode(payload),
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "nilan/dhw/set":
		settings := nilan.Settings{
			DHWProductionPaused: config.BoolAddr(payload == "OFF"),
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "nilan/heating/set":
		settings := nilan.Settings{
			CentralHeatingPaused: config.BoolAddr(payload == "OFF"),
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "nilan/room_temp/set":
		temp, _ := strconv.Atoi(payload)
		temp *= 10
		settings := nilan.Settings{
			DesiredRoomTemperature: &temp,
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "nilan/dhw/temp/set":
		temp, _ := strconv.Atoi(payload)
		temp *= 10
		settings := nilan.Settings{
			DesiredDHWTemperature: &temp,
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	case "nilan/supply/set":
		temp, _ := strconv.Atoi(payload)
		temp *= 10
		settings := nilan.Settings{
			SetpointSupplyTemperature: &temp,
		}
		a.nilanController.SendSettings(settings)
		a.fetchSettings()
	}
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

func (a *NilanMQTTAdapter) sendConfig(topic string, config interface{}) {
	d, _ := json.Marshal(config)
	t := a.mqttClient.Publish(topic, 0, false, d)
	t.Wait()
}

func (a *NilanMQTTAdapter) sendAllConfigs() {
	a.sendConfig("homeassistant/sensor/nilan/1/config", config.RoomTemperature())
	a.sendConfig("homeassistant/sensor/nilan/2/config", config.OutdoorTemperature())
	a.sendConfig("homeassistant/sensor/nilan/3/config", config.HumidityAvg())
	a.sendConfig("homeassistant/sensor/nilan/4/config", config.Humidity())
	a.sendConfig("homeassistant/sensor/nilan/5/config", config.DHWTemperatureTop())
	a.sendConfig("homeassistant/sensor/nilan/6/config", config.DHWTemperatureBottom())
	a.sendConfig("homeassistant/sensor/nilan/7/config", config.SupplyFlowTemperature())
	a.sendConfig("homeassistant/fan/nilan/config", config.NilanVentilation())
	a.sendConfig("homeassistant/switch/nilan/1/config", config.DHWSwitch())
	a.sendConfig("homeassistant/switch/nilan/2/config", config.CentralHeatingSwitch())
	a.sendConfig("homeassistant/number/nilan/1/config", config.RoomTemperatureSetpoint())
	a.sendConfig("homeassistant/number/nilan/2/config", config.DHWTemperatureSetpoint())
	a.sendConfig("homeassistant/number/nilan/3/config", config.SupplyFlowSetpoint())
	a.sendConfig("homeassistant/binary_sensor/nilan/1/config", config.OldFilterSensor())
	a.sendConfig("homeassistant/binary_sensor/nilan/2/config", config.ErrorSensor())
}

func (a *NilanMQTTAdapter) startFetchingNilanData() {
	for a.running {
		a.fetchReadings()
		a.fetchSettings()
		a.fetchErrors()
		time.Sleep(time.Second * time.Duration(a.updateFrequency))
	}
	close(a.readingsChan)
}

func (a *NilanMQTTAdapter) fetchReadings() {
	readings, err := a.nilanController.FetchReadings()
	if err != nil {
		log.Default().Printf("error (fetch readings) - %s\n", err)
		return
	}
	a.readingsChan <- *readings
}

func (a *NilanMQTTAdapter) fetchSettings() {
	settings, err := a.nilanController.FetchSettings()
	if err != nil {
		log.Default().Printf("error (fetch settings) - %s\n", err)
		return
	}
	a.settingsChan <- *settings
}

func (a *NilanMQTTAdapter) fetchErrors() {
	errors, err := a.nilanController.FetchErrors()
	if err != nil {
		log.Default().Printf("error (fetch errors) - %s\n", err)
		return
	}
	a.errorsChan <- *errors
}

func (a *NilanMQTTAdapter) startPublishingReadings() {
	for readings := range a.readingsChan {
		readingsDTO := dto.CreateReadingsDTO(readings)
		a.publish("nilan/readings", readingsDTO)
	}
}

func (a *NilanMQTTAdapter) startPublishingSettings() {
	for settings := range a.settingsChan {
		settingsDTO := dto.CreateSettingsDTO(settings)
		a.publish("nilan/settings", settingsDTO)
	}
}

func (a *NilanMQTTAdapter) startPublishingErrors() {
	for errors := range a.errorsChan {
		errorsDTO := dto.CreateErrorsDTO(errors)
		a.publish("nilan/errors", errorsDTO)
	}
}

func (a *NilanMQTTAdapter) publish(topic string, v interface{}) {
	d, _ := json.Marshal(v)
	t := a.mqttClient.Publish(topic, 0, false, d)
	t.Wait()
}
