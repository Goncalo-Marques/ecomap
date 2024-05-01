-- Municipalities foreign keys.
ALTER TABLE employees DROP COLUMN municipality_id;
ALTER TABLE containers DROP COLUMN municipality_id;
ALTER TABLE trucks DROP COLUMN municipality_id;
ALTER TABLE warehouses DROP COLUMN municipality_id;

-- Municipalities.
DROP TABLE municipalities;
