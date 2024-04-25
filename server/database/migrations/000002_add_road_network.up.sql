-- Data converted from https://download.geofabrik.de/europe/portugal.html using https://osm2po.de/
-- Due to the size of the generated SQL files, it is expected that the scripts in `server/database/scripts` are used to populate the road network tables.

-- Extension.
CREATE EXTENSION pg_trgm;

-- Road network.
CREATE TABLE road_network(id integer, osm_id bigint, osm_name character varying, osm_meta character varying, osm_source_id bigint, osm_target_id bigint, clazz integer, flags integer, source integer, target integer, km double precision, kmh integer, cost double precision, reverse_cost double precision, x1 double precision, y1 double precision, x2 double precision, y2 double precision);
SELECT AddGeometryColumn('road_network', 'geom_way', 4326, 'LINESTRING', 2);

ALTER TABLE road_network ADD CONSTRAINT road_network_pkey PRIMARY KEY(id);
CREATE INDEX road_network_source_idx ON road_network(source);
CREATE INDEX road_network_target_idx ON road_network(target);
CREATE INDEX road_network_osm_source_id_idx ON road_network(osm_source_id);
CREATE INDEX road_network_geom_way_idx ON road_network USING gist (geom_way);
CREATE INDEX road_network_osm_target_id_idx ON road_network(osm_target_id);
CREATE INDEX road_network_osm_name_idx ON road_network USING gin (osm_name gin_trgm_ops);

-- Road network vertex.
CREATE TABLE road_network_vertex(id integer, clazz integer, osm_id bigint, osm_name character varying, ref_count integer, restrictions character varying);
SELECT AddGeometryColumn('road_network_vertex', 'geom_vertex', 4326, 'POINT', 2);

ALTER TABLE road_network_vertex ADD CONSTRAINT road_network_vertex_pkey PRIMARY KEY(id);
CREATE INDEX road_network_vertex_osm_id_idx ON road_network_vertex(osm_id);
CREATE INDEX road_network_vertex_geom_vertex_idx ON road_network_vertex USING gist (geom_vertex);

-- Road network foreign keys.
ALTER TABLE employees ADD road_id integer;
ALTER TABLE employees ADD CONSTRAINT employees_road_id_fkey FOREIGN KEY (road_id) REFERENCES road_network (id);

ALTER TABLE containers ADD road_id integer;
ALTER TABLE containers ADD CONSTRAINT containers_road_id_fkey FOREIGN KEY (road_id) REFERENCES road_network (id);

ALTER TABLE warehouses ADD road_id integer;
ALTER TABLE warehouses ADD CONSTRAINT warehouses_road_id_fkey FOREIGN KEY (road_id) REFERENCES road_network (id);
