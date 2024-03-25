-- Extensions.
CREATE EXTENSION postgis;
CREATE EXTENSION pgrouting;

-- Functions.
CREATE FUNCTION public.enforce_lower_case_username() 
RETURNS TRIGGER AS $$
BEGIN
    new.username = LOWER(new.username);
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION public.update_modified_time() 
RETURNS TRIGGER AS $$
BEGIN
    new.modified_time = CURRENT_TIMESTAMP;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

-- Users.
CREATE TABLE public.users (
    id              uuid        NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    username        varchar(50) NOT NULL,
    password        varchar(60) NOT NULL,
    first_name      varchar(50) NOT NULL,
    last_name       varchar(50) NOT NULL,
    created_time    timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_time   timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey           PRIMARY KEY (id),
    CONSTRAINT users_username_key   UNIQUE (username)
);

CREATE TRIGGER users_enforce_lower_case_username
    BEFORE INSERT OR UPDATE ON public.users
    FOR EACH ROW
    EXECUTE PROCEDURE public.enforce_lower_case_username();

CREATE TRIGGER users_update_modified_time
    BEFORE UPDATE ON public.users
    FOR EACH ROW
    EXECUTE PROCEDURE public.update_modified_time();

-- Employees.
CREATE TYPE public.employees_employee_type AS ENUM ('waste_operator', 'manager');

CREATE TABLE public.employees (
    id              uuid                            NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    username        varchar(50)                     NOT NULL,
    password        varchar(60)                     NOT NULL,
    first_name      varchar(50)                     NOT NULL,
    last_name       varchar(50)                     NOT NULL,
    type            public.employees_employee_type  NOT NULL,
    date_of_birth   date                            NOT NULL,
    phone_number    varchar(20)                     NOT NULL,
    geom            geometry                        NOT NULL,
    schedule_start  time                            NOT NULL,
    schedule_end    time                            NOT NULL,
    created_time    timestamp                       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_time   timestamp                       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT employees_pkey           PRIMARY KEY (id),
    CONSTRAINT employees_username_key   UNIQUE (username)
);

CREATE TRIGGER employees_enforce_lower_case_username
    BEFORE INSERT OR UPDATE ON public.employees
    FOR EACH ROW
    EXECUTE PROCEDURE public.enforce_lower_case_username();

CREATE TRIGGER employees_update_modified_time
    BEFORE UPDATE ON public.employees
    FOR EACH ROW
    EXECUTE PROCEDURE public.update_modified_time();

-- Containers.
CREATE TYPE public.containers_container_type AS ENUM ('general', 'paper', 'plastic', 'metal', 'glass', 'organic', 'hazardous');

CREATE TABLE public.containers (
    id              uuid                                NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    type            public.containers_container_type    NOT NULL,
    geom            geometry                            NOT NULL,
    created_time    timestamp                           NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_time   timestamp                           NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT containers_pkey  PRIMARY KEY (id)
);

CREATE TRIGGER containers_update_modified_time
    BEFORE UPDATE ON public.containers
    FOR EACH ROW
    EXECUTE PROCEDURE public.update_modified_time();

-- Container reports.
CREATE TYPE public.containers_reports_issue_type AS ENUM ('full', 'vandalized', 'misplaced', 'non-existent', 'other');

CREATE TABLE public.containers_reports (
    id              uuid                                    NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    container_id    uuid                                    NOT NULL,
    issue_type      public.containers_reports_issue_type    NOT NULL,
    description     varchar(500),
    attachment      bytea,
    issuer_id       uuid                                    NOT NULL,
    resolver_id     uuid,
    resolved        boolean                                 NOT NULL    DEFAULT FALSE,
    created_time    timestamp                               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_time   timestamp                               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT containers_reports_pkey              PRIMARY KEY (id),
    CONSTRAINT containers_reports_container_id_fkey FOREIGN KEY (container_id)  REFERENCES public.containers (id),
    CONSTRAINT containers_reports_issuer_id_fkey    FOREIGN KEY (issuer_id)     REFERENCES public.users (id),
    CONSTRAINT containers_reports_resolver_id_fkey  FOREIGN KEY (resolver_id)   REFERENCES public.employees (id)
);

CREATE TRIGGER containers_reports_update_modified_time
    BEFORE UPDATE ON public.containers_reports
    FOR EACH ROW
    EXECUTE PROCEDURE public.update_modified_time();

CREATE INDEX containers_reports_issuer_id_idx ON public.containers_reports (issuer_id);

