package tool

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func FormatJson(status int, err error, data interface{}, c *gin.Context) {
	//发生错误
	if err != nil {
		c.JSON(status, gin.H{
			"status": status,
			"err":    fmt.Sprint(err),
		})
		return
	}
	//无错误
	c.JSON(status, gin.H{
		"status": status,
		"data":   data,
	})
}
