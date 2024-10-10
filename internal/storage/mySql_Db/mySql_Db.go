package mySql_Db

import (
	"database/sql"
	"fmt"
	"github.com/SudipSarkar1193/students-API-Go/internal/config"
	"github.com/SudipSarkar1193/students-API-Go/internal/types"
	_ "github.com/go-sql-driver/mysql"
	// Import your types package where your Student struct is defined
)

// New function to create a new MySQL DB connection
func New(cfg *config.Config) (*sql.DB, error) {

	// Data Source Name (DSN): user:password@tcp(host:port)/
	// No database is specified initially
	dsn := cfg.Dsn

	// Connect to MySQL server without specifying a database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the MySQL server to verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to MySQL server!")

	// Create the student database if it doesn't exist
	query := "CREATE DATABASE IF NOT EXISTS student"
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	fmt.Println("Database 'student' created or already exists.")

	// Close the current connection and reconnect, specifying the 'student' database
	dsnWithDB := dsn + "student"
	dbWithDB, err := sql.Open("mysql", dsnWithDB)
	if err != nil {
		return nil, err
	}

	// Ping the MySQL server to verify the connection to the 'student' database
	if err := dbWithDB.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to the 'student' database!")

	// Create the students table if it doesn't exist
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS students (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL 
    );`

	_, err = dbWithDB.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	fmt.Println("Table 'students' created or already exists.")

	return dbWithDB, nil
}

func NewOnline(cfg *config.Config) (*sql.DB, error) {
	// Data Source Name (DSN) - connecting to the provided database directly
	dsn := cfg.Dsn + "?parseTime=true" // Your DSN should already have the database name

	// Connect to the MySQL server using the provided database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the MySQL server to ensure the connection is established
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to the database!")

	// Create the students table if it doesn't exist
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS students (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
    );`

	// Execute the query to create the table
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	fmt.Println("Table 'students' created or already exists.")

	// Return the connected database object
	return db, nil
}

// Function to insert a new student into the database
func InsertStudent(db *sql.DB, student *types.Student) error {
	query := "INSERT INTO students (name, email, password) VALUES (?, ?,?)"

	// Executing the query
	result, err := db.Exec(query, student.Name, student.Email, student.Password)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted student
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	student.Id = id

	fmt.Printf("New student inserted with ID: %d\n", id)
	return nil
}

// Function to insert a new student into the database
func GetAllStudents(db *sql.DB) ([]types.Student, error) {
	query := "SELECT * FROM students"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

type Error string

func (e Error) Error() string {
	return string(e)
}

func GetStudentsByIdOrEmail(db *sql.DB, email string, name string) (types.Student, error) {
	var student types.Student

	// Check if email is provided and not empty
	if email != "" {
		query := "SELECT name, email, id,password FROM students WHERE email=?"
		err := db.QueryRow(query, email).Scan(&student.Name, &student.Email, &student.Id, &student.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				return types.Student{}, fmt.Errorf("no student found with email: %v", email)
			}
			return types.Student{}, err
		}
		fmt.Printf("Found student with email: %v\n", student)
		return student, nil
	}

	// Check if name is provided and not empty
	if name != "" {
		query := "SELECT name, email, id,password FROM students WHERE name=?"
		err := db.QueryRow(query, name).Scan(&student.Name, &student.Email, &student.Id, &student.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				return types.Student{}, fmt.Errorf("no student found with name: %v", name)
			}
			return types.Student{}, err
		}
		fmt.Printf("Found student with name: %v\n", student)
		return student, nil
	}

	return types.Student{}, fmt.Errorf("either email or name must be provided")
}
