package periodetahun

type PeriodeCreateRequest struct {
	Id         int      `json:"id"`
	TahunAwal  string   `json:"tahun_awal"`
	TahunAkhir string   `json:"tahun_akhir"`
	TahunList  []string `json:"tahun_list"`
}

type TahunPeriodeCreateRequest struct {
	Tahun string `json:"tahun"`
}
