package model

// Success response helper
func Success(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// Error response helper
func Error(message string, err interface{}) BaseResponse {
	return BaseResponse{
		Status:  "error",
		Message: message,
		Error:   err,
	}
}
