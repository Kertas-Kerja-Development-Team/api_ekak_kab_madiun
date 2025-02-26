package domain

import (
	"database/sql"
	"time"
)

type Indikator struct {
	Id               string
	PokinId          string
	ProgramId        string
	RencanaKinerjaId string
	KegiatanId       string
	SubKegiatanId    string
	TujuanOpdId      int
	TujuanPemdaId    int
	SasaranPemdaId   int
	Indikator        string
	Tahun            string
	CloneFrom        string
	Sumber           string
	ParentId         int
	ParentOpdId      string
	AsalIku          string
	ParentName       string
	ManualIK         *ManualIKSasaranOpd
	CreatedAt        time.Time
	TahunAwal        string
	TahunAkhir       string
	JenisPeriode     string
	Target           []Target
	RencanaKinerja   RencanaKinerja
	RumusPerhitungan sql.NullString
	SumberData       sql.NullString
}
