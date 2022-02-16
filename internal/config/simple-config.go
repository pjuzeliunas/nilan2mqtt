package config

type SimpleConfig struct {
	DeviceClass       string `json:"device_class"`
	Name              string `json:"name"`
	StateTopic        string `json:"state_topic"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	ValueTemplate     string `json:"value_template"`
	UniqueId          string `json:"unique_id"`
}

func RoomTemperature() SimpleConfig {
	return SimpleConfig{
		DeviceClass:       "temperature",
		Name:              "NILAN Room temperature",
		StateTopic:        "nilan/readings/state",
		UnitOfMeasurement: "°C",
		ValueTemplate:     "{{ value_json.room_temp }}",
		UniqueId:          "e8caf088-f8df-4721-942c-aa2f1a1c6f9f",
	}
}

func OutdoorTemperature() SimpleConfig {
	return SimpleConfig{
		DeviceClass:       "temperature",
		Name:              "NILAN Outdoor temperature",
		StateTopic:        "nilan/readings/state",
		UnitOfMeasurement: "°C",
		ValueTemplate:     "{{ value_json.out_temp }}",
		UniqueId:          "54836e31-25ab-4b1e-834a-168403a7658f",
	}
}

func Humidity() SimpleConfig {
	return SimpleConfig{
		DeviceClass:       "humidity",
		Name:              "NILAN Room humidity (actual)",
		StateTopic:        "nilan/readings/state",
		UnitOfMeasurement: "%",
		ValueTemplate:     "{{ value_json.humidity }}",
		UniqueId:          "0fce7936-e340-4fe0-9609-640ce6635d12",
	}
}

func HumidityAvg() SimpleConfig {
	return SimpleConfig{
		DeviceClass:       "humidity",
		Name:              "NILAN Room humidity (average)",
		StateTopic:        "nilan/readings/state",
		UnitOfMeasurement: "%",
		ValueTemplate:     "{{ value_json.humidity_avg }}",
		UniqueId:          "eb87a9cc-fe1e-4694-b86a-cd23fd156f95",
	}
}

func DHWTemperatureTop() SimpleConfig {
	return SimpleConfig{
		DeviceClass:       "temperature",
		Name:              "NILAN DHW tank temperature (top)",
		StateTopic:        "nilan/readings/state",
		UnitOfMeasurement: "°C",
		ValueTemplate:     "{{ value_json.dhw_top_temp }}",
		UniqueId:          "00b07bfe-8e71-459c-baab-fea155c99d12",
	}
}

func DHWTemperatureBottom() SimpleConfig {
	return SimpleConfig{
		DeviceClass:       "temperature",
		Name:              "NILAN DHW tank temperature (bottom)",
		StateTopic:        "nilan/readings/state",
		UnitOfMeasurement: "°C",
		ValueTemplate:     "{{ value_json.dhw_bottom_temp }}",
		UniqueId:          "64e52f1f-2896-489b-8ba6-4dad3ac6a767",
	}
}

func SupplyFlowTemperature() SimpleConfig {
	return SimpleConfig{
		DeviceClass:       "temperature",
		Name:              "NILAN Supply flow temperature",
		StateTopic:        "nilan/readings/state",
		UnitOfMeasurement: "°C",
		ValueTemplate:     "{{ value_json.supply_temp }}",
		UniqueId:          "c510308e-d465-4c96-af9d-74f299a55266",
	}
}
