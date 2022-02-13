package internal

import "github.com/pjuzeliunas/nilan"

type VentilationDTO struct {
	State string `json:"state"`
	Speed int    `json:"speed"`
	Mode  string `json:"mode"`
}

func CreateVentilationDTO(settings nilan.Settings) VentilationDTO {
	return VentilationDTO{
		State: ventilationState(settings),
		Speed: ventilationSpeed(settings),
		Mode:  ventilationMode(settings),
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
