package kegiatan

import (
	"time"
)

type KegiatanResponse struct {
	Id           string    `json:"id"`
	NamaKegiatan string    `json:"nama_kegiatan"`
	KodeKegiatan string    `json:"kode_kegiatan"`
	KodeOPD      string    `json:"kode_opd"`
	CreatedAt    time.Time `json:"created_at"`
	Indikator    []IndikatorResponse
}

type IndikatorResponse struct {
	Id        string           `json:"id"`
	ProgramId string           `json:"program_id"`
	Indikator string           `json:"indikator"`
	Tahun     string           `json:"tahun"`
	Target    []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
