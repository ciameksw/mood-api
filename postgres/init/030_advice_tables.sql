\connect advice_db

CREATE TABLE IF NOT EXISTS public.advice_type (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) UNIQUE NOT NULL,
	description TEXT
);

CREATE TABLE IF NOT EXISTS public.advice (
	id SERIAL PRIMARY KEY,
	advice_type_id INT REFERENCES public.advice_type(id),
	title VARCHAR(200),
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

