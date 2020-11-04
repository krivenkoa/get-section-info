package service

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"skat-vending.com/selection-info/internal/dal"

	"skat-vending.com/selection-info/pkg/api"
)

const (
	operationSuccess = "Операция выполнена успешно"
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
func (s *Sections) Get(ctx context.Context, req api.SectionRequest) (api.Section, error) {
	sec, desc, err := s.dal.GetSectionBaseParams(ctx, req.IdRazdel)
	if err != nil {
		section := api.Section{
			Success: false,
			Description: []api.Description{
				*desc,
			},
		}
		return section, errors.Wrapf(err, "failed get base section params")
	}

	innerThemes, desc, err := s.dal.InnerThemesList(ctx, req.IdRazdel, req.IdOperator, req.IdOtdel)
	if err != nil {
		section := api.Section{
			Success: false,
			Description: []api.Description{
				*desc,
			},
		}
		return section, errors.Wrapf(err, "failed get inner themes")
	}

	outerThemes, desc, err := s.dal.OuterThemesList(ctx, req.IdRazdel)
	if err != nil {
		section := api.Section{
			Success: false,
			Description: []api.Description{
				*desc,
			},
		}
		return section, errors.Wrapf(err, "failed get outer themes")
	}

	otdelRazdel, desc, err := s.dal.GetOtdelRazdel(ctx, req.IdRazdel, req.IdOtdel)
	if err != nil {
		section := api.Section{
			Success: false,
			Description: []api.Description{
				*desc,
			},
		}
		return section, errors.Wrapf(err, "failed get otdel razdel")
	}

	section := api.Section{
		Success: true,
		Description: []api.Description{
			{
				Message: operationSuccess,
			},
		},
		NameRazdel:       sec.NameRazdel,
		Archive:          sec.Archive,
		DateArchive:      sec.DateArchive,
		CountInnerThemes: len(innerThemes),
		CountOuterThemes: len(outerThemes),
		InnerThemes:      innerThemes,
		OuterThemes:      outerThemes,
		OtdelRazdel:      otdelRazdel,
	}

	return section, nil
}
