package helper

import (
	"context"
	"ekak_kabupaten_madiun/model/web"
	"net/http"
	"strings"
)

const (
	RoleSuperAdmin = "super_admin"
	RoleAdminOpd   = "admin_opd"
)

// GetUserClaimsFromContext mengambil JWT claims dari request context
// Mengembalikan claims dan boolean yang menandakan apakah claims valid
func GetUserClaimsFromContext(ctx context.Context) (web.JWTClaim, bool) {
	claims, ok := ctx.Value(UserInfoKey).(web.JWTClaim)
	if !ok {
		return web.JWTClaim{}, false
	}
	return claims, true
}

// HasRole mengecek apakah user memiliki role tertentu
func HasRole(claims web.JWTClaim, role string) bool {
	for _, userRole := range claims.Roles {
		if strings.EqualFold(userRole, role) {
			return true
		}
	}
	return false
}

// HasAnyRole mengecek apakah user memiliki salah satu role dari daftar role
func HasAnyRole(claims web.JWTClaim, roles ...string) bool {
	for _, role := range roles {
		if HasRole(claims, role) {
			return true
		}
	}
	return false
}

// CheckSuperAdminRole memeriksa apakah user memiliki role super_admin
// Jika tidak, akan menulis response error ke writer
// Mengembalikan true jika user adalah super_admin, false jika tidak
func CheckSuperAdminRole(writer http.ResponseWriter, request *http.Request) bool {
	claims, ok := GetUserClaimsFromContext(request.Context())
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak valid",
		}
		WriteToResponseBody(writer, webResponse)
		return false
	}

	if !HasRole(claims, RoleSuperAdmin) {
		webResponse := web.WebResponse{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Data:   "Akses ditolak. Hanya super_admin yang dapat mengakses fitur ini",
		}
		WriteToResponseBody(writer, webResponse)
		return false
	}

	return true
}

// CheckRole memeriksa apakah user memiliki role tertentu
// Versi generic yang bisa digunakan untuk role apapun
func CheckRole(writer http.ResponseWriter, request *http.Request, requiredRole string, errorMessage string) bool {
	claims, ok := GetUserClaimsFromContext(request.Context())
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak valid",
		}
		WriteToResponseBody(writer, webResponse)
		return false
	}

	if !HasRole(claims, requiredRole) {
		if errorMessage == "" {
			errorMessage = "Akses ditolak. Anda tidak memiliki akses untuk fitur ini"
		}
		webResponse := web.WebResponse{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Data:   errorMessage,
		}
		WriteToResponseBody(writer, webResponse)
		return false
	}

	return true
}

// CheckAnyRole memeriksa apakah user memiliki salah satu role dari daftar role
func CheckAnyRole(writer http.ResponseWriter, request *http.Request, requiredRoles []string, errorMessage string) bool {
	claims, ok := GetUserClaimsFromContext(request.Context())
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak valid",
		}
		WriteToResponseBody(writer, webResponse)
		return false
	}

	if !HasAnyRole(claims, requiredRoles...) {
		if errorMessage == "" {
			errorMessage = "Akses ditolak. Anda tidak memiliki akses untuk fitur ini"
		}
		webResponse := web.WebResponse{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Data:   errorMessage,
		}
		WriteToResponseBody(writer, webResponse)
		return false
	}

	return true
}

// CheckSuperAdminOrAdminOpdRole memeriksa apakah user memiliki role super_admin atau admin_opd
// Mengembalikan claims dan boolean yang menandakan apakah user memiliki salah satu role tersebut
func CheckSuperAdminOrAdminOpdRole(writer http.ResponseWriter, request *http.Request) (web.JWTClaim, bool) {
	claims, ok := GetUserClaimsFromContext(request.Context())
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak valid",
		}
		WriteToResponseBody(writer, webResponse)
		return web.JWTClaim{}, false
	}

	hasSuperAdmin := HasRole(claims, RoleSuperAdmin)
	hasAdminOpd := HasRole(claims, RoleAdminOpd)

	if !hasSuperAdmin && !hasAdminOpd {
		webResponse := web.WebResponse{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Data:   "Akses ditolak. Hanya super_admin atau admin_opd yang dapat mengakses fitur ini",
		}
		WriteToResponseBody(writer, webResponse)
		return web.JWTClaim{}, false
	}

	return claims, true
}

// GetFilteredKodeOpd mengembalikan kode_opd yang sesuai berdasarkan role user
// - Jika super_admin: mengembalikan kode_opd dari query parameter (bisa kosong untuk semua)
// - Jika admin_opd: mengembalikan kode_opd dari token JWT (mengabaikan query parameter)
func GetFilteredKodeOpd(claims web.JWTClaim, queryKodeOpd string) string {
	if HasRole(claims, RoleSuperAdmin) {
		// Super admin bisa akses semua, gunakan kode_opd dari query atau kosong
		return queryKodeOpd
	}

	if HasRole(claims, RoleAdminOpd) {
		// Admin OPD hanya bisa akses sesuai kode_opd mereka
		return claims.KodeOpd
	}

	return queryKodeOpd
}

// ValidateKodeOpdAccess memvalidasi apakah user dapat mengakses kode_opd tertentu
// - Jika super_admin: selalu mengembalikan true (bebas akses)
// - Jika admin_opd: hanya mengembalikan true jika kode_opd sama dengan di token JWT
func ValidateKodeOpdAccess(claims web.JWTClaim, targetKodeOpd string) bool {
	if HasRole(claims, RoleSuperAdmin) {
		return true // Super admin bebas akses
	}

	if HasRole(claims, RoleAdminOpd) {
		return claims.KodeOpd == targetKodeOpd
	}

	return false
}

// ValidateKodeOpdAccessWithError memvalidasi kode_opd dan mengembalikan error response jika tidak valid
// Mengembalikan true jika valid, false jika tidak valid (dan sudah menulis response error)
func ValidateKodeOpdAccessWithError(writer http.ResponseWriter, claims web.JWTClaim, targetKodeOpd string, errorMessage string) bool {
	if ValidateKodeOpdAccess(claims, targetKodeOpd) {
		return true
	}

	if errorMessage == "" {
		errorMessage = "Akses ditolak. Anda hanya dapat mengakses data dari OPD Anda sendiri"
	}

	webResponse := web.WebResponse{
		Code:   http.StatusForbidden,
		Status: "FORBIDDEN",
		Data:   errorMessage,
	}
	WriteToResponseBody(writer, webResponse)
	return false
}
