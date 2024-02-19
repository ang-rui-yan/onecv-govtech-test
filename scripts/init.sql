CREATE DATABASE studentadmindb;

\c studentadmindb;

-- create tables
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

BEGIN;

-- INSERT TEACHERS
INSERT INTO teachers (email) VALUES
('teacherken@gmail.com'),
('teacherjoe@gmail.com')
ON CONFLICT (email) DO NOTHING;

-- INSERT STUDENTS
-- DEFAULT FOR suspended IS FALSE
INSERT INTO students (email) VALUES
('studentjon@gmail.com'),
('studenthon@gmail.com'),
('commonstudent1@gmail.com'),
('commonstudent2@gmail.com'),
('student_only_under_teacher_ken@gmail.com'),
('studentmary@gmail.com'),
('studentbob@gmail.com'),
('studentagnes@gmail.com'),
('studentmiche@gmail.com')
ON CONFLICT (email) DO NOTHING;

-- INSERT INTO teacher_students (teacher_id, student_id) VALUES
-- (1, 1),
-- (1, 2),
-- (1, 3),
-- (1, 4),
-- (1, 5),
-- (2, 4),
-- (2, 5),
-- (2, 6),
-- (2, 7);

COMMIT;


-- \dt
-- drop table teacher_students, teachers, students;