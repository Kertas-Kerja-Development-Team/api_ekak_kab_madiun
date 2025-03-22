package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
)

type ProgramRepositoryImpl struct {
}

func NewProgramRepositoryImpl() *ProgramRepositoryImpl {
	return &ProgramRepositoryImpl{}
}

func (repository *ProgramRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, program domainmaster.ProgramKegiatan) (domainmaster.ProgramKegiatan, error) {
	scriptProgram := "INSERT INTO tb_master_program (id, kode_program, nama_program, tahun, is_active) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, scriptProgram, program.Id, program.KodeProgram, program.NamaProgram, program.Tahun, program.IsActive)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, err
	}
	scriptIndikator := "INSERT INTO tb_indikator (id, program_id, indikator, tahun) VALUES (?, ?, ?, ?)"
	for _, indikator := range program.Indikator {
		_, err = tx.ExecContext(ctx, scriptIndikator, indikator.Id, indikator.ProgramId, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domainmaster.ProgramKegiatan{}, err
		}

		scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
		for _, target := range indikator.Target {
			_, err = tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domainmaster.ProgramKegiatan{}, err
			}
		}
	}
	return program, nil
}

func (repository *ProgramRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, program domainmaster.ProgramKegiatan) (domainmaster.ProgramKegiatan, error) {
	scriptProgram := "UPDATE tb_master_program SET kode_program = ?, nama_program = ?, tahun = ?, is_active = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, scriptProgram, program.KodeProgram, program.NamaProgram, program.Tahun, program.IsActive, program.Id)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE program_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, program.Id)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, err
	}

	scriptIndikator := "INSERT INTO tb_indikator (id, program_id, indikator, tahun) VALUES (?, ?, ?, ?)"
	for _, indikator := range program.Indikator {
		_, err = tx.ExecContext(ctx, scriptIndikator, indikator.Id, indikator.ProgramId, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domainmaster.ProgramKegiatan{}, err
		}

		scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTarget, indikator.Id)
		if err != nil {
			return domainmaster.ProgramKegiatan{}, err
		}

		scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
		for _, target := range indikator.Target {
			_, err = tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domainmaster.ProgramKegiatan{}, err
			}
		}
	}

	return program, nil
}

func (repository *ProgramRepositoryImpl) FindIndikatorByProgramId(ctx context.Context, tx *sql.Tx, programId string) ([]domain.Indikator, error) {
	script := "SELECT id, program_id, indikator, tahun FROM tb_indikator WHERE program_id = ?"
	params := []interface{}{programId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.ProgramId, &indikator.Indikator, &indikator.Tahun)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *ProgramRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	params := []interface{}{indikatorId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (repository *ProgramRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.ProgramKegiatan, error) {
	// Query untuk mendapatkan data program
	scriptProgram := "SELECT id, kode_program, nama_program, tahun, is_active FROM tb_master_program WHERE id = ?"
	rows, err := tx.QueryContext(ctx, scriptProgram, id)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, fmt.Errorf("gagal mengambil data program: %v", err)
	}
	defer rows.Close()

	var program domainmaster.ProgramKegiatan
	if !rows.Next() {
		return domainmaster.ProgramKegiatan{}, fmt.Errorf("program dengan id %s tidak ditemukan", id)
	}

	// Scan data program
	err = rows.Scan(
		&program.Id,
		&program.KodeProgram,
		&program.NamaProgram,
		&program.Tahun,
		&program.IsActive,
	)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, fmt.Errorf("gagal scan data program: %v", err)
	}

	return program, nil
}

func (repository *ProgramRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	scriptIndikator := "DELETE FROM tb_indikator WHERE program_id = ?"
	_, err := tx.ExecContext(ctx, scriptIndikator, id)
	if err != nil {
		return err
	}

	scriptProgram := "DELETE FROM tb_master_program WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptProgram, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *ProgramRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.ProgramKegiatan, error) {
	script := "SELECT id, kode_program, nama_program, tahun, is_active FROM tb_master_program ORDER BY id"

	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []domainmaster.ProgramKegiatan
	for rows.Next() {
		var program domainmaster.ProgramKegiatan
		err := rows.Scan(
			&program.Id,
			&program.KodeProgram,
			&program.NamaProgram,
			&program.Tahun,
			&program.IsActive,
		)
		if err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}

	return programs, nil
}
