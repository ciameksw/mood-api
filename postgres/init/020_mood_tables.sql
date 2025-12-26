\connect mood_db

CREATE TABLE IF NOT EXISTS public.mood_type (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) UNIQUE NOT NULL,
	description TEXT
);

INSERT INTO public.mood_type (name, description) VALUES
('Happy', 'Feeling joyful, content, and positive about the day'),
('Sad', 'Feeling down, melancholic, or experiencing a sense of loss'),
('Anxious', 'Feeling worried, nervous, or uneasy about present or future events'),
('Calm', 'Feeling peaceful, relaxed, and in a state of tranquility'),
('Energetic', 'Feeling full of energy, motivated, and ready to take on challenges'),
('Tired', 'Feeling exhausted, drained, or lacking physical or mental energy'),
('Angry', 'Feeling frustrated, irritated, or experiencing strong displeasure'),
('Grateful', 'Feeling thankful and appreciative of people or circumstances'),
('Stressed', 'Feeling overwhelmed by pressures, demands, or responsibilities'),
('Neutral', 'Feeling balanced with no strong emotions, just going through the day');

CREATE TABLE IF NOT EXISTS public.mood (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	mood_type_id INT REFERENCES public.mood_type(id),
	note TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

