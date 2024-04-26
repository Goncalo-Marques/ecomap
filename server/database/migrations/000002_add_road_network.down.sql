-- Road network foreign keys.
ALTER TABLE employees DROP COLUMN road_id;
ALTER TABLE containers DROP COLUMN road_id;
ALTER TABLE warehouses DROP COLUMN road_id;

-- Road network vertex.
DROP TABLE road_network_vertex;

-- Road network.
DROP TABLE road_network;

-- Extension.
DROP EXTENSION pg_trgm;
