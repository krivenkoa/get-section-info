package api

// Section represents info about section
type SectionRequest struct {
	IdOtdel    string `json:"id_otdel" example:"123"`
	IdRazdel   string `json:"id_razdel" example:"123123"`
	IdOperator string `json:"id_operator" example:"2"`
}
