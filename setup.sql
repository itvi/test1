-- department
CREATE TABLE department(
    id INTEGER PRIMARY KEY autoincrement,
    code VARCHAR(7) NOT NULL UNIQUE,
    name VARCHAR(20)
);
CREATE TABLE asset_category(
    id INTEGER PRIMARY KEY autoincrement,
    code VARCHAR(4) NOT NULL UNIQUE,
    name VARCHAR(10) NOT NULL UNIQU
);
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
-- assets movement document
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
