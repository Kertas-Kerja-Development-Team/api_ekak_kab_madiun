package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/sasaranpemda"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SasaranPemdaServiceImpl struct {
	SasaranPemdaRepository repository.SasaranPemdaRepository
	PeriodeRepository      repository.PeriodeRepository
	PohonKinerjaRepository repository.PohonKinerjaRepository
	TujuanPemdaRepository  repository.TujuanPemdaRepository
	DB                     *sql.DB
}

func NewSasaranPemdaServiceImpl(sasaranPemdaRepository repository.SasaranPemdaRepository, periodeRepository repository.PeriodeRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, tujuanPemdaRepository repository.TujuanPemdaRepository, DB *sql.DB) *SasaranPemdaServiceImpl {
	return &SasaranPemdaServiceImpl{
		SasaranPemdaRepository: sasaranPemdaRepository,
		PeriodeRepository:      periodeRepository,
		PohonKinerjaRepository: pohonKinerjaRepository,
		TujuanPemdaRepository:  tujuanPemdaRepository,
		DB:                     DB,
	}
}

func (service *SasaranPemdaServiceImpl) generateRandomId(ctx context.Context, tx *sql.Tx) int {
	rand.Seed(time.Now().UnixNano())
	for {
		// Generate random number between 10000-99999
		id := rand.Intn(90000) + 10000
		if !service.SasaranPemdaRepository.IsIdExists(ctx, tx, id) {
			return id
		}
	}
}

func generateIndikatorIdSasaran() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5] // Mengambil 5 karakter pertama dari UUID
	return fmt.Sprintf("IND-SAS-PMD-%s-%s", currentYear, uuid)
}

func generateTargetIdSasaran() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5]
	return fmt.Sprintf("TRG-SAS-PMD-%s-%s", currentYear, uuid)
}

func (service *SasaranPemdaServiceImpl) Create(ctx context.Context, request sasaranpemda.SasaranPemdaCreateRequest) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi pohon kinerja
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}

	// Validasi level pohon kinerja (1-3)
	if pokinData.LevelPohon < 1 || pokinData.LevelPohon > 3 {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("level pohon kinerja harus berada di antara 1-3, level saat ini: %d", pokinData.LevelPohon)
	}

	// Validasi subtema id belum digunakan
	if service.SasaranPemdaRepository.IsSubtemaIdExists(ctx, tx, request.SubtemaId) {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja dengan id %d sudah digunakan untuk sasaran pemda lain", request.SubtemaId)
	}

	// Validasi tujuan pemda exists
	exists := service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.TujuanPemdaId)
	if !exists {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda dengan id %d tidak ditemukan", request.TujuanPemdaId)
	}

	// Validasi periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	// Persiapkan slice indikator untuk domain
	var indikators []domain.Indikator

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		var targets []domain.Target

		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true

			// Tambahkan target
			target := domain.Target{
				Id:     generateTargetIdSasaran(),
				Target: targetReq.Target,
				Satuan: targetReq.Satuan,
				Tahun:  targetReq.Tahun,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:        generateIndikatorIdSasaran(),
			Indikator: indikatorReq.Indikator,
			RumusPerhitungan: sql.NullString{
				String: indikatorReq.RumusPerhitungan,
				Valid:  true,
			},
			SumberData: sql.NullString{
				String: indikatorReq.SumberData,
				Valid:  true,
			},
			Target: targets,
		}
		indikators = append(indikators, indikator)
	}

	sasaranPemda := domain.SasaranPemda{
		Id:            service.generateRandomId(ctx, tx),
		SubtemaId:     request.SubtemaId,
		TujuanPemdaId: request.TujuanPemdaId,
		SasaranPemda:  request.SasaranPemda,
		PeriodeId:     request.PeriodeId,
		Indikator:     indikators,
	}

	result, err := service.SasaranPemdaRepository.Create(ctx, tx, sasaranPemda)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal membuat sasaran pemda: %v", err)
	}

	// Konversi ke response secara manual
	var indikatorResponses []sasaranpemda.IndikatorResponse
	for _, indikator := range result.Indikator {
		var targetResponses []sasaranpemda.TargetResponse
		// Sort target berdasarkan tahun
		sort.Slice(indikator.Target, func(i, j int) bool {
			return indikator.Target[i].Tahun < indikator.Target[j].Tahun
		})

		for _, target := range indikator.Target {
			targetResponses = append(targetResponses, sasaranpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, sasaranpemda.IndikatorResponse{
			Id:               indikator.Id,
			Indikator:        indikator.Indikator,
			RumusPerhitungan: indikator.RumusPerhitungan.String,
			SumberData:       indikator.SumberData.String,
			Target:           targetResponses,
		})
	}

	// Sort indikator berdasarkan ID
	sort.Slice(indikatorResponses, func(i, j int) bool {
		return indikatorResponses[i].Id < indikatorResponses[j].Id
	})

	return sasaranpemda.SasaranPemdaResponse{
		Id:            result.Id,
		TujuanPemdaId: result.TujuanPemdaId,
		SubtemaId:     result.SubtemaId,
		NamaSubtema:   pokinData.NamaPohon,
		SasaranPemda:  result.SasaranPemda,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}, nil
}
func (service *SasaranPemdaServiceImpl) Update(ctx context.Context, request sasaranpemda.SasaranPemdaUpdateRequest) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Convert request.Id string ke int
	sasaranPemdaId := request.Id

	// Validasi tujuan pemda exists
	_, err = service.TujuanPemdaRepository.FindById(ctx, tx, request.TujuanPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda tidak ditemukan: %v", err)
	}

	if !service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.TujuanPemdaId) {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda dengan id %d tidak ditemukan", request.TujuanPemdaId)
	}
	// Validasi sasaran pemda pemda exists
	sasaranPemda, err := service.SasaranPemdaRepository.FindById(ctx, tx, sasaranPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	//validasi pohon kinerja
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}

	// Validasi level pohon kinerja (1-3)
	if pokinData.LevelPohon < 1 || pokinData.LevelPohon > 3 {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("level pohon kinerja harus berada di antara 1-3, level saat ini: %d", pokinData.LevelPohon)
	}

	periode, err := service.PeriodeRepository.FindById(ctx, tx, sasaranPemda.PeriodeId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true
		}
	}

	// Update data sasaran pemda
	sasaranPemda.TujuanPemdaId = request.TujuanPemdaId
	sasaranPemda.SasaranPemda = request.SasaranPemda

	// Proses indikator
	var indikators []domain.Indikator
	for _, indikatorReq := range request.Indikator {
		var targets []domain.Target

		// Proses target untuk setiap indikator
		for _, targetReq := range indikatorReq.Target {
			targetId := targetReq.Id
			if targetId == "" || targetId == "-" {
				// Generate ID baru jika ID kosong atau "-"
				targetId = generateTargetIdSasaran()
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorReq.Id,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       targetReq.Tahun,
			}
			targets = append(targets, target)
		}

		// Buat atau gunakan ID indikator yang ada
		indikatorId := indikatorReq.Id
		if indikatorId == "" {
			indikatorId = generateIndikatorIdSasaran()
		}

		indikator := domain.Indikator{
			Id:            indikatorId,
			TujuanPemdaId: sasaranPemda.TujuanPemdaId,
			Indikator:     indikatorReq.Indikator,
			RumusPerhitungan: sql.NullString{
				String: indikatorReq.RumusPerhitungan,
				Valid:  true,
			},
			SumberData: sql.NullString{
				String: indikatorReq.SumberData,
				Valid:  true,
			},
			Target: targets,
		}
		indikators = append(indikators, indikator)
	}

	sasaranPemda.Indikator = indikators

	// Simpan semua perubahan
	result, err := service.SasaranPemdaRepository.Update(ctx, tx, sasaranPemda)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	// Ambil data periode untuk response
	periode, err = service.PeriodeRepository.FindById(ctx, tx, sasaranPemda.PeriodeId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data periode: %v", err)
	}

	return sasaranpemda.SasaranPemdaResponse{
		Id:            result.Id,
		TujuanPemdaId: result.TujuanPemdaId,
		SubtemaId:     result.SubtemaId,
		NamaSubtema:   pokinData.NamaPohon,
		SasaranPemda:  result.SasaranPemda,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir,
		},
		Indikator: convertToIndikatorUpdateResponses(indikators),
	}, nil
}

