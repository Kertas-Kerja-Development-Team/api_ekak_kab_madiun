CREATE TABLE tb_urusan (
    id VARCHAR(225) NOT NULL,
    kode_urusan VARCHAR(225),
    nama_urusan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE = InnoDB;
