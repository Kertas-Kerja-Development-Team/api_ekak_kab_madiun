package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/urusanrespon"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UrusanControllerImpl struct {
	UrusanService service.UrusanService
}

func NewUrusanControllerImpl(urusanService service.UrusanService) *UrusanControllerImpl {
	return &UrusanControllerImpl{
		UrusanService: urusanService,
	}
}

func (controller *UrusanControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	urusanCreateRequest := urusanrespon.UrusanCreateRequest{}
	helper.ReadFromRequestBody(request, &urusanCreateRequest)

	urusanResponse, err := controller.UrusanService.Create(request.Context(), urusanCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   urusanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UrusanControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	urusanUpdateRequest := urusanrespon.UrusanUpdateRequest{}
	helper.ReadFromRequestBody(request, &urusanUpdateRequest)

	urusanUpdateRequest.Id = params.ByName("id")

	urusanResponse, err := controller.UrusanService.Update(request.Context(), urusanUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update urusan",
		Data:   urusanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UrusanControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	id := params.ByName("id")

	urusanResponse, err := controller.UrusanService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find urusan by id",
		Data:   urusanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UrusanControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	urusanResponses, err := controller.UrusanService.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find all urusan",
		Data:   urusanResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UrusanControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	id := params.ByName("id")

	err := controller.UrusanService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success delete urusan",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UrusanControllerImpl) FindByKodeOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	kodeOpd := helper.GetFilteredKodeOpd(claims, params.ByName("kode_opd"))

	urusanResponses, err := controller.UrusanService.FindByKodeOpd(request.Context(), kodeOpd)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find urusan by kode opd",
		Data:   urusanResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UrusanControllerImpl) FindUrusanAndBidangByKodeOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	kodeOpd := helper.GetFilteredKodeOpd(claims, params.ByName("kode_opd"))

	urusanResponses, err := controller.UrusanService.FindUrusanAndBidangByKodeOpd(request.Context(), kodeOpd)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find urusan and bidang by kode opd",
		Data:   urusanResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
