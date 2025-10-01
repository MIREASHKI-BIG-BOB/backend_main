-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS medicals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL,
    address VARCHAR(100) NOT NULL,
    phone VARCHAR(11) NOT NULL,
    email VARCHAR(30) NOT NULL,
    license_number INTEGER NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS doctors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(15) NOT NULL,
    phone VARCHAR(11) NOT NULL,
    specialization VARCHAR(20) NOT NULL,
    license_number INTEGER NOT NULL UNIQUE,
    med_id INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (med_id) REFERENCES medicals(id)
    );

CREATE TABLE IF NOT EXISTS examinations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    med_id INTEGER NOT NULL,
    doctor_id INTEGER NOT NULL,
    notes VARCHAR(300),
    status INTEGER DEFAULT 0,
    cloud_id INTEGER,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_by INTEGER NOT NULL,
    FOREIGN KEY (med_id) REFERENCES medicals(id),
    FOREIGN KEY (doctor_id) REFERENCES doctors(id),
    FOREIGN KEY (created_by) REFERENCES doctors(id),
    FOREIGN KEY (updated_by) REFERENCES doctors(id)
    );

CREATE TABLE IF NOT EXISTS ctg (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    examination_id INTEGER NOT NULL,
    sec_from_start REAL NOT NULL DEFAULT 0.0,
    uuid VARCHAR(36) NOT NULL,
    bpm REAL DEFAULT 0.0,
    uterus REAL DEFAULT 0.0,
    spasms REAL DEFAULT 0.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (examination_id) REFERENCES examinations(id)
    );

CREATE INDEX IF NOT EXISTS idx_examinations_client_id ON examinations(client_id);
CREATE INDEX IF NOT EXISTS idx_examinations_doctor_id ON examinations(doctor_id);
CREATE INDEX IF NOT EXISTS idx_examinations_med_id ON examinations(med_id);
CREATE INDEX IF NOT EXISTS idx_examinations_start_time ON examinations(start_time);
CREATE INDEX IF NOT EXISTS idx_ctg_examination_id ON ctg(examination_id);
CREATE INDEX IF NOT EXISTS idx_ctg_uuid ON ctg(uuid);
CREATE INDEX IF NOT EXISTS idx_ctg_created_at ON ctg(created_at);
CREATE INDEX IF NOT EXISTS idx_doctors_med_id ON doctors(med_id);
CREATE INDEX IF NOT EXISTS idx_doctors_is_active ON doctors(is_active);

INSERT INTO medicals (name, address, phone, email, license_number) VALUES
    ('Чухановская больница №1', 'ул. Медицинская, 15', '84956756975', 'hospital1@med.ru', 12345);

INSERT INTO doctors (name, phone, specialization, license_number, med_id) VALUES
    ('Стринченко Кеша', '79161234567', 'Акушер-гинеколог', 67890, 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_doctors_is_active;
DROP INDEX IF EXISTS idx_doctors_med_id;
DROP INDEX IF EXISTS idx_ctg_created_at;
DROP INDEX IF EXISTS idx_ctg_uuid;
DROP INDEX IF EXISTS idx_ctg_examination_id;
DROP INDEX IF EXISTS idx_examinations_start_time;
DROP INDEX IF EXISTS idx_examinations_med_id;
DROP INDEX IF EXISTS idx_examinations_doctor_id;
DROP INDEX IF EXISTS idx_examinations_client_id;

DROP TABLE IF EXISTS ctg;
DROP TABLE IF EXISTS examinations;
DROP TABLE IF EXISTS doctors;
DROP TABLE IF EXISTS medicals;

-- +goose StatementEnd