-- User container bookmarks.
CREATE TABLE public.users_container_bookmarks (
    user_id         uuid        NOT NULL,
    container_id    uuid        NOT NULL,
    created_time    timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_container_bookmarks_pkey               PRIMARY KEY (user_id, container_id),
    CONSTRAINT users_container_bookmarks_user_id_fkey       FOREIGN KEY (user_id)               REFERENCES public.users (id),
    CONSTRAINT users_container_bookmarks_container_id_fkey  FOREIGN KEY (container_id)          REFERENCES public.containers (id)
);

-- Trucks.
CREATE TABLE public.trucks (
    id              uuid        NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    make            varchar(50) NOT NULL,
    model           varchar(50) NOT NULL,
    license_plate   varchar(30) NOT NULL,
    person_capacity integer     NOT NULL,
    created_time    timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_time   timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT trucks_pkey                              PRIMARY KEY (id),
    CONSTRAINT trucks_person_capacity_positive_check    CHECK (person_capacity > 0)
);

CREATE TRIGGER trucks_update_modified_time
    BEFORE UPDATE ON public.trucks
    FOR EACH ROW
    EXECUTE PROCEDURE public.update_modified_time();

-- Warehouses.
CREATE TABLE public.warehouses (
    id              uuid        NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    geom            geometry    NOT NULL,
    truck_capacity  integer     NOT NULL,
    created_time    timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_time   timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT warehouses_pkey                          PRIMARY KEY (id),
    CONSTRAINT warehouses_truck_capacity_positive_check CHECK (truck_capacity > 0)
);

CREATE TRIGGER warehouses_update_modified_time
    BEFORE UPDATE ON public.warehouses
    FOR EACH ROW
    EXECUTE PROCEDURE public.update_modified_time();

-- Warehouse trucks.
CREATE TABLE public.warehouses_trucks (
    warehouse_id    uuid        NOT NULL,
    truck_id        uuid        NOT NULL,
    created_time    timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT warehouses_trucks_pkey               PRIMARY KEY (warehouse_id, truck_id),
    CONSTRAINT warehouses_trucks_warehouse_id_fkey  FOREIGN KEY (warehouse_id)              REFERENCES public.warehouses (id),
    CONSTRAINT warehouses_trucks_truck_id_fkey      FOREIGN KEY (truck_id)                  REFERENCES public.trucks (id),
    CONSTRAINT warehouses_trucks_truck_id_key       UNIQUE (truck_id)
);

-- Routes.
CREATE TABLE public.routes (
    id                      uuid        NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    truck_id                uuid        NOT NULL,
    departure_warehouse_id  uuid        NOT NULL,
    arrival_warehouse_id    uuid        NOT NULL,
    created_time            timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT routes_pkey                          PRIMARY KEY (id),
    CONSTRAINT routes_truck_id_fkey                 FOREIGN KEY (truck_id)                  REFERENCES public.trucks (id),
    CONSTRAINT routes_departure_warehouse_id_fkey   FOREIGN KEY (departure_warehouse_id)    REFERENCES public.warehouses (id),
    CONSTRAINT routes_arrival_warehouse_id_fkey     FOREIGN KEY (arrival_warehouse_id)      REFERENCES public.warehouses (id)
);

-- Route containers.
CREATE TABLE public.routes_containers (
    route_id        uuid        NOT NULL,
    container_id    uuid        NOT NULL,
    emptied         boolean,
    washed          boolean,
    responsible_id  uuid,
    created_time    timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_time   timestamp   NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT routes_containers_pkey                   PRIMARY KEY (route_id, container_id),
    CONSTRAINT routes_containers_route_id_fkey          FOREIGN KEY (route_id)                  REFERENCES public.routes (id),
    CONSTRAINT routes_containers_container_id_fkey      FOREIGN KEY (container_id)              REFERENCES public.containers (id),
    CONSTRAINT routes_containers_responsible_id_fkey    FOREIGN KEY (responsible_id)            REFERENCES public.employees (id)
);

CREATE TRIGGER routes_containers_update_modified_time
    BEFORE UPDATE ON public.routes_containers
    FOR EACH ROW
    EXECUTE PROCEDURE public.update_modified_time();

CREATE INDEX routes_containers_created_time_idx ON public.routes_containers (created_time);

-- Route employees.
CREATE TYPE public.routes_employees_employee_role AS ENUM ('driver', 'collector');

CREATE TABLE public.routes_employees (
    route_id        uuid                                    NOT NULL,
    employee_id     uuid                                    NOT NULL,
    employee_role   public.routes_employees_employee_role   NOT NULL,
    created_time    timestamp                               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT routes_employees_pkey                PRIMARY KEY (route_id, employee_id),
    CONSTRAINT routes_employees_route_id_fkey       FOREIGN KEY (route_id)              REFERENCES public.routes (id),
    CONSTRAINT routes_employees_employee_id_fkey    FOREIGN KEY (employee_id)           REFERENCES public.employees (id)
);
