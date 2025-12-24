CREATE TABLE
    users (
        id bigserial NOT NULL,
        name VARCHAR NOT NULL,
        CONSTRAINT user_pkey PRIMARY KEY (id)
    );