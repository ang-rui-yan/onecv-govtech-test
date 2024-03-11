package cli

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	rootCmd     = &cobra.Command{Use: "myapp"}
	ctx         = context.Background()
	dbPool = 	&pgxpool.Pool{}
)

var createTableCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise the database by creating tables and inserting data",
	Run: func(cmd *cobra.Command, args []string) {
		createTables(dbPool)
		insertDummaryData(dbPool)
		log.Printf("Database has been initialise")
		os.Exit(1)
	},
}

var resetDataCmd = &cobra.Command{
	Use:   "reset",
	Short: "Apply migration on the database",
	Run: func(cmd *cobra.Command, args []string) {
		resetData(dbPool)
		log.Printf("All data has been reset")
		os.Exit(1)
	},
}

func Execute(postgresAddress string) {

	db, err := sql.Open("postgres", postgresAddress)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    var version string
    if err := db.QueryRow("select version()").Scan(&version); err != nil {
        panic(err)
    }

    fmt.Printf("version=%s\n", version)

	ctx := context.Background()

	// Create a connection pool to the studentadmindb database
	dbPool, err = pgxpool.New(ctx, postgresAddress)
	if err != nil {
		log.Fatalf("Unable to connect to studentadmindb: %v\n", err)
	}
	defer dbPool.Close()
	
	// Add the subcommands to the root command
	rootCmd.AddCommand(createTableCmd)
	rootCmd.AddCommand(resetDataCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}

func createTables(dbPool *pgxpool.Pool) {
	query := `
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL
);
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    suspended BOOLEAN NOT NULL DEFAULT false
);
CREATE TABLE teacher_students (
    teacher_id INT NOT NULL,
    student_id INT NOT NULL,
    PRIMARY KEY (teacher_id, student_id),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);
`
	_, err := dbPool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Could not create tables: %v\n", err)
	}
}

func insertDummaryData(dbPool *pgxpool.Pool) {
	query := `
INSERT INTO teachers (email) VALUES
('teacherken@gmail.com'),
('teacherjoe@gmail.com')
ON CONFLICT (email) DO NOTHING;

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
	`
	_, err := dbPool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Could not insert sample data: %v\n", err)
	}
}

func resetData(dbPool *pgxpool.Pool) {
	query := `
DROP TABLE IF EXISTS teacher_students;
DROP TABLE IF EXISTS teachers;
DROP TABLE IF EXISTS students;
	`
	_, err := dbPool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Could not drop table: %v\n", err)
	}
}