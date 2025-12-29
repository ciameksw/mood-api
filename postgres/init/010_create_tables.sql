\connect mood_api_db

CREATE TABLE IF NOT EXISTS public.users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.mood_type (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) UNIQUE NOT NULL,
	description TEXT
);

CREATE TABLE IF NOT EXISTS public.mood (
	id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES public.users(id),
    mood_date DATE NOT NULL DEFAULT CURRENT_DATE,
	mood_type_id INT REFERENCES public.mood_type(id),
	note TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (user_id, mood_date)
);

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

CREATE TABLE IF NOT EXISTS public.mood_advice_type_mapping (
    id SERIAL PRIMARY KEY,
    mood_type_id INT NOT NULL REFERENCES public.mood_type(id),
    advice_type_id INT NOT NULL REFERENCES public.advice_type(id),
    priority INT DEFAULT 1, -- Lower number = higher priority
    UNIQUE(mood_type_id, advice_type_id)
);

CREATE TABLE IF NOT EXISTS public.user_advice_periods (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES public.users(id),
	advice_id INT NOT NULL REFERENCES public.advice(id),
	period_from DATE NOT NULL,
	period_to DATE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (user_id, advice_id, period_from, period_to)
);