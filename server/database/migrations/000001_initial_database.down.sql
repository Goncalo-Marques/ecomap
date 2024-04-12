-- Route employees.
DROP TABLE routes_employees;
DROP TYPE routes_employees_employee_role;

-- Route containers.
DROP TABLE routes_containers;

-- Routes.
DROP TABLE routes;

-- Warehouse trucks.
DROP TABLE warehouses_trucks;

-- Warehouses.
DROP TABLE warehouses;

-- Trucks.
DROP TABLE trucks;

-- User container bookmarks.
DROP TABLE users_container_bookmarks;

-- Container reports.
DROP TABLE containers_reports;
DROP TYPE containers_reports_issue_type;

-- Containers.
DROP TABLE containers;
DROP TYPE containers_category;

-- Employees.
DROP TABLE employees;
DROP TYPE employees_role;

-- Users.
DROP TABLE users;

-- Functions.
DROP FUNCTION enforce_lower_case_username();
DROP FUNCTION update_modified_at();

-- Extensions.
DROP EXTENSION pgrouting CASCADE;
DROP EXTENSION postgis CASCADE;
