package pohonkinerja

import (
	"ekak_kabupaten_madiun/model/web/opdmaster"
)

type CascadingOpdResponse struct {
	KodeOpd      string                          `json:"kode_opd"`
	NamaOpd      string                          `json:"nama_opd"`
	Tahun        string                          `json:"tahun"`
	TujuanOpd    []TujuanOpdCascadingResponse    `json:"tujuan_opd"`
	BidangUrusan []BidangUrusanResponse          `json:"bidang_urusan"`
	Strategics   []StrategicCascadingOpdResponse `json:"childs"`
}

type TujuanOpdCascadingResponse struct {
	Id         int                       `json:"id"`
	KodeOpd    string                    `json:"kode_opd"`
	Tujuan     string                    `json:"tujuan"`
	KodeBidang string                    `json:"kode_bidang_urusan"`
	NamaBidang string                    `json:"nama_bidang_urusan"`
	Indikator  []IndikatorTujuanResponse `json:"indikator"`
}

type StrategicCascadingOpdResponse struct {
	Id                     int                            `json:"id"`
	Parent                 *int                           `json:"parent"`
	Strategi               string                         `json:"nama_pohon"`
	JenisPohon             string                         `json:"jenis_pohon"`
	LevelPohon             int                            `json:"level_pohon"`
	Keterangan             string                         `json:"keterangan"`
	KeteranganCrosscutting *string                        `json:"keterangan_crosscutting,omitempty"`
	Status                 string                         `json:"status"`
	CountReview            int                            `json:"jumlah_review"`
	KodeOpd                opdmaster.OpdResponseForAll    `json:"perangkat_daerah"`
	Program                []ProgramResponse              `json:"program"`
	IsActive               bool                           `json:"is_active"`
	RencanaKinerja         []RencanaKinerjaResponse       `json:"rencana_kinerja"`
	Indikator              []IndikatorResponse            `json:"indikator"`
	Tacticals              []TacticalCascadingOpdResponse `json:"childs,omitempty"`
	Review                 []ReviewResponse               `json:"review,omitempty"`
}

type TacticalCascadingOpdResponse struct {
	Id                     int                               `json:"id"`
	Parent                 int                               `json:"parent"`
	Strategi               string                            `json:"nama_pohon"`
	JenisPohon             string                            `json:"jenis_pohon"`
	LevelPohon             int                               `json:"level_pohon"`
	Keterangan             string                            `json:"keterangan"`
	KeteranganCrosscutting *string                           `json:"keterangan_crosscutting,omitempty"`
	Status                 string                            `json:"status"`
	CountReview            int                               `json:"jumlah_review"`
	Program                []ProgramResponse                 `json:"program"`
	KodeOpd                opdmaster.OpdResponseForAll       `json:"perangkat_daerah"`
	IsActive               bool                              `json:"is_active"`
	RencanaKinerja         []RencanaKinerjaResponse          `json:"rencana_kinerja"`
	Indikator              []IndikatorResponse               `json:"indikator"`
	Operationals           []OperationalCascadingOpdResponse `json:"childs,omitempty"`
	Review                 []ReviewResponse                  `json:"review,omitempty"`
}

type OperationalCascadingOpdResponse struct {
	Id                     int                                 `json:"id"`
	Parent                 int                                 `json:"parent"`
	Strategi               string                              `json:"nama_pohon"`
	JenisPohon             string                              `json:"jenis_pohon"`
	LevelPohon             int                                 `json:"level_pohon"`
	Keterangan             string                              `json:"keterangan"`
	KeteranganCrosscutting *string                             `json:"keterangan_crosscutting,omitempty"`
	Status                 string                              `json:"status"`
	CountReview            int                                 `json:"jumlah_review"`
	KodeOpd                opdmaster.OpdResponseForAll         `json:"perangkat_daerah"`
	IsActive               bool                                `json:"is_active"`
	RencanaKinerja         []RencanaKinerjaOperationalResponse `json:"rencana_kinerja"`
	Indikator              []IndikatorResponse                 `json:"indikator"`
	Childs                 []OperationalNOpdCascadingResponse  `json:"childs,omitempty"`
	Review                 []ReviewResponse                    `json:"review,omitempty"`
}

