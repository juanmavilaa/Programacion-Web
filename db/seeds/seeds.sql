-- Limpieza opcional de tablas
TRUNCATE TABLE workouts, clients RESTART IDENTITY CASCADE;

-- Clientes iniciales
INSERT INTO clients (name, email) VALUES
('Carlos Gómez', 'carlos.gomez@example.com'),
('Lucía Fernández', 'lucia.fernandez@example.com'),
('Martín Rodríguez', 'martin.rodriguez@example.com');

-- Entrenamientos iniciales
INSERT INTO workouts (client_id, exercise, sets, repetitions, weight, created_at) VALUES
(1, 'Sentadilla con barra', 4, 6, 120.00, NOW() - INTERVAL '3 days'),
(1, 'Press de banca plano', 4, 8, 85.00, NOW() - INTERVAL '2 days'),
(2, 'Remo con barra', 4, 10, 45.00, NOW() - INTERVAL '2 days'),
(3, 'Prensa de piernas', 4, 12, 180.00, NOW() - INTERVAL '1 day');

