package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type SasaranPemdaRepositoryImpl struct {
}

func NewSasaranPemdaRepositoryImpl() *SasaranPemdaRepositoryImpl {
	return &SasaranPemdaRepositoryImpl{}
}

func (repository *SasaranPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error) {
	query := "INSERT INTO tb_sasaran_pemda(id, sasaran_pemda_id, rumus_perhitungan, sumber_data, periode_id) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.Id, sasaranPemda.SasaranPemdaId, sasaranPemda.RumusPerhitungan, sasaranPemda.SumberData, sasaranPemda.PeriodeId)
	if err != nil {
		return sasaranPemda, err
	}
	return sasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) CreateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) (domain.Indikator, error) {
	query := "INSERT INTO tb_indikator(id, sasaran_pemda_id, indikator) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, indikator.Id, indikator.SasaranPemdaId, indikator.Indikator)
	if err != nil {
		return indikator, err
	}
	return indikator, nil
}

func (repository *SasaranPemdaRepositoryImpl) CreateTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error {
	query := "INSERT INTO tb_target(id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, target.Id, target.IndikatorId, target.Target, target.Satuan, target.Tahun)
	return err
}

func (repository *SasaranPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error) {
	// Update tujuan pemda
	query := "UPDATE tb_sasaran_pemda SET sasaran_pemda_id = ?, rumus_perhitungan = ?, sumber_data = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.SasaranPemdaId, sasaranPemda.RumusPerhitungan, sasaranPemda.SumberData, sasaranPemda.Id)
	if err != nil {
		return sasaranPemda, err
	}

	// Proses indikator
	for _, indikator := range sasaranPemda.Indikator {
		// Update atau insert indikator
		scriptUpdateIndikator := `
            INSERT INTO tb_indikator (id, sasaran_pemda_id, indikator) 
            VALUES (?, ?, ?)
            ON DUPLICATE KEY UPDATE indikator = VALUES(indikator)`

		_, err := tx.ExecContext(ctx, scriptUpdateIndikator,
			indikator.Id,
			sasaranPemda.Id,
			indikator.Indikator)
		if err != nil {
			return sasaranPemda, err
		}

		// Hapus target lama untuk indikator ini
		scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTargets, indikator.Id)
		if err != nil {
			return sasaranPemda, err
		}

		// Insert target baru
		for _, target := range indikator.Target {
			scriptInsertTarget := `
                INSERT INTO tb_target 
                    (id, indikator_id, target, satuan, tahun)
                VALUES 
                    (?, ?, ?, ?, ?)`

			_, err := tx.ExecContext(ctx, scriptInsertTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return sasaranPemda, err
			}
		}
	}

	return sasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) error {
	// 1. Hapus target berdasarkan indikator yang terkait dengan tujuan pemda
	queryDeleteTarget := `
        DELETE t FROM tb_target t
        INNER JOIN tb_indikator i ON t.indikator_id = i.id
        WHERE i.sasaran_pemda_id = ?`
	_, err := tx.ExecContext(ctx, queryDeleteTarget, sasaranPemdaId)
	if err != nil {
		return err
	}

	// 2. Hapus indikator yang terkait dengan tujuan pemda
	queryDeleteIndikator := "DELETE FROM tb_indikator WHERE sasaran_pemda_id = ?"
	_, err = tx.ExecContext(ctx, queryDeleteIndikator, sasaranPemdaId)
	if err != nil {
		return err
	}

	// 3. Hapus tujuan pemda
	queryDeleteSasaranPemda := "DELETE FROM tb_sasaran_pemda WHERE id = ?"
	_, err = tx.ExecContext(ctx, queryDeleteSasaranPemda, sasaranPemdaId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *SasaranPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) (domain.SasaranPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.sasaran_pemda_id,
            tp.rumus_perhitungan,
            tp.sumber_data,
            tp.periode_id,
            COALESCE(p.tahun_awal, '') as tahun_awal,
            COALESCE(p.tahun_akhir, '') as tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_sasaran_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_indikator i ON tp.id = i.sasaran_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        WHERE tp.id = ?
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query, sasaranPemdaId)
	if err != nil {
		return domain.SasaranPemda{}, err
	}
	defer rows.Close()

	var sasaranPemda domain.SasaranPemda
	var firstRow = true
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var id int
		var sasaranPemdaId int
		var rumusPerhitungan, sumberData string
		var periodeId int
		var tahunAwal, tahunAkhir string
		var indikatorId, indikatorText sql.NullString
		var targetId, targetValue, targetSatuan, targetTahun sql.NullString

		err := rows.Scan(
			&id,
			&sasaranPemdaId,
			&rumusPerhitungan,
			&sumberData,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return domain.SasaranPemda{}, err
		}

		if firstRow {
			sasaranPemda = domain.SasaranPemda{
				Id:               id,
				SasaranPemdaId:   sasaranPemdaId,
				PeriodeId:        periodeId,
				RumusPerhitungan: rumusPerhitungan,
				SumberData:       sumberData,
				Periode: domain.Periode{
					TahunAwal:  tahunAwal,
					TahunAkhir: tahunAkhir,
				},
				Indikator: []domain.Indikator{},
			}

			// Update periode hanya jika ada dan valid
			if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
				sasaranPemda.Periode.TahunAwal = tahunAwal
				sasaranPemda.Periode.TahunAkhir = tahunAkhir
			}

			firstRow = false
		}

		if indikatorId.Valid && indikatorText.Valid {
			currentIndikator, exists := indikatorMap[indikatorId.String]

			if !exists {
				indikator := domain.Indikator{
					Id:             indikatorId.String,
					SasaranPemdaId: id,
					Indikator:      indikatorText.String,
					Target:         make([]domain.Target, 0),
				}

				// Tambahkan target kosong hanya jika periode valid
				if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
					tahunAwalInt, _ := strconv.Atoi(tahunAwal)
					tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

					for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						indikator.Target = append(indikator.Target, domain.Target{
							Id:          "",
							IndikatorId: indikatorId.String,
							Target:      "-",
							Satuan:      "-",
							Tahun:       tahunStr,
						})
					}
				}

				sasaranPemda.Indikator = append(sasaranPemda.Indikator, indikator)
				indikatorMap[indikatorId.String] = &sasaranPemda.Indikator[len(sasaranPemda.Indikator)-1]
				currentIndikator = &sasaranPemda.Indikator[len(sasaranPemda.Indikator)-1]
			}

			// Update target hanya jika periode valid
			if periodeId != 0 && targetId.Valid && targetValue.Valid && targetTahun.Valid {
				// Update target yang sudah ada dengan data sebenarnya
				for i := range currentIndikator.Target {
					if currentIndikator.Target[i].Tahun == targetTahun.String {
						currentIndikator.Target[i] = domain.Target{
							Id:          targetId.String,
							IndikatorId: indikatorId.String,
							Target:      targetValue.String,
							Satuan:      targetSatuan.String,
							Tahun:       targetTahun.String,
						}
						break
					}
				}
			}
		}
	}

	if sasaranPemda.Id == 0 {
		return sasaranPemda, errors.New("sasaran pemda not found")
	}

	return sasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.SasaranPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.sasaran_pemda_id,
            tp.rumus_perhitungan,
            tp.sumber_data,
            tp.periode_id,
             COALESCE(p.tahun_awal, '') as tahun_awal,
            COALESCE(p.tahun_akhir, '') as tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_sasaran_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_indikator i ON tp.id = i.sasaran_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sasaranPemdaMap := make(map[int]*domain.SasaranPemda)

	for rows.Next() {
		var id int
		var sasaranPemdaId int
		var rumusPerhitungan, sumberData string
		var periodeId int
		var tahunAwal, tahunAkhir string
		var indikatorId, indikatorText sql.NullString
		var targetId, targetValue, targetSatuan, targetTahun sql.NullString

		err := rows.Scan(
			&id,
			&sasaranPemdaId,
			&rumusPerhitungan,
			&sumberData,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return nil, err
		}

		sasaranPemda, exists := sasaranPemdaMap[id]
		if !exists {
			sasaranPemda = &domain.SasaranPemda{
				Id:               id,
				SasaranPemdaId:   sasaranPemdaId,
				PeriodeId:        periodeId,
				RumusPerhitungan: rumusPerhitungan,
				SumberData:       sumberData,
				Periode: domain.Periode{
					TahunAwal:  tahunAwal,
					TahunAkhir: tahunAkhir,
				},
				Indikator: []domain.Indikator{},
			}

			// Update periode hanya jika ada dan valid
			if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
				sasaranPemda.Periode.TahunAwal = tahunAwal
				sasaranPemda.Periode.TahunAkhir = tahunAkhir
			}

			sasaranPemdaMap[id] = sasaranPemda
		}

		if indikatorId.Valid && indikatorText.Valid {
			var currentIndikator *domain.Indikator

			// Cari indikator yang sudah ada
			for i := range sasaranPemda.Indikator {
				if sasaranPemda.Indikator[i].Id == indikatorId.String {
					currentIndikator = &sasaranPemda.Indikator[i]
					break
				}
			}

			// Buat indikator baru jika belum ada
			if currentIndikator == nil {
				newIndikator := domain.Indikator{
					Id:             indikatorId.String,
					SasaranPemdaId: id,
					Indikator:      indikatorText.String,
					Target:         make([]domain.Target, 0),
				}

				// Tambahkan target kosong untuk semua tahun dalam periode
				if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
					tahunAwalInt, _ := strconv.Atoi(tahunAwal)
					tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

					for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						newIndikator.Target = append(newIndikator.Target, domain.Target{
							Id:          "",
							IndikatorId: indikatorId.String,
							Target:      "-",
							Satuan:      "-",
							Tahun:       tahunStr,
						})
					}
				}

				sasaranPemda.Indikator = append(sasaranPemda.Indikator, newIndikator)
				currentIndikator = &sasaranPemda.Indikator[len(sasaranPemda.Indikator)-1]
			}

			// Update target dengan data sebenarnya jika ada
			if targetId.Valid && targetValue.Valid && targetTahun.Valid {
				for i := range currentIndikator.Target {
					if currentIndikator.Target[i].Tahun == targetTahun.String {
						currentIndikator.Target[i] = domain.Target{
							Id:          targetId.String,
							IndikatorId: indikatorId.String,
							Target:      targetValue.String,
							Satuan:      targetSatuan.String,
							Tahun:       targetTahun.String,
						}
						break
					}
				}
			}
		}
	}

	// Convert map to slice dan urutkan hasilnya
	result := make([]domain.SasaranPemda, 0, len(sasaranPemdaMap))
	for _, sasaranPemda := range sasaranPemdaMap {
		result = append(result, *sasaranPemda)
	}

	// Urutkan berdasarkan ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})

	return result, nil
}

