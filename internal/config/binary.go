package config

import "github.com/pjuzeliunas/nilan2mqtt/internal/dto"

type BinarySensor struct {
	DeviceClass   string `json:"device_class"`
	Device        Device `json:"device"`
	Name          string `json:"name"`
	StateTopic    string `json:"state_topic"`
	UniqueId      string `json:"unique_id"`
	ValueTemplate string `json:"value_template"`
}

func OldFilterSensor() BinarySensor {
	return BinarySensor{
		DeviceClass:   "problem",
		Device:        NilanDevice(),
		Name:          "NILAN Filter status",
		StateTopic:    dto.ErrorsTopic,
		UniqueId:      "fece20c3-c48a-4c10-afa8-f981c8fc31ac",
		ValueTemplate: "{{ value_json.old_filter }}",
	}
}

func ErrorSensor() BinarySensor {
	return BinarySensor{
		DeviceClass:   "problem",
		Device:        NilanDevice(),
		Name:          "NILAN Error status",
		StateTopic:    dto.ErrorsTopic,
		UniqueId:      "651ab1c7-a553-4088-aa37-ac226fed38e7",
		ValueTemplate: "{{  value_json.other_errors }}",
	}
}
