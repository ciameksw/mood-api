\connect mood_db

CREATE TABLE IF NOT EXISTS public.mood_type (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) UNIQUE NOT NULL,
	description TEXT
);

CREATE TABLE IF NOT EXISTS public.mood (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	mood_type_id INT REFERENCES public.mood_type(id),
	note TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

