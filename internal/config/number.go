package config

type Number struct {
	CommandTopic      string  `json:"command_topic"`
	Device            Device  `json:"device"`
	Icon              string  `json:"icon"`
	Min               float32 `json:"min"`
	Max               float32 `json:"max"`
	Name              string  `json:"name"`
	StateTopic        string  `json:"state_topic"`
	ValueTemplate     string  `json:"value_template"`
	Step              float32 `json:"step"`
	UnitOfMeasurement string  `json:"unit_of_measurement"`
	UniqueId          string  `json:"unique_id"`
}

func RoomTemperatureSetpoint() Number {
	return Number{
		CommandTopic:      "nilan/room_temp/set",
		Device:            NilanDevice(),
		Icon:              "mdi:hvac",
		Min:               5.0,
		Max:               40.0,
		Name:              "NILAN Room temperature setting",
		StateTopic:        "nilan/settings",
		ValueTemplate:     "{{ value_json.room_temp_setpoint }}",
		Step:              1.0,
		UnitOfMeasurement: "°C",
		UniqueId:          "071347f8-2eef-4e21-9cb0-f35b6dbb3f5b",
	}
}

func DHWTemperatureSetpoint() Number {
	return Number{
		CommandTopic:      "nilan/dhw/temp/set",
		Device:            NilanDevice(),
		Icon:              "mdi:water-pump",
		Min:               10.0,
		Max:               60.0,
		Name:              "NILAN DHW temperature setting",
		StateTopic:        "nilan/settings",
		ValueTemplate:     "{{ value_json.dhw_setpoint }}",
		Step:              1.0,
		UnitOfMeasurement: "°C",
		UniqueId:          "67d522bb-21ff-4c22-ab25-a312150c6132",
	}
}

func SupplyFlowSetpoint() Number {
	return Number{
		CommandTopic:      "nilan/supply/set",
		Device:            NilanDevice(),
		Icon:              "mdi:pipe",
		Min:               5.0,
		Max:               50.0,
		Name:              "NILAN Supply flow temperature setting",
		StateTopic:        "nilan/settings",
		ValueTemplate:     "{{ value_json.supply_flow_setpoint }}",
		Step:              1.0,
		UnitOfMeasurement: "°C",
		UniqueId:          "177b37a2-ca7c-4188-894b-6a797b95de34",
	}
}