func (repository *SasaranPemdaRepositoryImpl) DeleteIndikator(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) error {
	query := "DELETE FROM tb_indikator WHERE sasaran_pemda_id = ?"
	_, err := tx.ExecContext(ctx, query, sasaranPemdaId)
	return err
}

func (repository *SasaranPemdaRepositoryImpl) IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool {
	query := "SELECT COUNT(*) FROM tb_sasaran_pemda WHERE id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return true
	}
	return count > 0
}

func (repository *SasaranPemdaRepositoryImpl) UpdatePeriode(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error) {
	// Update hanya periode_id
	query := "UPDATE tb_sasaran_pemda SET periode_id = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, sasaranPemda.PeriodeId, sasaranPemda.Id)
	if err != nil {
		return domain.SasaranPemda{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.SasaranPemda{}, err
	}

	if rowsAffected == 0 {
		return domain.SasaranPemda{}, fmt.Errorf("periode dengan id %d sudah digunakan", sasaranPemda.PeriodeId)
	}

	// Ambil data terbaru setelah update
	query = `
        SELECT 
            tp.id,
            tp.sasaran_pemda_id,
            tp.rumus_perhitungan,
            tp.sumber_data,
            tp.periode_id,
            COALESCE(p.tahun_awal, 'Pilih periode') as tahun_awal,
            COALESCE(p.tahun_akhir, 'Pilih periode') as tahun_akhir
        FROM 
            tb_sasaran_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
        WHERE tp.id = ?`

	var updatedSasaranPemda domain.SasaranPemda
	err = tx.QueryRowContext(ctx, query, sasaranPemda.Id).Scan(
		&updatedSasaranPemda.Id,
		&updatedSasaranPemda.SasaranPemdaId,
		&updatedSasaranPemda.RumusPerhitungan,
		&updatedSasaranPemda.SumberData,
		&updatedSasaranPemda.PeriodeId,
		&updatedSasaranPemda.Periode.TahunAwal,
		&updatedSasaranPemda.Periode.TahunAkhir,
	)
	if err != nil {
		return domain.SasaranPemda{}, fmt.Errorf("gagal mengambil data setelah update: %v", err)
	}

	return updatedSasaranPemda, nil
}
