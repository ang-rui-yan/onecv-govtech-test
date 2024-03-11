# OneCV GovTech Asessment

## Description

The objective of this assessment is to build a Golang backend application with PostgreSQL database. This application allows teachers to perform administrative functions for their students. Teachers and students are uniquely identified by their email addresses.

## Installation

1. Add in appropriate environment variables

   ```
   COPY .env.example .env
   ```

1. Update the postgresql connection URL in the .env

1. Install go mod

   ```
   go mod download
   ```

1. Initialise the database

   ```
   go run cmd/main.go init
   ```

1. Start the server

   ```
   go run .
   ```

1. To clear the database

   ```
   go run cmd/main.go reset
   ```

## API endpoints

1. POST /api/register

   **Headers**: Content-Type: application/json

   **Success response status**: HTTP 204

   **Example request body**

   ```json
   {
     "teacher": "teacherken@gmail.com",
     "students": ["studentjon@gmail.com", "studenthon@gmail.com"]
   }
   ```

1. GET /api/commonstudents

   **Headers**: Content-Type: application/json

   **Success response status**: HTTP 200

   **Querystring**:

   - teacher, could be multiple
     - GET /api/commonstudents?teacher=teacherken%40gmail.com
     - GET /api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com

   **Example response body**:

   ```json
   {
     "students": [
       "commonstudent1@gmail.com",
       "commonstudent2@gmail.com",
       "student_only_under_teacher_ken@gmail.com"
     ]
   }
   ```

1. POST /api/suspend

   **Headers**: Content-Type: application/json

   **Success response status**: HTTP 204

   **Example request body**:

   ```json
   {
     "student": "studentmary@gmail.com"
   }
   ```

1. POST /api/retrievefornotifications

   **Headers**: Content-Type: application/json

   **Success response status**: HTTP 200

   **Example request body**:

   ```json
   {
     "teacher": "teacherken@gmail.com",
     "notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
   }
   ```

   **Example response body**:

   ```json
   {
     "recipients": [
       "studentbob@gmail.com",
       "studentagnes@gmail.com",
       "studentmiche@gmail.com"
     ]
   }
   ```

## Running the Tests

```bash
go test ./...
```

## Author

Ang Rui Yan

## Acknowledgements

- The folder and file structure of this project was inspired by the organizational approach demonstrated in the `pgx-test-example` repository by [reshimahendra](https://github.com/reshimahendra/pgx-test-example/tree/master). I found the structuring and organization particularly effective for managing `pgx` connections and testing in Go applications.
