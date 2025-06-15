-- Extensions.
CREATE EXTENSION postgis;
CREATE EXTENSION pgrouting;

-- Functions.
CREATE FUNCTION enforce_lower_case_username() 
RETURNS TRIGGER AS $$
BEGIN
    new.username = LOWER(new.username);
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION update_modified_at() 
RETURNS TRIGGER AS $$
BEGIN
    new.modified_at = CURRENT_TIMESTAMP;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

-- Users.
CREATE TABLE users (
    id          uuid        NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    username    varchar(50) NOT NULL,
    password    varchar(60) NOT NULL,
    first_name  varchar(50) NOT NULL,
    last_name   varchar(50) NOT NULL,
    created_at  timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey           PRIMARY KEY (id),
    CONSTRAINT users_username_key   UNIQUE (username)
);

CREATE TRIGGER users_enforce_lower_case_username
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE enforce_lower_case_username();

CREATE TRIGGER users_update_modified_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();

-- Employees.
CREATE TYPE employees_role AS ENUM ('waste_operator', 'manager');

CREATE TABLE employees (
    id              uuid                    NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    username        varchar(50)             NOT NULL,
    password        varchar(60)             NOT NULL,
    first_name      varchar(50)             NOT NULL,
    last_name       varchar(50)             NOT NULL,
    role            employees_role          NOT NULL,
    date_of_birth   date                    NOT NULL,
    phone_number    varchar(20)             NOT NULL,
    geom            geometry('POINT', 4326) NOT NULL,
    schedule_start  time                    NOT NULL,
    schedule_end    time                    NOT NULL,
    created_at      timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT employees_pkey           PRIMARY KEY (id),
    CONSTRAINT employees_username_key   UNIQUE (username)
);

CREATE TRIGGER employees_enforce_lower_case_username
    BEFORE INSERT OR UPDATE ON employees
    FOR EACH ROW
    EXECUTE PROCEDURE enforce_lower_case_username();

CREATE TRIGGER employees_update_modified_at
    BEFORE UPDATE ON employees
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();

-- Insert default admin manager employee with 'admin' password.
INSERT INTO employees(username, password, first_name, last_name, role, date_of_birth, phone_number, geom, schedule_start, schedule_end)
    VALUES ('admin', '$2a$14$3YbWepKf4uralK8RS4Gi5eiIwKOrUU0dHWXaTCGpEmWfr.gdjEu96', 'Super', 'Admin', 'manager', '1970-01-01', '', ST_GeomFromGeoJSON('{"coordinates": [0, 0], "type": "Point"}'), '00:00:00', '00:00:00');

-- Containers.
CREATE TYPE containers_category AS ENUM ('general', 'paper', 'plastic', 'metal', 'glass', 'organic', 'hazardous');

CREATE TABLE containers (
    id          uuid                    NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    category    containers_category     NOT NULL,
    geom        geometry('POINT', 4326) NOT NULL,
    created_at  timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT containers_pkey  PRIMARY KEY (id)
);

CREATE TRIGGER containers_update_modified_at
    BEFORE UPDATE ON containers
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();

-- Container reports.
CREATE TYPE containers_reports_issue_type AS ENUM ('full', 'vandalized', 'misplaced', 'non-existent', 'other');

CREATE TABLE containers_reports (
    id              uuid                            NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    container_id    uuid                            NOT NULL,
    issue_type      containers_reports_issue_type   NOT NULL,
    description     varchar(500),
    attachment      bytea,
    issuer_id       uuid,
    resolver_id     uuid,
    resolved        boolean                         NOT NULL    DEFAULT FALSE,
    created_at      timestamp                       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp                       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT containers_reports_pkey              PRIMARY KEY (id),
    CONSTRAINT containers_reports_container_id_fkey FOREIGN KEY (container_id)  REFERENCES containers (id)  ON DELETE CASCADE,
    CONSTRAINT containers_reports_issuer_id_fkey    FOREIGN KEY (issuer_id)     REFERENCES users (id)       ON DELETE SET NULL,
    CONSTRAINT containers_reports_resolver_id_fkey  FOREIGN KEY (resolver_id)   REFERENCES employees (id)   ON DELETE SET NULL
);

CREATE TRIGGER containers_reports_update_modified_at
    BEFORE UPDATE ON containers_reports
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();

CREATE INDEX containers_reports_issuer_id_idx ON containers_reports (issuer_id);

