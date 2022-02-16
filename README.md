NOTE: Development is in progress.

**nilan2mqtt** fetches Nilan (R) heatpump CTS700 readings and settings via modbus and sends them to MQTT broker in regular intervals. It also executes commands retrieved from MQTT. Compatible with Home Assistant.

## Environment variables:
- `NILAN_ADDR` - Nilan heatpump address, including port. E.g. `192.168.1.15:502`
- `MQTT_ADDR` - Address of MQTT broker, e.g. `"192.168.1.18:1883"`
- `MQTT_USER` - (optional) MQTT broker username
- `MQTT_PWD` - (optional) MQTT broker password

## Local build and run

```
go get all
go build -o nilan2mqtt cmd/nilan2mqtt/main.go
NILAN_ADDR="<heatpump addr>" MQTT_ADDR="<MQTT broker addr>" ./nilan2mqtt 
```

## Docker build

```
docker build -t pjuzeliunas/nilan2mqtt:dev .
```

## Docker run
```
docker run -e NILAN_ADDR="<heatpump addr>" -e MQTT_ADDR="<MQTT broker addr>" --restart=unless-stopped --net=host --name=nilan2mqtt -d pjuzeliunas/nilan2mqtt:dev
```

Also available in DockerHub: [pjuzeliunas/nilan2mqtt](https://hub.docker.com/repository/docker/pjuzeliunas/nilan2mqtt)
