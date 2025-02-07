package sasaranpemda

type SasaranPemdaUpdateRequest struct {
	Id               int                      `json:"id"`
	SasaranPemdaId   int                      `json:"sasaran_pemda_id"`
	PeriodeId        int                      `json:"periode_id"`
	RumusPerhitungan string                   `json:"rumus_perhitungan"`
	SumberData       string                   `json:"sumber_data"`
	Indikator        []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id             string                `json:"id"`
	SasaranPemdaId string                `json:"sasaran_id"`
	Indikator      string                `json:"indikator"`
	Target         []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}
