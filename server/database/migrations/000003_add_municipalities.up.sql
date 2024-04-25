-- Data converted from https://www.dgterritorio.gov.pt/cartografia/cartografia-tematica/caop
-- Due to the size of the generated SQL file, it is expected that the script in `server/database/scripts` is used to populate the municipalities table.

-- Municipalities.
CREATE TABLE municipalities (
    id              serial                      NOT NULL,
    geom            geometry('Polygon', 4326)   NOT NULL,
    fid             bigint                      NOT NULL,
    name            varchar(50)                 NOT NULL,
    district        varchar(50)                 NOT NULL,
    nutsiii         varchar(50)                 NOT NULL,
    nutsii          varchar(50)                 NOT NULL,
    nutsi           varchar(50)                 NOT NULL,
    area_ha         double precision            NOT NULL,
    perimeter_km    bigint                      NOT NULL,
    CONSTRAINT municipalities_pkey PRIMARY KEY (id)
);

CREATE INDEX municipalities_geom_idx ON municipalities USING gist (geom);

-- Municipalities foreign keys.
ALTER TABLE employees ADD municipality_id integer;
ALTER TABLE employees ADD CONSTRAINT employees_municipality_id_fkey FOREIGN KEY (municipality_id) REFERENCES municipalities (id);

ALTER TABLE containers ADD municipality_id integer;
ALTER TABLE containers ADD CONSTRAINT containers_municipality_id_fkey FOREIGN KEY (municipality_id) REFERENCES municipalities (id);

ALTER TABLE warehouses ADD municipality_id integer;
ALTER TABLE warehouses ADD CONSTRAINT warehouses_municipality_id_fkey FOREIGN KEY (municipality_id) REFERENCES municipalities (id);
