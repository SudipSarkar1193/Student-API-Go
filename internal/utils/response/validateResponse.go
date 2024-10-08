package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

func ValidateResponse(w http.ResponseWriter, errs interface{}) error {
	for _, err := range errs.(validator.ValidationErrors) {

		fmt.Println()
		fmt.Println()
		fmt.Println("err.Namespace():", err.Namespace())
		fmt.Println("err.Field():", err.Field())
		fmt.Println("err.StructNamespace():", err.StructNamespace())
		fmt.Println("err.StructField():", err.StructField())
		fmt.Println("err.Tag():", err.Tag())
		fmt.Println("err.ActualTag():", err.ActualTag())
		fmt.Println("err.Kind():", err.Kind())
		fmt.Println("err.Type():", err.Type())
		fmt.Println("err.Value():", err.Value())
		fmt.Println("err.Param():", err.Param())
		fmt.Println()
		fmt.Println()

		var errMsgs []string

		switch err.ActualTag() {
		case "email":
			errMsgs = append(errMsgs, "invalid email")
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", err.Field()))

		default:
			errMsgs = append(errMsgs, "Invalid Input !")
		}

		http.Error(w, strings.Join(errMsgs, ","), http.StatusInternalServerError)

	}

	return nil
}
