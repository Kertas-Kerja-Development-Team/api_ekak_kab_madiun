package urusanrespon

import "ekak_kabupaten_madiun/model/web/bidangurusanresponse"

type UrusanResponse struct {
	Id           string                                      `json:"id"`
	KodeUrusan   string                                      `json:"kode_urusan"`
	NamaUrusan   string                                      `json:"nama_urusan"`
	BidangUrusan []bidangurusanresponse.BidangUrusanResponse `json:"bidang_urusan"`
}
