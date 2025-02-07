package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type TujuanPemdaRepository interface {
	Create(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error)
	CreateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) (domain.Indikator, error)
	CreateTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error
	Update(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error)
	Delete(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error
	FindById(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) (domain.TujuanPemda, error)
	FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.TujuanPemda, error)
	DeleteIndikator(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error
	IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool
	UpdatePeriode(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error)
	FindAllWithPokin(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.TujuanPemdaWithPokin, error)
}
