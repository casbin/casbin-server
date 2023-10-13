package server

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type Repository struct {
	tracer trace.Tracer
	db     *gorm.DB
}

type IRepository interface {
	RemoveFilteredNamedPolicy(ctx context.Context, ptype string, fieldIndex int, fieldValues ...string) (bool, error)
	RemoveNamedPolicy(ctx context.Context, ptype string, rule []string) (bool, error)
	AddNamedPolicy(ctx context.Context, ptype string, rule []string) (bool, error)
	Enforce(ctx context.Context, rule []string) (bool, error)
}

func NewRepository(tracer trace.Tracer, db *gorm.DB) *Repository {
	return &Repository{
		tracer: tracer,
		db:     db,
	}
}

// TODO: implement method RemoveFilteredNamedPolicy in repository directly query to db
func (r *Repository) RemoveFilteredNamedPolicy(ctx context.Context, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

// TODO: implement method AddNamedPolicy in repository directly query to db
func (r *Repository) RemoveNamedPolicies(ctx context.Context, ptype string, rules [][]string) error {
	return nil
}

// TODO: implement method AddNamedPolicy in repository directly query to db
func (r *Repository) AddNamedPolicy(ctx context.Context, ptype string, rules [][]string) error {
	return nil
}

// TODO: implement method Enforce in repository directly query to db
func (r *Repository) Enforce(ctx context.Context, rules [][]string) (bool, error) {
	return false, nil
}
