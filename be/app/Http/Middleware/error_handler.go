package Middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ErrorHandle() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		ok := errors.Cause(err.Err)
		if ok != nil {
			c.JSON(400, gin.H{
				"error": "Blah blahhh",
			})
			return
		}
	}
}
