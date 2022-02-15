package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pjuzeliunas/nilan2mqtt/internal/adapter"
)

func main() {
	nilanAddr := os.Getenv("NILAN_ADDR")
	mqttAddr := os.Getenv("MQTT_ADDR")
	mqttUser := os.Getenv("MQTT_USER")
	mqttPwd := os.Getenv("MQTT_PWD")

	adapter := adapter.NewNilanMQTTAdapter(
		nilanAddr, // Nilan address
		mqttAddr,  // MQTT broker address
		mqttUser, mqttPwd,
	)
	adapter.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	adapter.Stop()
}
