package dto

import "github.com/pjuzeliunas/nilan"

type Readings struct {
	RoomTemperature          float32 `json:"room_temp"`
	OutdoorTemperature       float32 `json:"out_temp"`
	AverageHumidity          int     `json:"humidity_avg"`
	ActualHumidity           int     `json:"humidity"`
	DHWTankTopTemperature    float32 `json:"dhw_top_temp"`
	DHWTankBottomTemperature float32 `json:"dhw_bottom_temp"`
	SupplyFlowTemperature    float32 `json:"supply_temp"`
}

const ReadingsTopic string = "nilan/readings"

func CreateReadingsDTO(readings nilan.Readings) Readings {
	return Readings{
		RoomTemperature:          temperature(readings.RoomTemperature),
		OutdoorTemperature:       temperature(readings.OutdoorTemperature),
		AverageHumidity:          readings.AverageHumidity,
		ActualHumidity:           readings.ActualHumidity,
		DHWTankTopTemperature:    temperature(readings.DHWTankTopTemperature),
		DHWTankBottomTemperature: temperature(readings.DHWTankBottomTemperature),
		SupplyFlowTemperature:    temperature(readings.SupplyFlowTemperature),
	}
}

func temperature(rawValue int) float32 {
	return float32(rawValue) / 10.0
}
