package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"skat-vending.com/selection-info/internal/dal"

	"skat-vending.com/selection-info/pkg/api"
)

// Sections represents service for sections management
type Sections struct {
	dal *dal.Sections
}

// NewSections returns new instance of Sections service
func NewSections(db *sql.DB) *Sections {
	return &Sections{
		dal: dal.NewSection(db),
	}
}

// Get returns all sections info
func (s *Sections) Get(ctx context.Context, req api.SectionRequest) (*api.Section, error) {
	sec, err := s.dal.GetSectionBaseParams(ctx, req.IdRazdel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed get base section params")
	}

	innerThemes, err := s.dal.InnerThemesList(ctx, req.IdRazdel, req.IdOperator, req.IdOtdel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed get inner themes")
	}

	outerThemes, err := s.dal.OuterThemesList(ctx, req.IdRazdel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed get inner themes")
	}

	section := &api.Section{
		Success:          true,
		NameRazdel:       sec.NameRazdel,
		Archive:          sec.Archive,
		DateArchive:      sec.DateArchive,
		CountInnerThemes: len(innerThemes),
		CountOuterThemes: len(innerThemes),
		InnerThemes:      innerThemes,
		OuterThemes:      outerThemes,
		OtdelRazdel:      make(map[string]api.Otdel),
	}

	return section, nil
}
