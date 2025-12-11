package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProgramControllerImpl struct {
	ProgramService service.ProgramService
}

func NewProgramControllerImpl(programService service.ProgramService) *ProgramControllerImpl {
	return &ProgramControllerImpl{
		ProgramService: programService,
	}
}

func (controller *ProgramControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	programCreateRequest := programkegiatan.ProgramKegiatanCreateRequest{}
	helper.ReadFromRequestBody(request, &programCreateRequest)

	programResponse, err := controller.ProgramService.Create(request.Context(), programCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   201,
		Status: "Success",
		Data:   programResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	programUpdateRequest := programkegiatan.ProgramKegiatanUpdateRequest{}
	helper.ReadFromRequestBody(request, &programUpdateRequest)

	programUpdateRequest.Id = params.ByName("programId")
	programResponse, err := controller.ProgramService.Update(request.Context(), programUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success",
		Data:   programResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	programId := params.ByName("id")
	err := controller.ProgramService.Delete(request.Context(), programId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success",
		Data:   "Program berhasil dihapus",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	programId := params.ByName("id")
	programResponse, err := controller.ProgramService.FindById(request.Context(), programId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success",
		Data:   programResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ProgramControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	programResponses, err := controller.ProgramService.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success",
		Data:   programResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
