package dal

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"skat-vending.com/selection-info/internal/errs"
	"skat-vending.com/selection-info/internal/utils"
	"skat-vending.com/selection-info/pkg/api"
)

// Sections represent DAL over sql.DB type
type Sections struct {
	db *sql.DB
}

// NewSection creates new instance of Sections
func NewSection(db *sql.DB) *Sections {
	return &Sections{db: db}
}

// InnerThemesList returns InnerTheme list
func (s *Sections) InnerThemesList(ctx context.Context, idRazd, idOperator, idOtdel string) ([]api.InnerTheme, *api.Description, error) {
	sql := `SELECT id_theme, name_theme from themes where isnull(archive,0)=0 
				and id_theme in (select id_theme 
                 from razd_theme 
                    where id_razd = ?
                  AND ISNULL(stat_table, '') = ''
                  AND id_theme NOT IN (
                                       SELECT et.themeId FROM dbo.themes_threshold th
                                       JOIN emfc_test_result_th et ON et.themeId = th.themeId
                                       WHERE et.userId = ? 
                                         AND et.result < th.threshold 
                                         AND th.otdelId = ?
                                      ))`

	r, err := s.db.QueryContext(ctx, sql, idRazd, idOperator, idOtdel)
	if err != nil {
		description := &api.Description{
			Message: "retrieving inner themes",
			Reason:  utils.String(fmt.Sprintf("%s; idRazd=%s, idOperator=%s, idOtdel=%s", sql, idRazd, idOperator, idOtdel)),
		}
		return nil, description, errors.Wrapf(errs.ErrInternalDatabase, "retrieving inner themes: %v", err)
	}
	defer closeRows(r)

	result := make([]api.InnerTheme, 0)
	for r.Next() {
		innnerTheme, err := s.innerThemeFromRecord(r)
		if err != nil {
			logrus.WithError(err).Error("load inner theme record from database")
			continue
		}
		result = append(result, *innnerTheme)
	}

	return result, nil, nil
}

// OuterThemesList returns OuterTheme list
func (s *Sections) OuterThemesList(ctx context.Context, idRazd string) ([]api.OuterTheme, *api.Description, error) {
	sql := `SELECT id_theme,name_theme,tax
           FROM themes
           WHERE id_theme IN (
                              SELECT id_theme 
                              FROM razd_theme 
                              WHERE id_razd = ?
                              AND stat_table IS NOT NULL 
                              AND stat_table <> ''
                             )`

	r, err := s.db.QueryContext(ctx, sql, idRazd)
	if err != nil {
		description := &api.Description{
			Message: "retrieving outer themes",
			Reason:  utils.String(fmt.Sprintf("%s; idRazd=%s", sql, idRazd)),
		}
		return nil, description, errors.Wrapf(errs.ErrInternalDatabase, "retrieving outer themes: %v", err)
	}
	defer closeRows(r)

	result := make([]api.OuterTheme, 0)
	for r.Next() {
		outerTheme, err := s.outerThemeFromRecord(r)
		if err != nil {
			logrus.WithError(err).Error("load inner theme record from database")
			continue
		}
		result = append(result, *outerTheme)
	}

	return result, nil, nil
}

// GetSectionBaseParams returns name_razdel, archive, date_archive params in Section
func (s *Sections) GetSectionBaseParams(ctx context.Context, idRazd string) (*api.Section, *api.Description, error) {
	sql := `SELECT razdel, archive, date_archive FROM razdel WHERE id = ?`
	r, err := s.db.QueryContext(ctx, sql, idRazd)
	if err != nil {
		description := &api.Description{
			Message: "retrieving section base params",
			Reason:  utils.String(fmt.Sprintf("sql: %s; idRazd=%s", sql, idRazd)),
		}
		return nil, description, errors.Wrapf(errs.ErrInternalDatabase, "retrieving section base params: %v", err)
	}
	defer closeRows(r)

	if !r.Next() {
		description := &api.Description{
			Message: "section params not found",
			Reason:  utils.String(fmt.Sprintf("sql: %s; idRazd=%s", sql, idRazd)),
		}
		return nil, description, errors.Wrapf(errs.ErrNotFound, "section params not found: %v", idRazd)
	}

	res, err := s.sectionFromRecord(r)
	if err != nil {
		description := &api.Description{
			Message: "convert section base params from record",
			Reason:  utils.String(fmt.Sprintf("sql: %s; idRazd=%s", sql, idRazd)),
		}
		return nil, description, err
	}

	return res, nil, nil
}

