package config

type Fan struct {
	Name                         string   `json:"name"`
	CommandTopic                 string   `json:"command_topic"`
	PercentageStateTopic         string   `json:"percentage_state_topic"`
	PercentageStateValueTemplate string   `json:"percentage_value_template"`
	PercentageCommandTopic       string   `json:"percentage_command_topic"`
	PresetModeStateTopic         string   `json:"preset_mode_state_topic"`
	PresetModeValueTemplate      string   `json:"preset_mode_value_template"`
	PresetModeCommandTopic       string   `json:"preset_mode_command_topic"`
	PresetModes                  []string `json:"preset_modes"`
	SpeedRangeMin                int      `json:"speed_range_min"`
	SpeedRangeMax                int      `json:"speed_range_max"`
	StateTopic                   string   `json:"state_topic"`
	StateValueTemplate           string   `json:"state_value_template"`
	UniqueId                     string   `json:"unique_id"`
}

func NilanVentilation() Fan {
	return Fan{
		Name:                         "NILAN Ventilation",
		CommandTopic:                 "homeassistant/fan/nilan/set",
		PercentageStateTopic:         "homeassistant/fan/nilan/state",
		PercentageStateValueTemplate: "{{ value_json.speed }}",
		PercentageCommandTopic:       "homeassistant/fan/nilan/speed/set",
		PresetModeStateTopic:         "homeassistant/fan/nilan/state",
		PresetModeValueTemplate:      "{{ value_json.mode }}",
		PresetModeCommandTopic:       "homeassistant/fan/nilan/mode/set",
		PresetModes:                  []string{"auto", "heating", "cooling"},
		SpeedRangeMin:                1,
		SpeedRangeMax:                4,
		StateTopic:                   "homeassistant/fan/nilan/state",
		StateValueTemplate:           "{{ value_json.state }}",
		UniqueId:                     "3d5c2bc2-a192-4c4a-9171-a23b4ba6c16c",
	}
}
