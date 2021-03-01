package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

func ConvertToJson(model interface{}) []byte {
	jsonData, err := json.Marshal(model)
	if err != nil {
		fmt.Println("ConvertToJson error: ", err)
	} else {
		fmt.Println("data: ", jsonData)
	}

	return jsonData
}

// JSONokResponse make a ok response with a json body
func JSONokResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	return
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
	return
}

// NotFoundResponse make a not found response with a error message
func NotFoundResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("No results"))
	return
}

// UnprocessableEntityResponse make an unprocessable entity response with a specific message
func UnprocessableEntityResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(message))
	return
}

// InternalServerError log the caller code position and write to w 500 Internal Server Error
func InternalServerError(w http.ResponseWriter) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Println(file, ":", line)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func convertStringToUint(s string) uint {
	var uint64 uint64
	uint64, _ = strconv.ParseUint(s, 10, 64)
	return uint(uint64)
}
