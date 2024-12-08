package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type TujuanOpdRepository interface {
	Create(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) (domain.TujuanOpd, error)
	Update(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.TujuanOpd, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.TujuanOpd, error)
	FindIndikatorByTujuanId(ctx context.Context, tx *sql.Tx, tujuanOpdId int) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) ([]domain.Target, error)
}
