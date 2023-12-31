/*
    Схема БД для информационной системы
    отслеживания выполнения задач.

    docker build -t my-postgres .
    docker run -d --name my-postgres-container -p 5432:5432 my-postgres
*/
--CREATE DATABASE testdb;

DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;

-- пользователи системы
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- метки задач
CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- задачи
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT NOT NULL DEFAULT extract(epoch from now()), -- время создания задачи
    closed BIGINT DEFAULT 0, -- время выполнения задачи
    author_id INTEGER REFERENCES users(id) DEFAULT 1, -- автор задачи
    assigned_id INTEGER REFERENCES users(id) DEFAULT 1, -- ответственный
    title TEXT, -- название задачи
    content TEXT -- задачи
);

-- связь многие - ко- многим между задачами и метками
CREATE TABLE tasks_labels (
    task_id INTEGER REFERENCES tasks(id),
    label_id INTEGER REFERENCES labels(id)
);


-- наполнение БД начальными данными
INSERT INTO users (name) VALUES ('ollya');
INSERT INTO users (name) VALUES ('sergius');
INSERT INTO users (name) VALUES ('geralt');