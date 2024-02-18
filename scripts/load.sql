
-- \dt
-- drop table teacher_students, teachers, students;

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

INSERT INTO teacher_students (teacher_id, student_id) VALUES
(1, 1),
(1, 2),
(1, 3),
(1, 4),
(1, 5),
(2, 4),
(2, 5),
(2, 6),
(2, 7);