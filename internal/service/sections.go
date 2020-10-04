package service

import (
	"context"

	"skat-vending.com/selection-info/pkg/api"
)

// Sections represents service for sections management
type Sections struct {
	// TODO: store dao here
}

// NewSections returns new instance of Sections service
func NewSections() *Sections {
	return &Sections{}
}

// Get returns all sections info
func (s *Sections) Get(ctx context.Context) ([]api.Section, error) {
	ss := make([]api.Section, 0)
	return ss, nil
}