func (service *SasaranPemdaServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.SasaranPemdaRepository.Delete(ctx, tx, id)
}

func (service *SasaranPemdaServiceImpl) FindById(ctx context.Context, sasaranPemdaId int) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data sasaran pemda
	sasaranPemda, err := service.SasaranPemdaRepository.FindById(ctx, tx, sasaranPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	// Ambil data pohon kinerja untuk nama subtema
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, sasaranPemda.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	// Ambil data tujuan pemda
	tujuanPemda, err := service.TujuanPemdaRepository.FindById(ctx, tx, sasaranPemda.TujuanPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data tujuan pemda: %v", err)
	}

	// Konversi indikator ke response
	var indikatorResponses []sasaranpemda.IndikatorResponse
	for _, indikator := range sasaranPemda.Indikator {
		var targetResponses []sasaranpemda.TargetResponse
		for _, target := range indikator.Target {
			targetResponses = append(targetResponses, sasaranpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, sasaranpemda.IndikatorResponse{
			Id:               indikator.Id,
			Indikator:        indikator.Indikator,
			RumusPerhitungan: indikator.RumusPerhitungan.String,
			SumberData:       indikator.SumberData.String,
			Target:           targetResponses,
		})
	}

	return sasaranpemda.SasaranPemdaResponse{
		Id:            sasaranPemda.Id,
		TujuanPemdaId: sasaranPemda.TujuanPemdaId,
		TujuanPemda:   tujuanPemda.TujuanPemda,
		SubtemaId:     sasaranPemda.SubtemaId,
		NamaSubtema:   pokinData.NamaPohon,
		SasaranPemda:  sasaranPemda.SasaranPemda,
		JenisPohon:    sasaranPemda.JenisPohon,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  sasaranPemda.Periode.TahunAwal,
			TahunAkhir: sasaranPemda.Periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}, nil
}

func (service *SasaranPemdaServiceImpl) FindAll(ctx context.Context, tahun string) ([]sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranPemdaList, err := service.SasaranPemdaRepository.FindAll(ctx, tx, tahun)
	if err != nil {
		return []sasaranpemda.SasaranPemdaResponse{}, err
	}

	sasaranPemdaResponses := make([]sasaranpemda.SasaranPemdaResponse, 0, len(sasaranPemdaList))
	for _, sasaranPemda := range sasaranPemdaList {
		pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, sasaranPemda.SubtemaId)
		if err != nil {
			return []sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
		}

		sasaranPemdaResponses = append(sasaranPemdaResponses, sasaranpemda.SasaranPemdaResponse{
			Id:           sasaranPemda.Id,
			SubtemaId:    sasaranPemda.SubtemaId,
			NamaSubtema:  pokinData.NamaPohon,
			SasaranPemda: sasaranPemda.SasaranPemda,
		})
	}

	return sasaranPemdaResponses, nil
}

