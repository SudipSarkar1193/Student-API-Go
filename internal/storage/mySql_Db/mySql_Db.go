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
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE
    );`

	_, err = dbWithDB.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	fmt.Println("Table 'students' created or already exists.")

	return dbWithDB, nil

}

// Function to insert a new student into the database
func InsertStudent(db *sql.DB, student types.Student) error {
	query := "INSERT INTO students (name, email) VALUES (?, ?)"

	// Executing the query
	result, err := db.Exec(query, student.Name, student.Email)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted student
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("New student inserted with ID: %d\n", id)
	return nil
}

// Function to insert a new student into the database
func GetAllStudents(db *sql.DB) ([]types.Student, error) {
	query := "SELECT * FROM student.students"
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

func GetStudentsByIdOrEmail(db *sql.DB,email interface{} , id interface{}) (types.Student, error) {
	
	
	
	var student types.Student;

	if email != nil {
		query := "SELECT name,email,id FROM student.students WHERE email=?"
		err := db.QueryRow(query,email).Scan(&student.Name,&student.Email,&student.Id)
		if err != nil {
            return types.Student{},err
        }
		
		return student,nil

	}else if id!=nil {
		query := "SELECT name,email,id FROM student.students WHERE id=?"
		err := db.QueryRow(query,id).Scan(&student.Name,&student.Email,&student.Id)
		if err != nil {
            return types.Student{},err
        }
		
		return student,nil
	}else{
		
			return types.Student{}, Error("Either ID or Email must be provided")
	}
}
