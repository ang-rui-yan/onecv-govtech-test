CREATE DATABASE studentadmindb;

USE studentadmindb;

CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    suspended BOOLEAN NOT NULL DEFAULT false
);

-- many to many junction table
CREATE TABLE teacher_students (
    teacher_id INT NOT NULL,
    student_id INT NOT NULL,
    PRIMARY KEY (teacher_id, student_id),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);
