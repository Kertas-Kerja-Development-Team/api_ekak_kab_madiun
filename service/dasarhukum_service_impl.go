package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/dasarhukum"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type DasarHukumServiceImpl struct {
	DasarHukumRepository repository.DasarHukumRepository
	DB                   *sql.DB
}

func NewDasarHukumServiceImpl(dasarHukumRepository repository.DasarHukumRepository, DB *sql.DB) *DasarHukumServiceImpl {
	return &DasarHukumServiceImpl{
		DasarHukumRepository: dasarHukumRepository,
		DB:                   DB,
	}
}

func (service *DasarHukumServiceImpl) Create(ctx context.Context, request dasarhukum.DasarHukumCreateRequest) (dasarhukum.DasarHukumResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return dasarhukum.DasarHukumResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Membuat UUID dengan format yang diinginkan
	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("DASHU-REKIN-%s", randomDigits)

	// Mendapatkan urutan terakhir
	lastUrutan, err := service.DasarHukumRepository.GetLastUrutan(ctx, tx)
	if err != nil {
		return dasarhukum.DasarHukumResponse{}, err
	}

	// Menambahkan 1 ke urutan terakhir
	newUrutan := lastUrutan + 1

	dasarHukum := domain.DasarHukum{
		Id:               uuId,
		RekinId:          request.RekinId,
		PegawaiId:        request.PegawaiId,
		Urutan:           newUrutan,
		PeraturanTerkait: request.PeraturanTerkait,
		Uraian:           request.Uraian,
	}

	dasarHukum, err = service.DasarHukumRepository.Create(ctx, tx, dasarHukum)
	if err != nil {
		return dasarhukum.DasarHukumResponse{}, err
	}

	return helper.ToDasarHukumResponse(dasarHukum), nil
}

func (service *DasarHukumServiceImpl) Update(ctx context.Context, request dasarhukum.DasarHukumUpdateRequest) (dasarhukum.DasarHukumResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return dasarhukum.DasarHukumResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	dasarHukum, err := service.DasarHukumRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return dasarhukum.DasarHukumResponse{}, err
	}

	dasarHukum.PeraturanTerkait = request.PeraturanTerkait
	dasarHukum.Uraian = request.Uraian

	updatedDasarHukum, err := service.DasarHukumRepository.Update(ctx, tx, dasarHukum)
	if err != nil {
		return dasarhukum.DasarHukumResponse{}, err
	}

	return helper.ToDasarHukumResponse(updatedDasarHukum), nil
}

func (service *DasarHukumServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = service.DasarHukumRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *DasarHukumServiceImpl) FindAll(ctx context.Context, rekinId string, pegawaiId string) ([]dasarhukum.DasarHukumResponse, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	dasarHukums, err := service.DasarHukumRepository.FindAll(ctx, tx, rekinId, pegawaiId)
	if err != nil {
		return []dasarhukum.DasarHukumResponse{}, err
	}

	return helper.ToDasarHukumResponses(dasarHukums), nil
}

func (service *DasarHukumServiceImpl) FindById(ctx context.Context, id string) (dasarhukum.DasarHukumResponse, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	dasarHukum, err := service.DasarHukumRepository.FindById(ctx, tx, id)
	if err != nil {
		return dasarhukum.DasarHukumResponse{}, err
	}

	return helper.ToDasarHukumResponse(dasarHukum), nil
}
