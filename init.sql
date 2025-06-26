CREATE TABLE IF NOT EXISTS genomes (
    id SERIAL PRIMARY KEY,
    label TEXT NOT NULL
);

INSERT INTO genomes (label)
VALUES
    ('Homo sapiens genome A'),
    ('Homo sapiens genome B'),
    ('Sample genome test C');
