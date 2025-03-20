package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MatrixRenstraRepositoryImpl struct{}

func NewMatrixRenstraRepositoryImpl() *MatrixRenstraRepositoryImpl {
	return &MatrixRenstraRepositoryImpl{}
}

func (repository *MatrixRenstraRepositoryImpl) GetByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.SubKegiatanQuery, error) {
	query := `
        WITH RECURSIVE hierarchy AS (
            SELECT DISTINCT
                u.kode_urusan,
                u.nama_urusan,
                bu.kode_bidang_urusan,
                bu.nama_bidang_urusan,
                p.kode_program,
                p.nama_program,
                k.kode_kegiatan,
                k.nama_kegiatan,
                s.kode_subkegiatan,
                s.nama_subkegiatan,
                so.tahun as tahun_subkegiatan
            FROM tb_subkegiatan_opd so
            JOIN tb_subkegiatan s ON so.kode_subkegiatan = s.kode_subkegiatan
            JOIN tb_master_kegiatan k ON LEFT(s.kode_subkegiatan, LENGTH(k.kode_kegiatan)) = k.kode_kegiatan
            JOIN tb_master_program p ON LEFT(k.kode_kegiatan, LENGTH(p.kode_program)) = p.kode_program
            JOIN tb_bidang_urusan bu ON LEFT(p.kode_program, LENGTH(bu.kode_bidang_urusan)) = bu.kode_bidang_urusan
            JOIN tb_urusan u ON LEFT(bu.kode_bidang_urusan, LENGTH(u.kode_urusan)) = u.kode_urusan
            WHERE so.kode_opd = ?
        )
        SELECT 
            h.kode_urusan,
            h.nama_urusan,
            h.kode_bidang_urusan,
            h.nama_bidang_urusan,
            h.kode_program,
            h.nama_program,
            h.kode_kegiatan,
            h.nama_kegiatan,
            h.kode_subkegiatan,
            h.nama_subkegiatan,
            h.tahun_subkegiatan,
            i.id as indikator_id,
            i.kode as indikator_kode,
            i.indikator,
            i.tahun as indikator_tahun,
            i.kode_opd as indikator_kode_opd,
            t.id,
            t.target,
            t.satuan
        FROM hierarchy h
        LEFT JOIN tb_indikator i ON 
            (
                i.kode = h.kode_urusan OR 
                i.kode = h.kode_bidang_urusan OR
                i.kode = h.kode_program OR
                i.kode = h.kode_kegiatan OR
                i.kode = h.kode_subkegiatan
            )
            AND i.kode_opd = ?
            AND i.tahun BETWEEN ? AND ?
        LEFT JOIN tb_target t ON t.indikator_id = i.id
        ORDER BY 
            h.kode_urusan,
            h.kode_bidang_urusan,
            h.kode_program,
            h.kode_kegiatan,
            h.kode_subkegiatan,
            i.tahun
    `

	rows, err := tx.QueryContext(ctx, query, kodeOpd, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.SubKegiatanQuery
	for rows.Next() {
		var data domain.SubKegiatanQuery
		var indikatorId, indikatorKode, indikator, indikatorTahun, indikatorKodeOpd, targetId, target, satuan sql.NullString

		err := rows.Scan(
			&data.KodeUrusan,
			&data.NamaUrusan,
			&data.KodeBidangUrusan,
			&data.NamaBidangUrusan,
			&data.KodeProgram,
			&data.NamaProgram,
			&data.KodeKegiatan,
			&data.NamaKegiatan,
			&data.KodeSubKegiatan,
			&data.NamaSubKegiatan,
			&data.TahunSubKegiatan,
			&indikatorId,
			&indikatorKode,
			&indikator,
			&indikatorTahun,
			&indikatorKodeOpd,
			&targetId,
			&target,
			&satuan,
		)
		if err != nil {
			return nil, err
		}

		// Handle null values
		if indikatorId.Valid {
			data.IndikatorId = indikatorId.String
			data.IndikatorKode = indikatorKode.String
			data.Indikator = indikator.String
			data.IndikatorTahun = indikatorTahun.String
			data.IndikatorKodeOpd = indikatorKodeOpd.String
		}
		if target.Valid {
			data.TargetId = targetId.String
			data.Target = target.String
			data.Satuan = satuan.String
		}

		result = append(result, data)
	}

	return result, nil
}
