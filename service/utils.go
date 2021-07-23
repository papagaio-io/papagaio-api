package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

func ConvertToJson(model interface{}) []byte {
	jsonData, err := json.Marshal(model)
	if err != nil {
		log.Println("ConvertToJson error: ", err)
	} else {
		log.Println("data: ", jsonData)
	}

	return jsonData
}

// JSONokResponse make a ok response with a json body
func JSONokResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Println("encode error:", err)
	}
}

// MakeResponse prepare a response
func MakeResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	switch statusCode {
	case http.StatusNotFound:
		NotFoundResponse(w)
	case http.StatusUnprocessableEntity:
		message := "UnprocessableEntity"
		UnprocessableEntityResponse(w, message)
	case http.StatusOK:
		JSONokResponse(w, response)
	default:
		InternalServerError(w)
	}
}

// NotFoundResponse make a not found response with a error message
func NotFoundResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("No results"))

	if err != nil {
		log.Println("NotFoundResponse parsing error:", err)
	}
}

// UnprocessableEntityResponse make an unprocessable entity response with a specific message
func UnprocessableEntityResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println("UnprocessableEntityResponse parsing error:", err)
	}
}

// InternalServerError log the caller code position and write to w 500 Internal Server Error
func InternalServerError(w http.ResponseWriter) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Println(file, ":", line)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func getBoolParameter(r *http.Request, parameterName string) (bool, error) {
	forceCreateQuery, ok := r.URL.Query()[parameterName]
	forceCreate := false
	if ok {
		if len(forceCreateQuery[0]) == 0 {
			forceCreate = true
		} else {
			var parsError error
			forceCreate, parsError = strconv.ParseBool(forceCreateQuery[0])
			if parsError != nil {
				return false, errors.New(parameterName + " is not valid")
			}
		}
	}

	return forceCreate, nil
}
