package config

import "github.com/pjuzeliunas/nilan2mqtt/internal/dto"

type Switch struct {
	CommandTopic       string `json:"command_topic"`
	Device             Device `json:"device"`
	Icon               string `json:"icon"`
	Name               string `json:"name"`
	StateTopic         string `json:"state_topic"`
	StateValueTemplate string `json:"value_template"`
	UniqueId           string `json:"unique_id"`
}

func DHWSwitch() Switch {
	return Switch{
		CommandTopic:       "nilan/dhw/set",
		Device:             NilanDevice(),
		Icon:               "mdi:water-pump",
		Name:               "NILAN Domestic Hot Water",
		StateTopic:         dto.SettingsTopic,
		StateValueTemplate: "{{ value_json.dhw_state }}",
		UniqueId:           "b7e46b2f-5575-40a1-a690-233d55cd33bf",
	}
}

func CentralHeatingSwitch() Switch {
	return Switch{
		CommandTopic:       "nilan/heating/set",
		Device:             NilanDevice(),
		Icon:               "mdi:radiator",
		Name:               "NILAN Central Heating",
		StateTopic:         dto.SettingsTopic,
		StateValueTemplate: "{{ value_json.central_heating_state }}",
		UniqueId:           "fb1e0156-9250-4c78-a8c0-44113132062c",
	}
}

func OnOffVal(str string) *bool {
	switch str {
	case "ON":
		return BoolAddr(true)
	case "OFF":
		return BoolAddr(false)
	default:
		return nil
	}
}

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}
