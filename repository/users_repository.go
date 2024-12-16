package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, users domain.Users) (domain.Users, error)
	Update(ctx context.Context, tx *sql.Tx, users domain.Users) (domain.Users, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.Users, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Users, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindByEmailOrNip(ctx context.Context, tx *sql.Tx, username string) (domain.Users, error)
}
