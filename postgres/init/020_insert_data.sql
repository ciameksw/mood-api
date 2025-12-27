\connect mood_api_db

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

INSERT INTO public.advice_type (name, description) VALUES
('Motivation', 'Inspirational content to boost energy and drive'),
('Self-Care', 'Tips for taking care of your physical and mental well-being'),
('Mindfulness', 'Practices for being present and aware'),
('Productivity', 'Strategies to get things done efficiently'),
('Relaxation', 'Techniques to unwind and reduce tension'),
('Social Connection', 'Ideas for connecting with others'),
('Exercise', 'Physical activity suggestions'),
('Gratitude', 'Practices for appreciation and thankfulness'),
('Problem Solving', 'Approaches to tackle challenges'),
('Breathing Exercises', 'Techniques to calm the mind and body');

INSERT INTO public.mood_advice_type_mapping (mood_type_id, advice_type_id, priority) VALUES
-- Happy (id=1)
(1, 4, 1),  -- Happy -> Productivity
(1, 7, 2),  -- Happy -> Exercise
(1, 6, 3),  -- Happy -> Social Connection
(1, 8, 4),  -- Happy -> Gratitude

-- Sad (id=2)
(2, 2, 1),  -- Sad -> Self-Care
(2, 1, 2),  -- Sad -> Motivation
(2, 6, 3),  -- Sad -> Social Connection
(2, 8, 4),  -- Sad -> Gratitude

-- Anxious (id=3)
(3, 10, 1), -- Anxious -> Breathing Exercises
(3, 3, 2),  -- Anxious -> Mindfulness
(3, 5, 3),  -- Anxious -> Relaxation
(3, 9, 4),  -- Anxious -> Problem Solving

-- Calm (id=4)
(4, 3, 1),  -- Calm -> Mindfulness
(4, 8, 2),  -- Calm -> Gratitude
(4, 4, 3),  -- Calm -> Productivity
(4, 5, 4),  -- Calm -> Relaxation

-- Energetic (id=5)
(5, 7, 1),  -- Energetic -> Exercise
(5, 4, 2),  -- Energetic -> Productivity
(5, 6, 3),  -- Energetic -> Social Connection
(5, 1, 4),  -- Energetic -> Motivation

-- Tired (id=6)
(6, 5, 1),  -- Tired -> Relaxation
(6, 2, 2),  -- Tired -> Self-Care
(6, 3, 3),  -- Tired -> Mindfulness
(6, 10, 4), -- Tired -> Breathing Exercises

-- Angry (id=7)
(7, 10, 1), -- Angry -> Breathing Exercises
(7, 7, 2),  -- Angry -> Exercise
(7, 5, 3),  -- Angry -> Relaxation
(7, 9, 4),  -- Angry -> Problem Solving

-- Grateful (id=8)
(8, 8, 1),  -- Grateful -> Gratitude
(8, 6, 2),  -- Grateful -> Social Connection
(8, 3, 3),  -- Grateful -> Mindfulness
(8, 1, 4),  -- Grateful -> Motivation

-- Stressed (id=9)
(9, 10, 1), -- Stressed -> Breathing Exercises
(9, 5, 2),  -- Stressed -> Relaxation
(9, 9, 3),  -- Stressed -> Problem Solving
(9, 2, 4),  -- Stressed -> Self-Care

-- Neutral (id=10)
(10, 3, 1), -- Neutral -> Mindfulness
(10, 4, 2), -- Neutral -> Productivity
(10, 8, 3), -- Neutral -> Gratitude
(10, 7, 4); -- Neutral -> Exercise