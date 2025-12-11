package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/user"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository    repository.UserRepository
	RoleRepository    repository.RoleRepository
	PegawaiRepository repository.PegawaiRepository
	OpdRepository     repository.OpdRepository
	DB                *sql.DB
}

func NewUserServiceImpl(userRepository repository.UserRepository, roleRepository repository.RoleRepository, pegawaiRepository repository.PegawaiRepository, opdRepository repository.OpdRepository, db *sql.DB) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository:    userRepository,
		RoleRepository:    roleRepository,
		PegawaiRepository: pegawaiRepository,
		OpdRepository:     opdRepository,
		DB:                db,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request user.UserCreateRequest) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input dasar
	if request.Nip == "" {
		return user.UserResponse{}, errors.New("nip harus diisi")
	}
	if request.Password == "" {
		return user.UserResponse{}, errors.New("password harus diisi")
	}
	if len(request.Role) == 0 {
		return user.UserResponse{}, errors.New("role harus diisi")
	}

	// Validasi NIP dengan data pegawai
	_, err = service.PegawaiRepository.FindByNip(ctx, tx, request.Nip)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.UserResponse{}, errors.New("nip tidak terdaftar di data pegawai")
		}
		return user.UserResponse{}, err
	}

	// Siapkan slice untuk menyimpan roles
	var roles []domain.Roles

	// Validasi dan ambil semua role yang dipilih
	for _, roleRequest := range request.Role {
		role, err := service.RoleRepository.FindById(ctx, tx, roleRequest.RoleId)
		if err != nil {
			if err == sql.ErrNoRows {
				return user.UserResponse{}, errors.New("role tidak ditemukan")
			}
			return user.UserResponse{}, err
		}
		roles = append(roles, role)
	}

	userDomain := domain.Users{
		Nip:      request.Nip,
		Email:    helper.EmptyStringIfNull(request.Email),
		Password: request.Password,
		IsActive: request.IsActive,
		Role:     roles,
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDomain.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.UserResponse{}, err
	}
	userDomain.Password = string(hashedPassword)

	createdUser, err := service.UserRepository.Create(ctx, tx, userDomain)
	if err != nil {
		return user.UserResponse{}, err
	}

	// Konversi role ke response
	var roleResponses []user.RoleResponse
	for _, role := range createdUser.Role {
		roleResponses = append(roleResponses, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	response := user.UserResponse{
		Id:       createdUser.Id,
		Nip:      createdUser.Nip,
		Email:    createdUser.Email,
		IsActive: createdUser.IsActive,
		Role:     roleResponses,
	}

	return response, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request user.UserUpdateRequest) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi user exists
	existingUser, err := service.UserRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return user.UserResponse{}, err
	}
	if existingUser.Id == 0 {
		return user.UserResponse{}, errors.New("user tidak ditemukan")
	}

	// Validasi input dasar
	if request.Nip == "" {
		return user.UserResponse{}, errors.New("nip harus diisi")
	}
	if request.Email == "" {
		return user.UserResponse{}, errors.New("email harus diisi")
	}
	if len(request.Role) == 0 {
		return user.UserResponse{}, errors.New("role harus diisi")
	}

	// Siapkan slice untuk menyimpan roles
	var roles []domain.Roles

	// Validasi dan ambil semua role yang dipilih
	for _, roleRequest := range request.Role {
		role, err := service.RoleRepository.FindById(ctx, tx, roleRequest.RoleId)
		if err != nil {
			if err == sql.ErrNoRows {
				return user.UserResponse{}, errors.New("role tidak ditemukan")
			}
			return user.UserResponse{}, err
		}
		roles = append(roles, role)
	}

	userDomain := domain.Users{
		Id:       existingUser.Id,
		Nip:      request.Nip,
		Email:    request.Email,
		IsActive: request.IsActive,
		Role:     roles,
	}

	// Handle password update
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return user.UserResponse{}, err
		}
		userDomain.Password = string(hashedPassword)
	} else {
		userDomain.Password = existingUser.Password
	}

	updatedUser, err := service.UserRepository.Update(ctx, tx, userDomain)
	if err != nil {
		return user.UserResponse{}, err
	}

	// Konversi role ke response
	var roleResponses []user.RoleResponse
	for _, role := range updatedUser.Role {
		roleResponses = append(roleResponses, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	response := user.UserResponse{
		Id:       updatedUser.Id,
		Nip:      updatedUser.Nip,
		Email:    updatedUser.Email,
		IsActive: updatedUser.IsActive,
		Role:     roleResponses,
	}

	return response, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	existingUser, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}
	if existingUser.Id == 0 {
		return errors.New("user tidak ditemukan")
	}

	err = service.UserRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) FindAll(ctx context.Context, kodeOpd string) ([]user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}

	var userResponses []user.UserResponse
	for _, u := range users {
		var roles []user.RoleResponse

		// Sort roles berdasarkan ID untuk konsistensi
		sortedRoles := make([]domain.Roles, len(u.Role))
		copy(sortedRoles, u.Role)
		sort.Slice(sortedRoles, func(i, j int) bool {
			return sortedRoles[i].Id < sortedRoles[j].Id
		})

		for _, role := range sortedRoles {
			roles = append(roles, user.RoleResponse{
				Id:   role.Id,
				Role: role.Role,
			})
		}

		pegawaiDomain, _ := service.PegawaiRepository.FindByNip(ctx, tx, u.Nip)

		userResponse := user.UserResponse{
			Id:          u.Id,
			Nip:         u.Nip,
			Email:       u.Email,
			NamaPegawai: pegawaiDomain.NamaPegawai,
			IsActive:    u.IsActive,
			Role:        roles,
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}
func (service *UserServiceImpl) FindById(ctx context.Context, id int) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cari user berdasarkan ID
	userDomain, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return user.UserResponse{}, err
	}

	// Cek apakah user ditemukan
	if userDomain.Id == 0 {
		return user.UserResponse{}, errors.New("user tidak ditemukan")
	}

	// Konversi role domain ke role response
	var roles []user.RoleResponse
	for _, role := range userDomain.Role {
		roles = append(roles, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	pegawaiDomain, _ := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)

	// Convert ke response
	response := user.UserResponse{
		Id:          userDomain.Id,
		Nip:         userDomain.Nip,
		Email:       userDomain.Email,
		NamaPegawai: pegawaiDomain.NamaPegawai,
		IsActive:    userDomain.IsActive,
		Role:        roles,
	}

	return response, nil
}

