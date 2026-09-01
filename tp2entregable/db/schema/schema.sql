CREATE TABLE workouts (
    id SERIAL PRIMARY KEY,
    exercise VARCHAR(255) NOT NULL,
    sets INTEGER NOT NULL,
    repetitions INTEGER NOT NULL,
    weight NUMERIC(6,2) NOT NULL,
    workout_date DATE NOT NULL
);

