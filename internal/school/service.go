package school

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type School struct {
	ID         string `json:"id"`
	PublicID   string `json:"public_id"`
	Name       string `json:"name"`
	ShortName  string `json:"short_name,omitempty"`
	SchoolType string `json:"school_type"`
	City       string `json:"city,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Email      string `json:"email,omitempty"`
}

type CreateInput struct {
	PublicID   string `json:"public_id"`
	Name       string `json:"name"`
	ShortName  string `json:"short_name"`
	SchoolType string `json:"school_type"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type Service struct{ db *pgxpool.Pool }

func NewService(db *pgxpool.Pool) *Service { return &Service{db: db} }

func (s *Service) Create(ctx context.Context, in CreateInput) (School, error) {
	in.PublicID = strings.TrimSpace(in.PublicID)
	in.Name = strings.TrimSpace(in.Name)
	if in.PublicID == "" || in.Name == "" || in.SchoolType == "" {
		return School{}, fmt.Errorf("public_id, name and school_type are required")
	}
	var out School
	err := s.db.QueryRow(ctx, `
		INSERT INTO schools (public_id,name,short_name,school_type,city,phone,email)
		VALUES ($1,$2,NULLIF($3,''),$4,NULLIF($5,''),NULLIF($6,''),NULLIF($7,''))
		RETURNING id::text, public_id, name, COALESCE(short_name,''), school_type::text,
		          COALESCE(city,''), COALESCE(phone,''), COALESCE(email,'')`,
		in.PublicID, in.Name, in.ShortName, in.SchoolType, in.City, in.Phone, in.Email,
	).Scan(&out.ID, &out.PublicID, &out.Name, &out.ShortName, &out.SchoolType, &out.City, &out.Phone, &out.Email)
	return out, err
}

func (s *Service) GetByPublicID(ctx context.Context, publicID string) (School, error) {
	var out School
	err := s.db.QueryRow(ctx, `
		SELECT id::text, public_id, name, COALESCE(short_name,''), school_type::text,
		       COALESCE(city,''), COALESCE(phone,''), COALESCE(email,'')
		FROM schools WHERE public_id=$1`, publicID,
	).Scan(&out.ID, &out.PublicID, &out.Name, &out.ShortName, &out.SchoolType, &out.City, &out.Phone, &out.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return School{}, ErrNotFound
	}
	return out, err
}

var ErrNotFound = errors.New("school not found")
