package httpresponse

import "github.com/gin-gonic/gin"

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func SentErrorResponse(c *gin.Context, status int, message string, err error) {
	var resp Response
	if err != nil {
		resp = Response{
			Message: message,
			Error:   err.Error(),
		}
	} else {
		resp = Response{
			Message: message,
		}
	}

	c.AbortWithStatusJSON(status, resp)
}
