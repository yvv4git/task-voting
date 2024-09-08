-- Добавляем 10 пользователей в таблицу users с полем email
INSERT INTO users (name, email, created_at)
VALUES ('user1', 'user1@example.com', NOW()),
       ('user2', 'user2@example.com', NOW()),
       ('user3', 'user3@example.com', NOW()),
       ('user4', 'user4@example.com', NOW()),
       ('user5', 'user5@example.com', NOW()),
       ('user6', 'user6@example.com', NOW()),
       ('user7', 'user7@example.com', NOW()),
       ('user8', 'user8@example.com', NOW()),
       ('user9', 'user9@example.com', NOW()),
       ('user10', 'user10@example.com', NOW());

-- Добавляем учетные данные в таблицу auth
-- Предполагается, что пароли хешированы с использованием bcrypt
INSERT INTO auth (username, password, created_at)
VALUES ('user1', '$2a$10$ZlmY6kILxou32no9SoZLjeo94LmayVQdgm3SxGK7ED0DnlW/z9v7W', NOW()), -- password: password1
       ('user2', '$2a$10$zwXRgKtkS/wczbB3mmEOROzuEI1zw03vsbf4IyHmw9dfCB2F0AhRm', NOW()), -- password: password2
       ('user3', '$2a$10$1hlegiiX7ybCdn1EmvBSSuX5c7V6xVz4uwEukfxyLwezQC9ZZi5xS', NOW()), -- password: password3
       ('user4', '$2a$10$ZBK1mD.98uh3/EaoU3DUyOBacRJp9tuvMzRbcH47AHPNJ/D4Xlgre', NOW()), -- password: password4
       ('user5', '$2a$10$QoA6Q2eqPT1MVH6ocTTv4OBAUDf8ghhlgXXvi.keeUUCsgLwhxziC', NOW()), -- password: password5
       ('user6', '$2a$10$W83TtW4zVKApkA2tmzb.EuLX9ppplNhc8IfsMLV7ZJc7ix6e58A0m', NOW()), -- password: password6
       ('user7', '$2a$10$IfwQ3OsXpO4FzZI7zhRNCuMWc9gTziM9kBPiTGncLXu2pZcnZ3qR6', NOW()), -- password: password7
       ('user8', '$2a$10$CAl6J6YR7FaSHBqbdkUoAukHmhr4FTRUpg5Wtfa9.Z2jmAE6KlKKK', NOW()), -- password: password8
       ('user9', '$2a$10$1IXH2QfRrqxehGm0JAXomem0lK4qhSYhhxH0FF4lsbPFpcR1NZ0Qe', NOW()), -- password: password9
       ('user10', '$2a$10$/boH7yi4vXCoMIh0AuiBcOEDO/3ASTnBTDbmcwk/w6W1qbocU2GLu', NOW()); -- password: password10