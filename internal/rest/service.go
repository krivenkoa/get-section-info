package rest

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "skat-vending.com/selection-info/internal/docs"
	"skat-vending.com/selection-info/internal/service"
)

type Storage interface {
	CreateAccessToken(userId string) (string, error)
	CreateRefreshToken(userId string) (string, error)
}

// Service represents rest api
type Service struct {
	Storage  Storage
	Sections *service.Sections
}

//go:generate swag init --parseInternal --parseDependency --generalInfo ./service.go --output ../docs
//go:generate mv ../docs/swagger.json ../../api/openapi-spec
//go:generate mv ../docs/swagger.yaml ../../api/openapi-spec

// @title Content program rest gateway API
// @version 1.0
// @description This service allows webui to access content program functionality
// @BasePath /api/v1
// @license.name Proprietary
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (s *Service) Mount(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.Handler())
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/sections", s.getSections) // Post /sections
	})
}
