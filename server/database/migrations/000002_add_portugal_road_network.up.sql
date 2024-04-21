-- Data converted from https://download.geofabrik.de/europe/portugal.html using https://osm2po.de/
-- Due to the size of the generated SQL files, it is expected that the scripts in `server/database/scripts` are used to populate the road network tables.

CREATE TABLE road_network(id integer, osm_id bigint, osm_name character varying, osm_meta character varying, osm_source_id bigint, osm_target_id bigint, clazz integer, flags integer, source integer, target integer, km double precision, kmh integer, cost double precision, reverse_cost double precision, x1 double precision, y1 double precision, x2 double precision, y2 double precision);
SELECT AddGeometryColumn('road_network', 'geom_way', 4326, 'LINESTRING', 2);

ALTER TABLE road_network ADD CONSTRAINT road_network_pkey PRIMARY KEY(id);
CREATE INDEX road_network_source_idx ON road_network(source);
CREATE INDEX road_network_target_idx ON road_network(target);
-- CREATE INDEX road_network_osm_source_id_idx ON road_network(osm_source_id);
-- CREATE INDEX road_network_osm_target_id_idx ON road_network(osm_target_id);
-- CREATE INDEX road_network_geom_way_idx ON road_network USING GIST (geom_way);

CREATE TABLE road_network_vertex(id integer, clazz integer, osm_id bigint, osm_name character varying, ref_count integer, restrictions character varying);
SELECT AddGeometryColumn('road_network_vertex', 'geom_vertex', 4326, 'POINT', 2);

ALTER TABLE road_network_vertex ADD CONSTRAINT road_network_vertex_pkey PRIMARY KEY(id);
CREATE INDEX road_network_vertex_osm_id_idx ON road_network_vertex(osm_id);
-- CREATE INDEX road_network_vertex_geom_vertex_idx ON road_network_vertex USING GIST (geom_vertex);
