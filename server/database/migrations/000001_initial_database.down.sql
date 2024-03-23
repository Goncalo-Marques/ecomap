-- Route employees.
DROP TABLE public.routes_employees;
DROP TYPE public.routes_employees_employee_role;

-- Route containers.
DROP TABLE public.routes_containers;

-- Routes.
DROP TABLE public.routes;

-- Warehouse trucks.
DROP TABLE public.warehouses_trucks;

-- Warehouses.
DROP TABLE public.warehouses;

-- Trucks.
DROP TABLE public.trucks;

-- User container bookmarks.
DROP TABLE public.users_container_bookmarks;

-- Container reports.
DROP TABLE public.containers_reports;
DROP TYPE public.containers_reports_issue_type;

-- Containers.
DROP TABLE public.containers;
DROP TYPE public.containers_container_type;

-- Employees.
DROP TABLE public.employees;
DROP TYPE public.employees_employee_type;

-- Users.
DROP TABLE public.users;

-- Functions.
DROP FUNCTION public.enforce_lower_case_username();
DROP FUNCTION public.update_modified_time();

-- Extensions.
DROP EXTENSION pgrouting CASCADE;
DROP EXTENSION postgis CASCADE;
