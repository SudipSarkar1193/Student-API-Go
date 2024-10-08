package student

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/SudipSarkar1193/students-API-Go/internal/storage/mySql_Db"
	"github.com/SudipSarkar1193/students-API-Go/internal/types"
	"github.com/SudipSarkar1193/students-API-Go/internal/utils/response"

	"github.com/go-playground/validator/v10"
)

func New(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("%v HTTP method is not allowed", r.Method), http.StatusBadRequest)
			return
		}

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		// ⭐⭐ Explaination :

		/*
			1. r.Body contains the body of the HTTP request, typically in JSON format, that is sent to the server.

			2. json.NewDecoder(r.Body) creates a new JSON decoder to read and parse the JSON data from the r.Body.

			3. .Decode(&student) attempts to decode (unmarshal) the JSON data into the student struct. The &student is a pointer to the student struct, which means the decoded data will be stored directly into this struct.

			So, essentially, this line reads the JSON payload from the request body and decodes it into the Go struct named student.

		*/
		if err != nil {
			if errors.Is(err, io.EOF) {
				//io.EOF is a sentinel error in Go that indicates the end of input (end of a file or stream), commonly returned by functions when there is no more data to read.

				http.Error(w, fmt.Sprintf("no data to read: %v", err), http.StatusBadRequest)
				return

			} else {
				// Handle other decoding errors

				http.Error(w, fmt.Sprintf("failed to decode JSON: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Validate that all fields are filled
		// if student.Id == 0 || student.Name == "" || student.Email == "" {
		// 	http.Error(w, "all fields are required !", http.StatusBadRequest)
		// 	return

		// }

		var validate *validator.Validate

		validate = validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(&student); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err)
				return
			}

			response.ValidateResponse(w, err)
			return
		}

		//Everything is fine till now

		if err := mySql_Db.InsertStudent(db, student); err != nil {
			http.Error(w, fmt.Sprintf("Database error : %v", err), http.StatusInternalServerError)
			return
		}

		emptyResponse := response.CreateResponse(student, http.StatusCreated, "Student created Succesfully", "DeveloperMessage", "UserMessage", false, "Err")

		response.WriteResponse(w, emptyResponse)

	}
}

func GetAllStudents(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request ) {

		if r.Method != http.MethodGet {
			http.Error(w, fmt.Sprintf("%v HTTP method is not allowed", r.Method), http.StatusBadRequest)
			return
		}

		student, err := mySql_Db.GetAllStudents(db)
		if err != nil {
			http.Error(w, "Error fetching students data", http.StatusInternalServerError)
			return
		}

		emptyResponse := response.CreateResponse(student, http.StatusOK, "Student data fetched Succesfully", "DeveloperMessage", "UserMessage")

		response.WriteResponse(w, emptyResponse)

	}
}

func GetStudentsByIdOrEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("%v HTTP method is not allowed", r.Method), http.StatusBadRequest)
			return
		}

		type idOrEmail struct {
			Id    *int    `json:"id,omitempty"`
			Email *string `json:"email,omitempty"`
		}

		var reqBody idOrEmail

		// Decode the JSON request body
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		
		var student types.Student
		var err error
		if reqBody.Email != nil {
			student, err = mySql_Db.GetStudentsByIdOrEmail(db, *reqBody.Email, nil)
		} else if reqBody.Id != nil {
			student, err = mySql_Db.GetStudentsByIdOrEmail(db, nil, *reqBody.Id)
		}
		// fmt.Printf("reqBody.Id",*reqBody.Id," -> ",reqBody.Id)
		if err!=nil{
			if err == sql.ErrNoRows {
				if reqBody.Email!=nil {
					http.Error(w,"No student found with the email",http.StatusNotFound)
					return
				} else if reqBody.Id!=nil{
					http.Error(w,"No student found with the Id",http.StatusNotFound)
					return
				}else{
					http.Error(w,"Error fetching data",http.StatusBadRequest)
				}
			}else{
				http.Error(w,err.Error(),http.StatusBadRequest)
			}
		}
		
		if student.Id == 0 {
			student.Id = -1 // Or any placeholder value indicating it wasn't fetched
		}

		isSafe := r.Context().Value("IsSafe").(string)

		type respStruct struct {
			IsSafe any `json:"isSafe,omitempty"`
			types.Student
		}

		
		

		res := respStruct{
			IsSafe:  isSafe,
			Student: student, // This will copy all fields from the student to the response
		}

		


		emptyResponse := response.CreateResponse(res, http.StatusOK, "Student data fetched Succesfully", "DeveloperMessage", "UserMessage")

		response.WriteResponse(w, emptyResponse)


	}
}


func AddIsSafeMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Define logic to determine if the request is "safe"
		

		// Example: Check if a specific header exists
		// if r.Header.Get("X-Safe-Request") == "true" {
		// 	isSafe = true
		// }

		isSafeMsg := "Not Safe ! Na na ekkeibarei safe na !" // Default message
		// Add isSafe to the request context
		ctx := context.WithValue(r.Context(), "IsSafe",isSafeMsg)

		// Create a new request with the updated context
		r = r.WithContext(ctx)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

