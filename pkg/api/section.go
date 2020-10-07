package api

// GetDeviceResponse is a generic device response container
type GetSectionResponse struct {
	Response
	Sections []Section `json:"sections"`
}

// Section represents info about section
type Section struct {
	ID string `json:"id" example:"5f6516278fea4dfe56868aaf"`
}