type OperationalNOpdCascadingResponse struct {
	Id             int                                  `json:"id"`
	Parent         int                                  `json:"parent"`
	Strategi       string                               `json:"nama_pohon"`
	JenisPohon     string                               `json:"jenis_pohon"`
	LevelPohon     int                                  `json:"level_pohon"`
	Keterangan     string                               `json:"keterangan"`
	Status         string                               `json:"status"`
	CountReview    int                                  `json:"jumlah_review"`
	KodeOpd        opdmaster.OpdResponseForAll          `json:"perangkat_daerah"`
	IsActive       bool                                 `json:"is_active"`
	RencanaKinerja []RencanaKinerjaOperationalNResponse `json:"rencana_kinerja"`
	Indikator      []IndikatorResponse                  `json:"indikator"`
	Childs         []OperationalNOpdCascadingResponse   `json:"childs,omitempty"`
	Review         []ReviewResponse                     `json:"review,omitempty"`
}

type RencanaKinerjaResponse struct {
	Id                 string              `json:"id_rencana_kinerja,omitempty"`
	IdPohon            int                 `json:"id_pohon,omitempty"`
	NamaPohon          string              `json:"nama_pohon,omitempty"`
	NamaRencanaKinerja string              `json:"nama_rencana_kinerja,omitempty"`
	Tahun              string              `json:"tahun,omitempty"`
	PegawaiId          string              `json:"pegawai_id,omitempty"`
	NamaPegawai        string              `json:"nama_pegawai,omitempty"`
	Indikator          []IndikatorResponse `json:"indikator,omitempty"`
	Program            []ProgramResponse   `json:"program,omitempty"`
}

type RencanaKinerjaOperationalResponse struct {
	Id                 string              `json:"id_rencana_kinerja,omitempty"`
	IdPohon            int                 `json:"id_pohon,omitempty"`
	NamaPohon          string              `json:"nama_pohon,omitempty"`
	NamaRencanaKinerja string              `json:"nama_rencana_kinerja,omitempty"`
	Tahun              string              `json:"tahun,omitempty"`
	PegawaiId          string              `json:"pegawai_id,omitempty"`
	NamaPegawai        string              `json:"nama_pegawai,omitempty"`
	KodeSubkegiatan    string              `json:"kode_subkegiatan"`
	NamaSubkegiatan    string              `json:"nama_subkegiatan"`
	Indikator          []IndikatorResponse `json:"indikator,omitempty"`
	KodeKegiatan       string              `json:"kode_kegiatan"`
	NamaKegiatan       string              `json:"nama_kegiatan"`
}

type RencanaKinerjaOperationalNResponse struct {
	Id                 string              `json:"id_rencana_kinerja,omitempty"`
	IdPohon            int                 `json:"id_pohon,omitempty"`
	NamaPohon          string              `json:"nama_pohon,omitempty"`
	NamaRencanaKinerja string              `json:"nama_rencana_kinerja,omitempty"`
	Tahun              string              `json:"tahun,omitempty"`
	PegawaiId          string              `json:"pegawai_id,omitempty"`
	NamaPegawai        string              `json:"nama_pegawai,omitempty"`
	KodeSubkegiatan    string              `json:"kode_subkegiatan,omitempty"`
	NamaSubkegiatan    string              `json:"nama_subkegiatan,omitempty"`
	Indikator          []IndikatorResponse `json:"indikator,omitempty"`
}

type BidangUrusanResponse struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
}

type ProgramResponse struct {
	KodeProgram string `json:"kode_program"`
	NamaProgram string `json:"nama_program"`
}