// GetSectionBaseParams returns name_razdel, archive, date_archive params in Section
func (s *Sections) GetOtdelRazdel(ctx context.Context, idRazd, idOtdel string) (map[string]api.Otdel, *api.Description, error) {
	sql := `SELECT id_otdel,limit FROM otdel_razdel WHERE id_razdel = ?`
	r, err := s.db.QueryContext(ctx, sql, idRazd)
	if err != nil {
		description := &api.Description{
			Message: "retrieving inner themes",
			Reason:  utils.String(fmt.Sprintf("%s; idRazd=%s, idOtdel=%s", sql, idRazd, idOtdel)),
		}
		return nil, description, errors.Wrapf(errs.ErrInternalDatabase, "retrieving otdel razdel: %v", err)
	}
	defer closeRows(r)

	result := make(map[string]api.Otdel)
	for r.Next() {
		var (
			id    int
			limit string
		)
		if err := r.Scan(&id, &limit); err != nil {
			logrus.WithError(err).Error("scan idOtdel and limit from record")
			continue
		}

		windows, err := s.windowsList(ctx, idRazd, idOtdel)
		if err != nil {
			logrus.WithError(err).Error("scan idOtdel and limit from record")
			continue
		}

		result[string(id)] = api.Otdel{
			Windows: windows,
			Limit:   limit,
		}
	}

	return result, nil, nil
}

func (s *Sections) windowsList(ctx context.Context, idRazd, idOtdel string) ([]int, error) {
	sql := `SELECT id_wnd FROM window WHERE razdel = ? AND id_otdel = ?`
	r, err := s.db.QueryContext(ctx, sql, idRazd, idOtdel)
	if err != nil {
		return nil, errors.Wrapf(errs.ErrInternalDatabase, "retrieving windows: %v", err)
	}
	defer closeRows(r)

	result := make([]int, 0)
	for r.Next() {
		var id int
		if err := r.Scan(&id); err != nil {
			logrus.WithError(err).Error("load window id record from database")
			continue
		}
		result = append(result, id)
	}

	return result, nil
}

func (s *Sections) innerThemeFromRecord(r *sql.Rows) (*api.InnerTheme, error) {
	var (
		id   int
		name string
	)
	if err := r.Scan(&id, &name); err != nil {
		return nil, errors.Wrapf(errs.ErrInternalDatabase, "scan inner theme from record: %v", err)
	}

	return &api.InnerTheme{
		IdTheme:   id,
		NameTheme: name,
	}, nil
}

func (s *Sections) outerThemeFromRecord(r *sql.Rows) (*api.OuterTheme, error) {
	var (
		id   int
		name string
		tax  int
	)
	if err := r.Scan(&id, &name, &tax); err != nil {
		return nil, errors.Wrapf(errs.ErrInternalDatabase, "scan outer theme from record: %v", err)
	}

	return &api.OuterTheme{
		IdTheme:   id,
		NameTheme: name,
		Tax:       tax,
	}, nil
}

func (s *Sections) sectionFromRecord(r *sql.Rows) (*api.Section, error) {
	var (
		nameRazdel  string
		archive     bool
		dateArchive string
	)
	if err := r.Scan(&nameRazdel, &archive, &dateArchive); err != nil {
		return nil, errors.Wrapf(errs.ErrInternalDatabase, "scan section from record: %v", err)
	}

	return &api.Section{
		NameRazdel:  nameRazdel,
		Archive:     archive,
		DateArchive: dateArchive,
	}, nil
}

func closeRows(r *sql.Rows) {
	if r != nil {
		if err := r.Close(); err != nil {
			logrus.WithError(err).Error("close query")
		}
	}
}
