package dto

import "github.com/pjuzeliunas/nilan"

type ErrorsDTO struct {
	OldFilter   string `json:"old_filter"`
	OtherErrors string `json:"other_errors"`
}

func CreateErrorsDTO(errors nilan.Errors) ErrorsDTO {
	return ErrorsDTO{
		OldFilter:   onOffString(errors.OldFilterWarning),
		OtherErrors: onOffString(errors.OtherErrors),
	}
}

func onOffString(on bool) string {
	if on {
		return "ON"
	} else {
		return "OFF"
	}
}