// func (service *UserServiceImpl) Login(ctx context.Context, request user.UserLoginRequest) (user.UserLoginResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return user.UserLoginResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	if request.Username == "" {
// 		return user.UserLoginResponse{}, errors.New("email atau nip harus diisi")
// 	}
// 	if request.Password == "" {
// 		return user.UserLoginResponse{}, errors.New("password harus diisi")
// 	}

// 	userDomain, err := service.UserRepository.FindByEmailOrNip(ctx, tx, request.Username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return user.UserLoginResponse{}, errors.New("username atau password salah")
// 		}
// 		return user.UserLoginResponse{}, err
// 	}

// 	pegawaiDomain, err := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)
// 	if err != nil {
// 		return user.UserLoginResponse{}, err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(userDomain.Password), []byte(request.Password))
// 	if err != nil {
// 		return user.UserLoginResponse{}, errors.New("username atau password salah")
// 	}

// 	if !userDomain.IsActive {
// 		return user.UserLoginResponse{}, errors.New("akun tidak aktif")
// 	}

// 	roleNames := make([]string, 0, len(userDomain.Role))
// 	for _, role := range userDomain.Role {
// 		roleNames = append(roleNames, role.Role)
// 	}

// 	token := helper.CreateNewJWT(
// 		userDomain.Id,
// 		pegawaiDomain.Id,
// 		userDomain.Email,
// 		userDomain.Nip,
// 		pegawaiDomain.KodeOpd,
// 		roleNames,
// 	)

