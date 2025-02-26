CREATE TABLE tb_visi_pemda (
    id INT AUTO_INCREMENT PRIMARY KEY,
    visi TEXT,
    tahun_awal_periode VARCHAR(255) NOT NULL,
    tahun_akhir_periode VARCHAR(255) NOT NULL,
    jenis_periode VARCHAR(255) NOT NULL,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB;