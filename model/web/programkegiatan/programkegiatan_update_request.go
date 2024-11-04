package programkegiatan

type ProgramKegiatanUpdateRequest struct {
	Id          string                   `json:"id"`
	KodeProgram string                   `json:"kode_program"`
	NamaProgram string                   `json:"nama_program"`
	KodeOPD     string                   `json:"kode_opd"`
	Tahun       string                   `json:"tahun"`
	IsActive    bool                     `json:"is_active"`
	Indikator   []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id        string                `json:"id"`
	ProgramId string                `json:"program_id"`
	Indikator string                `json:"indikator"`
	Tahun     string                `json:"tahun"`
	Target    []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
