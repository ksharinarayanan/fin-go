package utils

type RequestResponse struct {
	Result string `json:"message"`
}

var BadRequestResponse = RequestResponse{
	Result: "Bad request",
}

var InternalServerResponse = RequestResponse{
	Result: "Internal server error occurred",
}
