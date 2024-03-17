-- TODO: Replace this first migration with the real one.
CREATE TABLE employee (
    id              UUID        NOT NULL    DEFAULT gen_random_uuid(),
    name            VARCHAR(70) NOT NULL,
    date_of_birth   DATE        NOT NULL,
    CONSTRAINT employee_pk PRIMARY KEY (id)
);
