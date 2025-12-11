package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/bidangurusanresponse"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type BidangUrusanControllerImpl struct {
	BidangUrusanService service.BidangUrusanService
}

func NewBidangUrusanControllerImpl(bidangUrusanService service.BidangUrusanService) *BidangUrusanControllerImpl {
	return &BidangUrusanControllerImpl{
		BidangUrusanService: bidangUrusanService,
	}
}

func (controller *BidangUrusanControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	bidangUrusanCreateRequest := bidangurusanresponse.BidangUrusanCreateRequest{}
	helper.ReadFromRequestBody(request, &bidangUrusanCreateRequest)

	bidangUrusanCreateResponse, err := controller.BidangUrusanService.Create(request.Context(), bidangUrusanCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bidangUrusanCreateResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *BidangUrusanControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	bidangUrusanUpdateRequest := bidangurusanresponse.BidangUrusanUpdateRequest{}
	helper.ReadFromRequestBody(request, &bidangUrusanUpdateRequest)

	bidangUrusanUpdateRequest.Id = params.ByName("id")

	bidangUrusanUpdateResponse, err := controller.BidangUrusanService.Update(request.Context(), bidangUrusanUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bidangUrusanUpdateResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *BidangUrusanControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	bidangUrusanId := params.ByName("id")

	bidangUrusanResponse, err := controller.BidangUrusanService.FindById(request.Context(), bidangUrusanId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bidangUrusanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *BidangUrusanControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	bidangUrusanResponses, err := controller.BidangUrusanService.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bidangUrusanResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *BidangUrusanControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	bidangUrusanId := params.ByName("id")

	err := controller.BidangUrusanService.Delete(request.Context(), bidangUrusanId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *BidangUrusanControllerImpl) FindByKodeOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	kodeOpd := params.ByName("kode_opd")

	bidangUrusanResponses, err := controller.BidangUrusanService.FindByKodeOpd(request.Context(), kodeOpd)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   bidangUrusanResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
