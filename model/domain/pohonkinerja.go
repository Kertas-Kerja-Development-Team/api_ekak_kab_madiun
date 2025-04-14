package domain

import (
	"database/sql"
	"time"
)

type PohonKinerja struct {
	IdCrosscutting         int
	Id                     int
	Parent                 int
	NamaPohon              string
	KodeOpd                string
	NamaOpd                string
	Keterangan             string
	KeteranganCrosscutting *string
	Tahun                  string
	JenisPohon             string
	LevelPohon             int
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Indikator              []Indikator
	Pelaksana              []PelaksanaPokin
	Status                 string
	CloneFrom              int
	Crosscutting           []Crosscutting
	PegawaiAction          interface{}
	CrosscuttingTo         int
	CountReview            int
	IsActive               bool
	//tambahan
	RencanaKinerja  []RencanaKinerja
	KegiatanId      sql.NullString
	SubkegiatanId   sql.NullString
	IsDeleted       bool
	NamaKegiatan    sql.NullString
	KodeKegiatan    sql.NullString
	NamaSubkegiatan sql.NullString
	KodeSubkegiatan sql.NullString
	Strategi        string
	PelaksanaIds    string
}

type PegawaiAction struct {
	ApproveBy *string
	RejectBy  *string
	ApproveAt *time.Time
	RejectAt  *time.Time
}
