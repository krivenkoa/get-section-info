package rest

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"skat-vending.com/selection-info/internal/coder"
)

// GetSections godoc
// @Summary Finds all sections
// @Description Finds all sections
// @ID get-sections
// @Tags sections
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSectionResponse
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /sections [get]
func (s *Service) getSections(w http.ResponseWriter, r *http.Request) {
	result, err := s.Sections.Get(context.Background())
	if err != nil {
		logrus.WithError(err).Errorf("getSections find all sections")
		if err := coder.WriteError(w, r, http.StatusInternalServerError, err.Error()); err != nil {
			logrus.WithError(err).Error("getSections write error")
		}
		return
	}

	if err := coder.WriteData(w, r, result, http.StatusOK); err != nil {
		logrus.WithError(err).Error("getSections writing response")
		return
	}
}
