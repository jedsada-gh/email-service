package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/email-service/data"
)

// PrintErrorMessage is response message error request
func PrintErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	var errorModel data.Error
	var errorDetailModel data.ErrorDetail
	errorDetailModel.Message = message
	errorModel.ErrorDetail = errorDetailModel
	error, err := json.Marshal(errorModel)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(error)
}

// PrintSuccessMessage is response message success request
func PrintSuccessMessage(w http.ResponseWriter, obj interface{}) {
	model, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(model)
}
