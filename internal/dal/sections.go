package dal

import (
	"context"
	"database/sql"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"skat-vending.com/selection-info/internal/errs"
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
func (s *Sections) InnerThemesList(ctx context.Context, idRazd, idOperator, idOtdel string) ([]api.InnerTheme, error) {
	r, err := s.db.QueryContext(ctx, `SELECT id_theme, name_theme from themes where isnull(archive,0)=0 
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
                                      )`, idRazd, idOperator, idOtdel)
	if err != nil {
		return nil, errors.Wrapf(errs.ErrInternalDatabase, "retrieving inner themes: %v", err)
	}
	defer closeRows(r)

	result := make([]api.InnerTheme, 0)
	for r.Next() {
		catalog, err := s.innerThemeFromRecord(r)
		if err != nil {
			log15.Error("load inner theme record from database", "err", err)
			continue
		}
		result = append(result, *catalog)
	}

	return result, nil
}

func (s *Sections) innerThemeFromRecord(r *sql.Rows) (*api.InnerTheme, error) {
	var (
		id   int
		name string
	)
	if err := r.Scan(&id, &name); err != nil {
		return nil, errors.Wrapf(errs.ErrInternalDatabase, "scan catalog from record: %v", err)
	}

	return &api.InnerTheme{
		IdTheme:   id,
		NameTheme: name,
	}, nil
}

func closeRows(r *sql.Rows) {
	if r != nil {
		if err := r.Close(); err != nil {
			log15.Error("close query", "error", err)
		}
	}
}
