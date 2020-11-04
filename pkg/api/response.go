package api

// Section is a response structure.
type Section struct {
	Success          bool             `json:"success" example:"false"`
	Description      []Description    `json:"description"`
	NameRazdel       string           `json:"name_razdel" example:"qwerty"`
	Archive          bool             `json:"archive" example:"false"`
	DateArchive      *string          `json:"date_archive" example:"2016-02-20"`
	CountInnerThemes int              `json:"count_inner_themes" example:"0"`
	CountOuterThemes int              `json:"count_outer_themes" example:"0"`
	InnerThemes      []InnerTheme     `json:"inner_themes" example:""`
	OuterThemes      []OuterTheme     `json:"outer_themes" example:""`
	OtdelRazdel      map[string]Otdel `json:"otdel_razdel" example:""`
}

type InnerTheme struct {
	IdTheme   int    `json:"id_theme" example:"123"`
	NameTheme string `json:"name_theme" example:""`
}

type OuterTheme struct {
	IdTheme   int    `json:"id_theme" example:"123"`
	NameTheme string `json:"name_theme" example:""`
	Tax       bool   `json:"tax" example:"true"`
}

type Otdel struct {
	Windows []int  `json:"windows" example:""`
	Limit   string `json:"limit" example:""`
}

type Description struct {
	Message    string  `json:"message" example:"Операция выполнена успешно"`
	Reason     *string `json:"reason" example:""`
	Stacktrace *string `json:"stacktrace" example:""`
}
