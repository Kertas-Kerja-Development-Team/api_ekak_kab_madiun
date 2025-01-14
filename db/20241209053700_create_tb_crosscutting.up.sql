CREATE TABLE tb_crosscutting (
    id INT AUTO_INCREMENT PRIMARY KEY,
    crosscutting_from INT,
    crosscutting_to INT UNIQUE,
    status VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
