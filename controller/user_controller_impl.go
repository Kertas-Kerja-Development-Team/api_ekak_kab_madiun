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

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserControllerImpl(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		userService: userService,
	}
}

func (controller *UserControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Pengecekan role super_admin atau admin_opd
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	userCreateRequest := user.UserCreateRequest{}
	helper.ReadFromRequestBody(request, &userCreateRequest)

	// Validasi kode_opd untuk admin_opd: hanya bisa create user dari OPD mereka
	if !helper.HasRole(claims, helper.RoleSuperAdmin) {
		// Ambil kode_opd dari NIP pegawai
		pegawaiKodeOpd, err := controller.userService.GetKodeOpdByNip(request.Context(), userCreateRequest.Nip)
		if err != nil {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed create user",
				Data:   "NIP tidak ditemukan atau tidak valid",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Validasi apakah kode_opd pegawai sama dengan kode_opd admin_opd
		if !helper.ValidateKodeOpdAccessWithError(writer, claims, pegawaiKodeOpd, "Anda hanya dapat membuat user untuk pegawai dari OPD Anda sendiri") {
			return
		}
	}

	userResponse, err := controller.userService.Create(request.Context(), userCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed create user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create user",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Pengecekan role super_admin atau admin_opd
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	// Parse ID dari URL parameter
	userId := params.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed update user",
			Data:   "invalid user id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	userUpdateRequest := user.UserUpdateRequest{}
	helper.ReadFromRequestBody(request, &userUpdateRequest)
	userUpdateRequest.Id = id

	// Validasi kode_opd untuk admin_opd: hanya bisa update user dari OPD mereka
	if !helper.HasRole(claims, helper.RoleSuperAdmin) {
		// Ambil kode_opd dari NIP pegawai
		pegawaiKodeOpd, err := controller.userService.GetKodeOpdByNip(request.Context(), userUpdateRequest.Nip)
		if err != nil {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed update user",
				Data:   "NIP tidak ditemukan atau tidak valid",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Validasi apakah kode_opd pegawai sama dengan kode_opd admin_opd
		if !helper.ValidateKodeOpdAccessWithError(writer, claims, pegawaiKodeOpd, "Anda hanya dapat mengupdate user untuk pegawai dari OPD Anda sendiri") {
			return
		}
	}

	userResponse, err := controller.userService.Update(request.Context(), userUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed update user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success update user",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Pengecekan role super_admin atau admin_opd
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	userId := params.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed delete user",
			Data:   "invalid user id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Validasi kode_opd untuk admin_opd: hanya bisa delete user dari OPD mereka
	if !helper.HasRole(claims, helper.RoleSuperAdmin) {
		// Ambil user yang akan dihapus untuk mendapatkan NIP
		existingUser, err := controller.userService.FindById(request.Context(), id)
		if err != nil {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed delete user",
				Data:   "user tidak ditemukan",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Ambil kode_opd dari NIP pegawai
		pegawaiKodeOpd, err := controller.userService.GetKodeOpdByNip(request.Context(), existingUser.Nip)
		if err != nil {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed delete user",
				Data:   "NIP tidak ditemukan atau tidak valid",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Validasi apakah kode_opd pegawai sama dengan kode_opd admin_opd
		if !helper.ValidateKodeOpdAccessWithError(writer, claims, pegawaiKodeOpd, "Anda hanya dapat menghapus user untuk pegawai dari OPD Anda sendiri") {
			return
		}
	}

	err = controller.userService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed delete user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success delete user",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Pengecekan role super_admin atau admin_opd
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	queryKodeOpd := request.URL.Query().Get("kode_opd")

	kodeOpd := helper.GetFilteredKodeOpd(claims, queryKodeOpd)
	userResponses, err := controller.userService.FindAll(request.Context(), kodeOpd)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find all user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find all user",
		Data:   userResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId := params.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find by id user",
			Data:   "invalid user id",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	userResponse, err := controller.userService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find by id user",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find by id user",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	loginRequest := user.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &loginRequest)

	loginResponse, err := controller.userService.Login(request.Context(), loginRequest)
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
		Data: map[string]interface{}{
			"token": loginResponse.Token,
		},
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindByKodeOpdAndRole(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Pengecekan role super_admin atau admin_opd
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	// Ambil kode_opd dari query parameter
	queryKodeOpd := request.URL.Query().Get("kode_opd")
	roleName := request.URL.Query().Get("role")

	// Filter kode_opd berdasarkan role:
	// - Super admin: bisa menggunakan query parameter atau kosong (untuk semua)
	// - Admin OPD: otomatis menggunakan kode_opd dari token JWT mereka
	kodeOpd := helper.GetFilteredKodeOpd(claims, queryKodeOpd)

	userResponses, err := controller.userService.FindByKodeOpdAndRole(request.Context(), kodeOpd, roleName)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find by kode opd and role",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find by kode opd and role",
		Data:   userResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindByNip(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Pengecekan role super_admin atau admin_opd
	claims, ok := helper.CheckSuperAdminOrAdminOpdRole(writer, request)
	if !ok {
		return
	}

	nip := params.ByName("nip")

	// Validasi kode_opd untuk admin_opd: hanya bisa find user dari OPD mereka
	if !helper.HasRole(claims, helper.RoleSuperAdmin) {
		// Ambil kode_opd dari NIP pegawai
		pegawaiKodeOpd, err := controller.userService.GetKodeOpdByNip(request.Context(), nip)
		if err != nil {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed find by nip",
				Data:   "NIP tidak ditemukan atau tidak valid",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Validasi apakah kode_opd pegawai sama dengan kode_opd admin_opd
		if !helper.ValidateKodeOpdAccessWithError(writer, claims, pegawaiKodeOpd, "Anda hanya dapat melihat user dari OPD Anda sendiri") {
			return
		}
	}

	userResponse, err := controller.userService.FindByNip(request.Context(), nip)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find by nip",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find by nip",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) CekAdminOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	response, err := controller.userService.CekAdminOpd(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success cek admin opd",
		Data:   response,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) ChangePassword(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil claims dari context untuk validasi role
	claims, ok := helper.GetUserClaimsFromContext(request.Context())
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak valid",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	changePasswordRequest := user.UserChangePasswordRequest{}
	helper.ReadFromRequestBody(request, &changePasswordRequest)

	// Validasi dan filter berdasarkan role
	isSuperAdmin := helper.HasRole(claims, helper.RoleSuperAdmin)
	isAdminOpd := helper.HasRole(claims, helper.RoleAdminOpd)

	if isSuperAdmin {
		// Super admin: bebas ubah password siapa saja, NIP dari request
		if changePasswordRequest.Nip == "" {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed change password",
				Data:   "NIP harus diisi",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}
		// Super admin tidak perlu validasi old_password
		changePasswordRequest.OldPassword = ""

	} else if isAdminOpd {
		// Admin OPD: bisa ubah password user dari OPD mereka
		if changePasswordRequest.Nip == "" {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed change password",
				Data:   "NIP harus diisi",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Validasi kode_opd: pastikan user yang akan diubah password adalah dari OPD yang sama
		pegawaiKodeOpd, err := controller.userService.GetKodeOpdByNip(request.Context(), changePasswordRequest.Nip)
		if err != nil {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed change password",
				Data:   "NIP tidak ditemukan atau tidak valid",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Validasi apakah kode_opd pegawai sama dengan kode_opd admin_opd
		if !helper.ValidateKodeOpdAccessWithError(writer, claims, pegawaiKodeOpd, "Anda hanya dapat mengubah password user dari OPD Anda sendiri") {
			return
		}
		// Admin OPD tidak perlu validasi old_password
		changePasswordRequest.OldPassword = ""

	} else {
		// Role lain: hanya bisa ubah password mereka sendiri
		// NIP diambil dari token (override request)
		changePasswordRequest.Nip = claims.Nip

		// Validasi old_password wajib untuk role selain admin
		if changePasswordRequest.OldPassword == "" {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "failed change password",
				Data:   "password lama harus diisi",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}
	}

	// Panggil service untuk change password
	err := controller.userService.ChangePassword(request.Context(), changePasswordRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed change password",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success change password",
		Data:   "Password berhasil diubah",
	}
	helper.WriteToResponseBody(writer, webResponse)
}
