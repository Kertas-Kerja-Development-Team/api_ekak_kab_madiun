package controller

import (
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PohonKinerjaOpdControllerImpl struct {
	PohonKinerjaOpdService service.PohonKinerjaOpdService
}

func NewPohonKinerjaOpdControllerImpl(pohonKinerjaOpdService service.PohonKinerjaOpdService) *PohonKinerjaOpdControllerImpl {
	return &PohonKinerjaOpdControllerImpl{
		PohonKinerjaOpdService: pohonKinerjaOpdService,
	}
}

func (controller *PohonKinerjaOpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonKinerjaCreateRequest := pohonkinerja.PohonKinerjaCreateRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaCreateRequest)

	response, err := controller.PohonKinerjaOpdService.Create(request.Context(), pohonKinerjaCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   response,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonKinerjaUpdateRequest := pohonkinerja.PohonKinerjaUpdateRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaUpdateRequest)

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	pohonKinerjaUpdateRequest.Id = id

	// Panggil service Update
	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.Update(request.Context(), pohonKinerjaUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Update Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := controller.PohonKinerjaOpdService.Delete(request.Context(), params.ByName("id"))
	if err != nil {

		webResponse := web.WebResponse{
			Code:   400,
			Status: "Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Delete Pohon Kinerja",
		Data:   nil,
	}

	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *PohonKinerjaOpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Dapatkan id dari params
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service FindById
	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Berhasil Mendapatkan Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	// Jika kodeOpd atau tahun kosong, kembalikan response null
	if kodeOpd == "" || tahun == "" {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service FindAll
	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindAll(request.Context(), kodeOpd, tahun)
	if err != nil {
		// Jika tidak ada data, kembalikan response sukses dengan data null
		if err == sql.ErrNoRows {
			webResponse := web.WebResponse{
				Code:   200,
				Status: "OK",
				Data:   nil,
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Untuk error lainnya
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get All Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) FindStrategicNoParent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	if kodeOpd == "" || tahun == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "kode_opd dan tahun harus diisi",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindStrategicNoParent(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get Strategic No Parent",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) DeletePelaksana(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := controller.PohonKinerjaOpdService.DeletePelaksana(request.Context(), params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusInternalServerError)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Delete Pelaksana",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaOpdControllerImpl) FindPokinByPelaksana(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pelaksanaId := params.ByName("pegawai_id")
	tahun := params.ByName("tahun")

	// Validasi parameter
	if pelaksanaId == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID Pegawai harus diisi",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	if tahun == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "Parameter tahun harus diisi",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pohonKinerjaResponse, err := controller.PohonKinerjaOpdService.FindPokinByPelaksana(request.Context(), pelaksanaId, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Berhasil Mendapatkan Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
