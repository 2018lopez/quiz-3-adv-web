--Filename: migrations/000001_todo.down.sql

CREATE TABLE
    IF NOT EXISTS todo(
        id bigserial PRIMARY KEY,
        created_at TIMESTAMP(0)
        with
            TIME Zone NOT null DEFAULT Now(),
            title text NOT NULL,
            description text NOT NULL,
            completed BOOLEAN
    );