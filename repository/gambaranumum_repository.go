package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type GambaranUmumRepository interface {
	Create(ctx context.Context, tx *sql.Tx, gambaranUmum domain.GambaranUmum) (domain.GambaranUmum, error)
	Update(ctx context.Context, tx *sql.Tx, gambaranUmum domain.GambaranUmum) (domain.GambaranUmum, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domain.GambaranUmum, error)
	FindAll(ctx context.Context, tx *sql.Tx, rekinId string, pegawaiId string) ([]domain.GambaranUmum, error)
	GetLastUrutan(ctx context.Context, tx *sql.Tx) (int, error)
}
