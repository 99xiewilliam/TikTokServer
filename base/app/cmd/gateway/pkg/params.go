package pkg

import (
	"github.com/gin-gonic/gin"
)

func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}
