package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/SudipSarkar1193/students-API-Go/internal/types"
	"github.com/SudipSarkar1193/students-API-Go/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "this HTTP method is not allowed !!!!!!!", http.StatusBadRequest)
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

				emptyResponse := response.Response{
					Data:         nil, // No data to return
					Message:      "there is no more data to read ",
					ErrorMessage: fmt.Sprintf("no data to read: %v", err),
					StatusCode:   http.StatusBadRequest, // 400
					ErrorCode:    http.StatusNoContent,  // 204
					IsError:      true,
				}

				response.WriteResponse(w, emptyResponse)

				return
			} else {
				// Handle other decoding errors
				emptyResponse := response.Response{
					Data:         nil,
					Message:      "invalid request body",
					ErrorMessage: fmt.Sprintf("failed to decode JSON: %v", err),
					StatusCode:   http.StatusBadRequest,
					ErrorCode:    http.StatusBadRequest,
					IsError:      true,
				}
				response.WriteResponse(w, emptyResponse)
				return
			}
		}

		// Validate that all fields are filled
		if student.Id == 0 || student.Name == "" || student.Email == "" {
			emptyResponse := response.Response{
				Data:         nil,
				Message:      "all fields are required",
				ErrorMessage: "id, name, and email must be filled",
				StatusCode:   http.StatusBadRequest,
				ErrorCode:    http.StatusBadRequest,
				IsError:      true,
			}
			response.WriteResponse(w, emptyResponse)

			return

		}

		//Everything is fine till now

		emptyResponse := response.Response{
			Data:       student,
			Message:    "Student created",
			StatusCode: http.StatusCreated,
			IsError:    false,
		}

		response.WriteResponse(w, emptyResponse)

	}
}