// 	response := user.UserLoginResponse{
// 		Token: token,
// 	}

// 	return response, nil
// }

type LoginAttempt struct {
	Count        int
	LastAttempt  time.Time
	BlockedUntil time.Time
}

var (
	// Map untuk menyimpan login attempts per NIP
	loginAttempts     = make(map[string]*LoginAttempt)
	loginAttemptsLock sync.RWMutex

	// Konfigurasi throttling
	maxLoginAttempts = 3                // Maksimal 3 kali salah
	blockDuration    = 3 * time.Minute  // Block selama 3 menit
	attemptWindow    = 15 * time.Minute // Reset counter setelah 15 menit tidak ada attempt
)

// ✅ Fungsi untuk cek apakah user sedang di-block
func isUserBlocked(nip string) (bool, time.Duration) {
	loginAttemptsLock.RLock()
	defer loginAttemptsLock.RUnlock()

	attempt, exists := loginAttempts[nip]
	if !exists {
		return false, 0
	}

	// Cek apakah masih dalam periode block
	if time.Now().Before(attempt.BlockedUntil) {
		remainingTime := time.Until(attempt.BlockedUntil)
		return true, remainingTime
	}

	return false, 0
}

// ✅ Fungsi untuk record login failure
func recordLoginFailure(nip string) (blocked bool, remainingTime time.Duration) {
	loginAttemptsLock.Lock()
	defer loginAttemptsLock.Unlock()

	now := time.Now()
	attempt, exists := loginAttempts[nip]

	if !exists {
		// First attempt
		loginAttempts[nip] = &LoginAttempt{
			Count:       1,
			LastAttempt: now,
		}
		return false, 0
	}

	// Reset counter jika sudah lewat dari attemptWindow
	if now.Sub(attempt.LastAttempt) > attemptWindow {
		attempt.Count = 1
		attempt.LastAttempt = now
		attempt.BlockedUntil = time.Time{}
		return false, 0
	}

	// Increment counter
	attempt.Count++
	attempt.LastAttempt = now

	// Block jika sudah mencapai max attempts
	if attempt.Count >= maxLoginAttempts {
		attempt.BlockedUntil = now.Add(blockDuration)
		remainingTime = blockDuration
		return true, remainingTime
	}

	return false, 0
}

// ✅ Fungsi untuk reset login attempts (dipanggil saat login berhasil)
func resetLoginAttempts(nip string) {
	loginAttemptsLock.Lock()
	defer loginAttemptsLock.Unlock()

	delete(loginAttempts, nip)
}

// ✅ Fungsi untuk cleanup old entries (optional, dipanggil secara periodik)
func cleanupOldLoginAttempts() {
	loginAttemptsLock.Lock()
	defer loginAttemptsLock.Unlock()

	now := time.Now()
	for nip, attempt := range loginAttempts {
		// Hapus entry yang sudah tidak relevan (lebih dari attemptWindow)
		if now.Sub(attempt.LastAttempt) > attemptWindow && now.After(attempt.BlockedUntil) {
			delete(loginAttempts, nip)
		}
	}
}

