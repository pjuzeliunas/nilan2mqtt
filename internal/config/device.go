package config

type Device struct {
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
	Identifiers  []string `json:"identifiers"`
}

func NilanDevice() Device {
	return Device{
		Name:         "nilan",
		Manufacturer: "Nilan",
		Model:        "CTS700",
		Identifiers:  []string{"501ff64a-4d56-45cd-8f6c-88c143efabc9"},
	}
}
