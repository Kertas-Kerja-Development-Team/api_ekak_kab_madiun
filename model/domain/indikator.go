package domain

type Indikator struct {
	Id               string
	RencanaKinerjaId string
	Indikator        string
	Tahun            string
	Target           []Target
}