func (service *UserServiceImpl) Login(ctx context.Context, request user.UserLoginRequest) (user.UserLoginResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserLoginResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input
	if request.Username == "" {
		return user.UserLoginResponse{}, errors.New("nip harus diisi")
	}
	if request.Password == "" {
		return user.UserLoginResponse{}, errors.New("password harus diisi")
	}

	// ✅ CEK APAKAH USER SEDANG DI-BLOCK
	blocked, remainingTime := isUserBlocked(request.Username)
	if blocked {
		minutes := int(remainingTime.Minutes())
		seconds := int(remainingTime.Seconds()) % 60
		return user.UserLoginResponse{}, fmt.Errorf(
			"akun diblokir karena terlalu banyak percobaan login gagal. silakan coba lagi dalam %d menit %d detik",
			minutes, seconds,
		)
	}

	// Cari user berdasarkan NIP saja
	userDomain, err := service.UserRepository.FindByEmailOrNip(ctx, tx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// ✅ RECORD FAILURE - NIP tidak ditemukan
			recordLoginFailure(request.Username)
			return user.UserLoginResponse{}, errors.New("nip atau password salah")
		}
		return user.UserLoginResponse{}, err
	}

	// Pastikan username yang digunakan adalah NIP
	if userDomain.Nip != request.Username {
		// ✅ RECORD FAILURE - Bukan NIP
		recordLoginFailure(request.Username)
		return user.UserLoginResponse{}, errors.New("silakan login menggunakan NIP")
	}

	pegawaiDomain, err := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)
	if err != nil {
		return user.UserLoginResponse{}, err
	}

	opdDomain, err := service.OpdRepository.FindByKodeOpd(ctx, tx, pegawaiDomain.KodeOpd)
	if err != nil {
		return user.UserLoginResponse{}, err
	}

	// ✅ VALIDASI PASSWORD
	err = bcrypt.CompareHashAndPassword([]byte(userDomain.Password), []byte(request.Password))
	if err != nil {
		// ✅ RECORD FAILURE - Password salah
		blocked, remainingTime := recordLoginFailure(request.Username)

		if blocked {
			minutes := int(remainingTime.Minutes())
			seconds := int(remainingTime.Seconds()) % 60
			return user.UserLoginResponse{}, fmt.Errorf(
				"terlalu banyak percobaan login gagal. akun diblokir selama %d menit %d detik",
				minutes, seconds,
			)
		}

		// Hitung berapa kali lagi bisa mencoba
		loginAttemptsLock.RLock()
		attempt := loginAttempts[request.Username]
		remainingAttempts := maxLoginAttempts - attempt.Count
		loginAttemptsLock.RUnlock()

		if remainingAttempts > 0 {
			return user.UserLoginResponse{}, fmt.Errorf(
				"nip atau password salah. sisa percobaan: %d kali",
				remainingAttempts,
			)
		}

		return user.UserLoginResponse{}, errors.New("nip atau password salah")
	}

	// ✅ VALIDASI AKUN AKTIF
	if !userDomain.IsActive {
		return user.UserLoginResponse{}, errors.New("akun tidak aktif")
	}

	// ✅ LOGIN BERHASIL - RESET ATTEMPTS
	resetLoginAttempts(request.Username)

	roleNames := make([]string, 0, len(userDomain.Role))
	for _, role := range userDomain.Role {
		roleNames = append(roleNames, role.Role)
	}

	token := helper.CreateNewJWT(
		userDomain.Id,
		pegawaiDomain.Id,
		userDomain.Email,
		userDomain.Nip,
		pegawaiDomain.KodeOpd,
		opdDomain.NamaOpd,
		pegawaiDomain.NamaPegawai,
		roleNames,
	)

	response := user.UserLoginResponse{
		Token: token,
	}

	return response, nil
}

