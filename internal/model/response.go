package model

// Success response helper
func Success(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  "Success",
		Message: message,
		Data:    data,
	}
}

// Error response helper
func Error(message string, err interface{}) BaseResponse {
	return BaseResponse{
		Status:  "Error",
		Message: message,
		Error:   err,
	}
}
