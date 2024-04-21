ALTER TABLE employees 
ADD way_osm_id bigint
CONSTRAINT employees_way_osm_id_fkey FOREIGN KEY (way_osm_id) REFERENCES road_network (osm_id);

ALTER TABLE containers 
ADD way_osm_id bigint
CONSTRAINT containers_way_osm_id_fkey FOREIGN KEY (way_osm_id) REFERENCES road_network (osm_id);

ALTER TABLE warehouses 
ADD way_osm_id bigint
CONSTRAINT warehouses_way_osm_id_fkey FOREIGN KEY (way_osm_id) REFERENCES road_network (osm_id);
