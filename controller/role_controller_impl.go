package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/user"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type RoleControllerImpl struct {
	RoleService service.RoleService
}

func NewRoleControllerImpl(roleService service.RoleService) *RoleControllerImpl {
	return &RoleControllerImpl{
		RoleService: roleService,
	}
}

func (controller *RoleControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}
	roleCreateRequest := user.RoleCreateRequest{}
	helper.ReadFromRequestBody(request, &roleCreateRequest)

	roleResponse, err := controller.RoleService.Create(request.Context(), roleCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed create role",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create role",
		Data:   roleResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RoleControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}

	roleId := params.ByName("id")
	id, err := strconv.Atoi(roleId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed update role",
			Data:   "invalid role id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	roleUpdateRequest := user.RoleUpdateRequest{}
	helper.ReadFromRequestBody(request, &roleUpdateRequest)
	roleUpdateRequest.Id = id

	roleResponse, err := controller.RoleService.Update(request.Context(), roleUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed update role",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success update role",
		Data:   roleResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RoleControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}

	roleId := params.ByName("id")
	id, err := strconv.Atoi(roleId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed delete role",
			Data:   "invalid role id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.RoleService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed delete role",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success delete role",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RoleControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}

	roleId := params.ByName("id")
	id, err := strconv.Atoi(roleId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed get role",
			Data:   "invalid role id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	roleResponse, err := controller.RoleService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed get role",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success get role",
		Data:   roleResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RoleControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helper.CheckSuperAdminRole(writer, request) {
		return
	}

	roleResponses, err := controller.RoleService.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed get roles",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success get roles",
		Data:   roleResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
