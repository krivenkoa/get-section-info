package api

// GetDeviceResponse is a generic device response container
type GetSectionResponse struct {
	Response
	Sections []Section `json:"sections"`
}

// Section represents info about section
type Section struct {
	IdOtdel    string `json:"id_otdel" example:"123"`
	IdRazdel   string `json:"id_razdel" example:"123123"`
	IdOperator string `json:"id_operator" example:"2"`
}
