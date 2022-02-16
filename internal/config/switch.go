package config

type Switch struct {
	CommandTopic string `json:"command_topic"`
	Device       Device `json:"device"`
	Icon         string `json:"icon"`
	Name         string `json:"name"`
	StateTopic   string `json:"state_topic"`
	UniqueId     string `json:"unique_id"`
}

func DHWSwitch() Switch {
	return Switch{
		CommandTopic: "nilan/dhw/set",
		Device:       NilanDevice(),
		Icon:         "mdi:water-pump",
		Name:         "NILAN Domestic Hot Water",
		StateTopic:   "nilan/dhw/state",
		UniqueId:     "b7e46b2f-5575-40a1-a690-233d55cd33bf",
	}
}

func CentralHeatingSwitch() Switch {
	return Switch{
		CommandTopic: "nilan/heating/set",
		Device:       NilanDevice(),
		Icon:         "mdi:radiator",
		Name:         "NILAN Central Heating",
		StateTopic:   "nilan/heating/state",
		UniqueId:     "fb1e0156-9250-4c78-a8c0-44113132062c",
	}
}

func OnOffString(on bool) string {
	if on {
		return "ON"
	} else {
		return "OFF"
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
