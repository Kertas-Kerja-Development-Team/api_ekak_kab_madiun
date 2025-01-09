package rencanakinerja

type ManualIKResponse struct {
	Id                  int        `json:"id"`
	IndikatorId         string     `json:"indikator_id"`
	Perspektif          string     `json:"perspektif"`
	TujuanRekin         string     `json:"tujuan_rekin"`
	Definisi            string     `json:"definisi"`
	KeyActivities       string     `json:"key_activities"`
	Formula             string     `json:"formula"`
	JenisIndikator      string     `json:"jenis_indikator"`
	OutputData          OutputData `json:"output_data"`
	UnitPenanggungJawab string     `json:"unit_penanggung_jawab"`
	UnitPenyediaJasa    string     `json:"unit_penyedia_jasa"`
	SumberData          string     `json:"sumber_data"`
	JangkaWaktuAwal     string     `json:"jangka_waktu_awal"`
	JangkaWaktuAkhir    string     `json:"jangka_waktu_akhir"`
	PeriodePelaporan    string     `json:"periode_pelaporan"`
}

type OutputData struct {
	Kinerja  bool `json:"kinerja"`
	Penduduk bool `json:"penduduk"`
	Spatial  bool `json:"spatial"`
}