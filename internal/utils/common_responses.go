package utils

type Response struct {
	Message string `json:"message"`
}

var SuccessResponse = Response{
	Message: "success",
}

var BadRequestResponse = Response{
	Message: "Bad request",
}

var InternalServerResponse = Response{
	Message: "Internal server error occurred",
}