func (service *SasaranPemdaServiceImpl) FindAllWithPokin(ctx context.Context, tahun string) ([]sasaranpemda.TematikResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	pokinData, err := service.SasaranPemdaRepository.FindAllWithPokin(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	// Konversi dari domain ke response
	var result []sasaranpemda.TematikResponse
	for _, tematik := range pokinData {
		tematikResponse := sasaranpemda.TematikResponse{
			TematikId:   tematik.TematikId,
			NamaTematik: tematik.NamaTematik,
			Subtematik:  []sasaranpemda.SubtematikResponse{},
		}

		// Konversi subtematik (hanya level 1-3)
		for _, subtematik := range tematik.Subtematik {
			subtematikResponse := sasaranpemda.SubtematikResponse{
				SubtematikId:   subtematik.SubtematikId,
				NamaSubtematik: subtematik.NamaSubtematik,
				JenisPohon:     subtematik.JenisPohon,
				LevelPohon:     subtematik.LevelPohon,
				SasaranPemda:   []sasaranpemda.SasaranPemdaWithpokinResponse{},
			}

			// Konversi sasaran pemda
			for _, sasaran := range subtematik.SasaranPemdaList {
				sasaranResponse := sasaranpemda.SasaranPemdaWithpokinResponse{
					IdSasaranPemda: sasaran.Id,
					SasaranPemda:   sasaran.SasaranPemda,
					Periode: sasaranpemda.PeriodeResponse{
						TahunAwal:  sasaran.Periode.TahunAwal,
						TahunAkhir: sasaran.Periode.TahunAkhir,
					},
					Indikator: []sasaranpemda.IndikatorSubtematikResponse{},
				}

				// Konversi indikator
				for _, indikator := range sasaran.Indikator {
					indikatorResponse := sasaranpemda.IndikatorSubtematikResponse{
						Id:               indikator.Id,
						Indikator:        indikator.Indikator,
						RumusPerhitungan: indikator.RumusPerhitungan.String,
						SumberData:       indikator.SumberData.String,
						Target:           []sasaranpemda.TargetResponse{},
					}

					// Konversi target
					for _, target := range indikator.Target {
						targetResponse := sasaranpemda.TargetResponse{
							Id:     target.Id,
							Target: target.Target,
							Satuan: target.Satuan,
							Tahun:  target.Tahun,
						}
						indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
					}

					// Sort target berdasarkan tahun
					sort.Slice(indikatorResponse.Target, func(i, j int) bool {
						return indikatorResponse.Target[i].Tahun < indikatorResponse.Target[j].Tahun
					})

					sasaranResponse.Indikator = append(sasaranResponse.Indikator, indikatorResponse)
				}

				// Sort indikator berdasarkan ID
				sort.Slice(sasaranResponse.Indikator, func(i, j int) bool {
					return sasaranResponse.Indikator[i].Id < sasaranResponse.Indikator[j].Id
				})

				subtematikResponse.SasaranPemda = append(subtematikResponse.SasaranPemda, sasaranResponse)
			}

			// Tambahkan subtematik jika memiliki sasaran pemda atau level 1-3
			if len(subtematikResponse.SasaranPemda) > 0 || (subtematik.LevelPohon >= 1 && subtematik.LevelPohon <= 3) {
				tematikResponse.Subtematik = append(tematikResponse.Subtematik, subtematikResponse)
			}
		}

		// Sort subtematik berdasarkan level dan ID
		sort.Slice(tematikResponse.Subtematik, func(i, j int) bool {
			if tematikResponse.Subtematik[i].LevelPohon != tematikResponse.Subtematik[j].LevelPohon {
				return tematikResponse.Subtematik[i].LevelPohon < tematikResponse.Subtematik[j].LevelPohon
			}
			return tematikResponse.Subtematik[i].SubtematikId < tematikResponse.Subtematik[j].SubtematikId
		})

		// Tambahkan tematik jika memiliki subtematik
		if len(tematikResponse.Subtematik) > 0 {
			result = append(result, tematikResponse)
		}
	}

	// Sort tematik berdasarkan nama
	sort.Slice(result, func(i, j int) bool {
		return result[i].NamaTematik < result[j].NamaTematik
	})

	return result, nil
}

func convertToIndikatorUpdateResponses(indikators []domain.Indikator) []sasaranpemda.IndikatorResponse {
	if len(indikators) == 0 {
		return nil
	}

	responses := make([]sasaranpemda.IndikatorResponse, len(indikators))
	for i, indikator := range indikators {
		targetResponses := make([]sasaranpemda.TargetResponse, len(indikator.Target))
		for j, target := range indikator.Target {
			targetResponses[j] = sasaranpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			}
		}

		// Sort target berdasarkan tahun
		sort.Slice(targetResponses, func(i, j int) bool {
			return targetResponses[i].Tahun < targetResponses[j].Tahun
		})

		responses[i] = sasaranpemda.IndikatorResponse{
			Id:               indikator.Id,
			Indikator:        indikator.Indikator,
			RumusPerhitungan: indikator.RumusPerhitungan.String,
			SumberData:       indikator.SumberData.String,
			Target:           targetResponses,
		}
	}

	// Sort indikator berdasarkan ID
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Id < responses[j].Id
	})

	return responses
}