-- User container bookmarks.
CREATE TABLE users_container_bookmarks (
    user_id         uuid        NOT NULL,
    container_id    uuid        NOT NULL,
    created_at      timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_container_bookmarks_pkey               PRIMARY KEY (user_id, container_id),
    CONSTRAINT users_container_bookmarks_user_id_fkey       FOREIGN KEY (user_id)               REFERENCES users (id)       ON DELETE CASCADE,
    CONSTRAINT users_container_bookmarks_container_id_fkey  FOREIGN KEY (container_id)          REFERENCES containers (id)
);

-- Trucks.
CREATE TABLE trucks (
    id              uuid                    NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    make            varchar(50)             NOT NULL,
    model           varchar(50)             NOT NULL,
    license_plate   varchar(30)             NOT NULL,
    person_capacity integer                 NOT NULL,
    geom            geometry('POINT', 4326) NOT NULL,
    created_at      timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT trucks_pkey                              PRIMARY KEY (id),
    CONSTRAINT trucks_person_capacity_positive_check    CHECK (person_capacity > 0)
);

CREATE TRIGGER trucks_update_modified_at
    BEFORE UPDATE ON trucks
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();

-- Warehouses.
CREATE TABLE warehouses (
    id              uuid                    NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    truck_capacity  integer                 NOT NULL,
    geom            geometry('POINT', 4326) NOT NULL,
    created_at      timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT warehouses_pkey                          PRIMARY KEY (id),
    CONSTRAINT warehouses_truck_capacity_positive_check CHECK (truck_capacity >= 0)
);

CREATE TRIGGER warehouses_update_modified_at
    BEFORE UPDATE ON warehouses
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();

-- Warehouse trucks.
CREATE TABLE warehouses_trucks (
    warehouse_id    uuid        NOT NULL,
    truck_id        uuid        NOT NULL,
    created_at      timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT warehouses_trucks_pkey               PRIMARY KEY (warehouse_id, truck_id),
    CONSTRAINT warehouses_trucks_warehouse_id_fkey  FOREIGN KEY (warehouse_id)              REFERENCES warehouses (id)  ON DELETE CASCADE,
    CONSTRAINT warehouses_trucks_truck_id_fkey      FOREIGN KEY (truck_id)                  REFERENCES trucks (id),
    CONSTRAINT warehouses_trucks_truck_id_key       UNIQUE (truck_id)
);

-- Routes.
CREATE TABLE routes (
    id                      uuid        NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    name                    varchar(50) NOT NULL,
    truck_id                uuid        NOT NULL,
    departure_warehouse_id  uuid        NOT NULL,
    arrival_warehouse_id    uuid        NOT NULL,
    created_at              timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at             timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT routes_pkey                          PRIMARY KEY (id),
    CONSTRAINT routes_truck_id_fkey                 FOREIGN KEY (truck_id)                  REFERENCES trucks (id),
    CONSTRAINT routes_departure_warehouse_id_fkey   FOREIGN KEY (departure_warehouse_id)    REFERENCES warehouses (id),
    CONSTRAINT routes_arrival_warehouse_id_fkey     FOREIGN KEY (arrival_warehouse_id)      REFERENCES warehouses (id)
);

CREATE TRIGGER routes_update_modified_at
    BEFORE UPDATE ON routes
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();

-- Route containers.
CREATE TABLE routes_containers (
    route_id        uuid        NOT NULL,
    container_id    uuid        NOT NULL,
    created_at      timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT routes_containers_pkey               PRIMARY KEY (route_id, container_id),
    CONSTRAINT routes_containers_route_id_fkey      FOREIGN KEY (route_id)                  REFERENCES routes (id)      ON DELETE CASCADE,
    CONSTRAINT routes_containers_container_id_fkey  FOREIGN KEY (container_id)              REFERENCES containers (id)
);

-- Route employees.
CREATE TYPE routes_employees_employee_role AS ENUM ('driver', 'collector');

CREATE TABLE routes_employees (
    route_id        uuid                            NOT NULL,
    employee_id     uuid                            NOT NULL,
    employee_role   routes_employees_employee_role  NOT NULL,
    created_at      timestamp                       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT routes_employees_pkey                PRIMARY KEY (route_id, employee_id),
    CONSTRAINT routes_employees_route_id_fkey       FOREIGN KEY (route_id)              REFERENCES routes (id)      ON DELETE CASCADE,
    CONSTRAINT routes_employees_employee_id_fkey    FOREIGN KEY (employee_id)           REFERENCES employees (id)
);
