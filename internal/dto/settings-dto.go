package dto

import (
	"github.com/pjuzeliunas/nilan"
)

type SettingsDTO struct {
	FanState            string `json:"fan_state"`
	Speed               int    `json:"fan_speed"`
	Mode                string `json:"fan_mode"`
	DHWState            string `json:"dhw_state"`
	CentralHeatingState string `json:"central_heating_state"`
	RoomTempSetpoint    int    `json:"room_temp_setpoint"`
	DHWSetpoint         int    `json:"dhw_setpoint"`
	SupplySetpoint      int    `json:"supply_flow_setpoint"`
}

const SettingsTopic string = "nilan/settings"

func CreateSettingsDTO(settings nilan.Settings) SettingsDTO {
	return SettingsDTO{
		FanState:            ventilationState(settings),
		Speed:               ventilationSpeed(settings),
		Mode:                ventilationMode(settings),
		DHWState:            onOffString(*settings.CentralHeatingIsOn),
		CentralHeatingState: onOffString(!*settings.CentralHeatingPaused),
		RoomTempSetpoint:    *settings.DesiredRoomTemperature / 10,
		DHWSetpoint:         *settings.DesiredDHWTemperature / 10,
		SupplySetpoint:      *settings.SetpointSupplyTemperature / 10,
	}
}

func ventilationState(settings nilan.Settings) string {
	if *settings.VentilationOnPause {
		return "OFF"
	} else {
		return "ON"
	}
}

func ventilationSpeed(settings nilan.Settings) int {
	if *settings.VentilationOnPause {
		return 0
	}
	switch *settings.FanSpeed {
	case nilan.FanSpeedLow:
		return 1
	case nilan.FanSpeedNormal:
		return 2
	case nilan.FanSpeedHigh:
		return 3
	case nilan.FanSpeedVeryHigh:
		return 4
	default:
		return 0
	}
}

func ventilationMode(settings nilan.Settings) string {
	switch *settings.VentilationMode {
	case 1:
		return "cooling"
	case 2:
		return "heating"
	default:
		return "auto"
	}
}

func FanSpeed(speed int) *nilan.FanSpeed {
	switch speed {
	case 1:
		return fanSpeedAddr(nilan.FanSpeedLow)
	case 2:
		return fanSpeedAddr(nilan.FanSpeedNormal)
	case 3:
		return fanSpeedAddr(nilan.FanSpeedHigh)
	case 4:
		return fanSpeedAddr(nilan.FanSpeedVeryHigh)
	default:
		return nil
	}
}

func fanSpeedAddr(fanSpeed nilan.FanSpeed) *nilan.FanSpeed {
	speed := fanSpeed
	return &speed
}

func Mode(mode string) *int {
	switch mode {
	case "auto":
		return intAddr(0)
	case "cooling":
		return intAddr(1)
	case "heating":
		return intAddr(2)
	default:
		return nil
	}
}

func intAddr(i int) *int {
	iVar := i
	return &iVar
}

func OnOffString(on bool) string {
	if on {
		return "ON"
	} else {
		return "OFF"
	}
}
