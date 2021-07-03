CREATE TABLE app_user (
    id INTEGER PRIMARY KEY autoincrement,
    sn VARCHAR(8) NOT NULL UNIQUE,
    name VARCHAR(6),
    email VARCHAR(50),
    hashed_password CHAR(60),
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
-- department
CREATE TABLE department(
    id INTEGER PRIMARY KEY autoincrement,
    code VARCHAR(7) NOT NULL UNIQUE,
    name VARCHAR(20)
);
-- Asset ...
CREATE TABLE asset(
    id INTEGER PRIMARY KEY autoincrement,
    number VARCHAR(9) NOT NULL UNIQUE,
    category_code VARCHAR(4),
    unit VARCHAR(1),
    supplier varchar(50),
    model VARCHAR(30),
    sn VARCHAR(30),
    warranty date,
    remark VARCHAR(100),
    created datetime DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE asset_category(
    id INTEGER PRIMARY KEY autoincrement,
    code VARCHAR(4) NOT NULL UNIQUE,
    name VARCHAR(10) NOT NULL UNIQU
);
CREATE TABLE asset_status(
    id INTEGER PRIMARY KEY autoincrement,
    name VARCHAR(2) NOT NULL UNIQUE
);
--- assets movement document
CREATE TABLE asset_mov(
    id INTEGER PRIMARY KEY autoincrement,
    a_number VARCHAR(9),
    mvt VARCHAR(3),
    qty int,
    from_loc varchar(15),
    to_loc varchar(15),
    from_employee VARCHAR(7),
    to_employee VARCHAR(7),
    doc_date date,
    created datetime DEFAULT(DATETIME ('now', 'localtime'))
);
-- computer config information
CREATE TABLE computer_cs (
    id INTEGER PRIMARY KEY autoincrement,
    asset_number VARCHAR(9),
    name VARCHAR(50),
    user_name VARCHAR(50),
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE computer_os(
    id INTEGER PRIMARY KEY autoincrement,
    asset_number VARCHAR(9),
    caption VARCHAR(80),
    version VARCHAR(20),
    install_date TEXT,
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE computer_cpu (
    id INTEGER PRIMARY KEY autoincrement,
    asset_number VARCHAR(9),
    name TEXT,
    number_of_cores INTEGER,
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE computer_disk(
    id INTEGER PRIMARY KEY autoincrement,
    asset_number VARCHAR(9),
    model VARCHAR(60),
    size FLOAT,
    sn VARCHAR(60),
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE computer_mem(
    id INTEGER PRIMARY KEY autoincrement,
    asset_number VARCHAR(9),
    manufacturer VARCHAR(20),
    capacity FLOAT,
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE computer_net(
    id INTEGER PRIMARY KEY autoincrement,
    asset_number VARCHAR(9),
    description VARCHAR(50),
    mac VARCHAR(17),
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
-- init computer config information
CREATE TABLE init_computer_cs (
    id INTEGER PRIMARY KEY autoincrement,
    ip VARCHAR(15),
    name VARCHAR(50),
    user_name VARCHAR(50),
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE init_computer_os(
    id INTEGER PRIMARY KEY autoincrement,
    ip VARCHAR(15),
    caption VARCHAR(80),
    version VARCHAR(20),
    install_date TEXT,
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE init_computer_cpu (
    id INTEGER PRIMARY KEY autoincrement,
    ip VARCHAR(15),
    name TEXT,
    number_of_cores INTEGER,
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE init_computer_disk(
    id INTEGER PRIMARY KEY autoincrement,
    ip VARCHAR(15),
    model VARCHAR(60),
    size FLOAT,
    sn VARCHAR(60),
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE init_computer_mem(
    id INTEGER PRIMARY KEY autoincrement,
    ip VARCHAR(15),
    manufacturer VARCHAR(20),
    capacity FLOAT,
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);
CREATE TABLE init_computer_net(
    id INTEGER PRIMARY KEY autoincrement,
    ip VARCHAR(15),
    description VARCHAR(50),
    mac VARCHAR(17),
    created TIMESTAMP DEFAULT(DATETIME ('now', 'localtime'))
);