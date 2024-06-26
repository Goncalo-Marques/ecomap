-- Landfills.
CREATE TABLE landfills (
    id              uuid                    NOT NULL    DEFAULT GEN_RANDOM_UUID(),
    geom            geometry('POINT', 4326) NOT NULL,
    road_id         integer,
    municipality_id integer,
    created_at      timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    modified_at     timestamp               NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT landfills_pkey                   PRIMARY KEY (id),
    CONSTRAINT landfills_road_id_fkey           FOREIGN KEY (road_id)           REFERENCES road_network (id),
    CONSTRAINT landfills_municipality_id_fkey   FOREIGN KEY (municipality_id)   REFERENCES municipalities (id)
);

CREATE TRIGGER landfills_update_modified_at
    BEFORE UPDATE ON landfills
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();
