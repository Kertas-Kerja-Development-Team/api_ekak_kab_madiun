ALTER TABLE tb_subkegiatan 
ADD COLUMN rekin_id VARCHAR(255) NOT NULL,
ADD COLUMN status VARCHAR(255) NOT NULL DEFAULT 'belum_diambil',
DROP COLUMN pegawai_id;