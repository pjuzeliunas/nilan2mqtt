package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pjuzeliunas/nilan2mqtt/internal/adapter"
)

func main() {
	adapter := adapter.NewNilanMQTTAdapter(
		"192.168.1.31:502",  // Nilan address
		"192.168.1.18:1883", // MQTT broker address
		"", "",
	)
	adapter.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