func (service *UserServiceImpl) FindByKodeOpdAndRole(ctx context.Context, kodeOpd string, roleName string) ([]user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input
	if kodeOpd == "" {
		return nil, errors.New("kode opd harus diisi")
	}
	if roleName == "" {
		return nil, errors.New("role harus diisi")
	}

	users, err := service.UserRepository.FindByKodeOpdAndRole(ctx, tx, kodeOpd, roleName)
	if err != nil {
		return nil, err
	}

	var userResponses []user.UserResponse
	for _, u := range users {
		var roles []user.RoleResponse
		for _, role := range u.Role {
			roles = append(roles, user.RoleResponse{
				Id:   role.Id,
				Role: role.Role,
			})
		}

		userResponse := user.UserResponse{
			Id:          u.Id,
			Nip:         u.Nip,
			IsActive:    u.IsActive,
			PegawaiId:   u.PegawaiId,
			NamaPegawai: u.NamaPegawai,
			Role:        roles,
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

func (service *UserServiceImpl) FindByNip(ctx context.Context, nip string) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	userDomain, err := service.UserRepository.FindByNip(ctx, tx, nip)
	if err != nil {
		return user.UserResponse{}, err
	}

	var roles []user.RoleResponse
	for _, role := range userDomain.Role {
		roles = append(roles, user.RoleResponse{
			Id:   role.Id,
			Role: role.Role,
		})
	}

	pegawaiDomain, err := service.PegawaiRepository.FindByNip(ctx, tx, userDomain.Nip)
	if err != nil {
		return user.UserResponse{}, err
	}

	userResponse := user.UserResponse{
		Nip:         userDomain.Nip,
		NamaPegawai: pegawaiDomain.NamaPegawai,
		IsActive:    userDomain.IsActive,
		Role:        roles,
	}

	return userResponse, nil
}

func (service *UserServiceImpl) CekAdminOpd(ctx context.Context) ([]user.CekAdminOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua OPD
	allOpd, err := service.OpdRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Ambil semua user dengan role admin_opd
	adminUsers, err := service.UserRepository.CekAdminOpd(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Buat map untuk grouping user berdasarkan kode_opd
	adminByOpd := make(map[string][]user.AdminOpdUserDetail)
	for _, u := range adminUsers {
		adminDetail := user.AdminOpdUserDetail{
			UserId:      u.Id,
			Nip:         u.Nip,
			NamaPegawai: u.NamaPegawai,
			Email:       u.Email,
			IsActive:    u.IsActive,
		}
		adminByOpd[u.KodeOpd] = append(adminByOpd[u.KodeOpd], adminDetail)
	}

	// Build response: semua OPD dengan admin users (atau array kosong jika tidak ada)
	var response []user.CekAdminOpdResponse
	for _, opd := range allOpd {
		opdResponse := user.CekAdminOpdResponse{
			KodeOpd:    opd.KodeOpd,
			NamaOpd:    opd.NamaOpd,
			AdminUsers: []user.AdminOpdUserDetail{}, // inisialisasi dengan array kosong
		}

		// Jika ada admin di OPD ini, masukkan datanya
		if admins, exists := adminByOpd[opd.KodeOpd]; exists {
			opdResponse.AdminUsers = admins
		}

		response = append(response, opdResponse)
	}

	return response, nil
}

func (service *UserServiceImpl) GetKodeOpdByNip(ctx context.Context, nip string) (string, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return "", err
	}
	defer helper.CommitOrRollback(tx)

	pegawai, err := service.PegawaiRepository.FindByNip(ctx, tx, nip)
	if err != nil {
		return "", err
	}

	return pegawai.KodeOpd, nil
}

func (service *UserServiceImpl) ChangePassword(ctx context.Context, request user.UserChangePasswordRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input dasar
	if request.Password1 == "" {
		return errors.New("password baru harus diisi")
	}
	if request.Password2 == "" {
		return errors.New("konfirmasi password harus diisi")
	}
	if request.Password1 != request.Password2 {
		return errors.New("password baru dan konfirmasi password tidak sama")
	}

	// Validasi panjang password
	if len(request.Password1) < 6 {
		return errors.New("password minimal 6 karakter")
	}

	// Validasi NIP harus diisi
	if request.Nip == "" {
		return errors.New("NIP harus diisi")
	}

	// Cari user berdasarkan NIP
	userDomain, err := service.UserRepository.FindByNip(ctx, tx, request.Nip)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return fmt.Errorf("gagal mencari user: %v", err)
	}

	if userDomain.Id == 0 {
		return errors.New("user tidak ditemukan")
	}

	// Validasi old_password jika diperlukan (untuk role selain super_admin dan admin_opd)
	if request.OldPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(userDomain.Password), []byte(request.OldPassword))
		if err != nil {
			return errors.New("password lama tidak sesuai")
		}
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password1), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password: %v", err)
	}

	// Update password menggunakan repository method
	err = service.UserRepository.UpdatePassword(ctx, tx, userDomain.Nip, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("gagal mengubah password: %v", err)
	}

	return nil
}
