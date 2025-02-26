package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type RencanaKinerjaRepository interface {
	Create(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error)
	FindAll(ctx context.Context, tx *sql.Tx, pegawaiId string, kodeOPD string, tahun string) ([]domain.RencanaKinerja, error)
	FindIndikatorbyRekinId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, targetId string) ([]domain.Target, error)
	FindById(ctx context.Context, tx *sql.Tx, id string, kodeOPD string, tahun string) (domain.RencanaKinerja, error)
	Update(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindAllRincianKak(ctx context.Context, tx *sql.Tx, rencanakinerjaid, pegawaiId string) ([]domain.RencanaKinerja, error)
	//sasaran opd
	CreateSasaranOpd(ctx context.Context, tx *sql.Tx, sasaranOpd domain.RencanaKinerja) (domain.RencanaKinerja, error)
	UpdateSasaranOpd(ctx context.Context, tx *sql.Tx, sasaranOpd domain.RencanaKinerja) (domain.RencanaKinerja, error)
	RekinsasaranOpd(ctx context.Context, tx *sql.Tx, pegawaiId string, kodeOPD string, tahun string) ([]domain.RencanaKinerja, error)
	FindIndikatorSasaranbyRekinId(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.Indikator, error)
	FindTargetByIndikatorIdAndTahun(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) ([]domain.Target, error)
}
