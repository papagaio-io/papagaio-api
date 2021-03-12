package api

func IsResponseOK(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
