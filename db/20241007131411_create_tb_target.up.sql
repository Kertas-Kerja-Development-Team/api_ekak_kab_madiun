CREATE TABLE tb_target (
    id VARCHAR(255) UNIQUE,
    indikator_id VARCHAR(255),
    target INT,
    satuan VARCHAR(255),
    tahun VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB;