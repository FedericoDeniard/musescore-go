package customErrors

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func HandleError(c *gin.Context, err error) {
	if err == nil {
		c.Abort()
		return
	}

	if httpErr, ok := err.(*HttpError); ok {
		fmt.Println(httpErr)
		c.JSON(httpErr.StatusCode, gin.H{"error": httpErr.Message})
		c.Abort()
	} else {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Error interno del servidor"})
		c.Abort()
	}
}
