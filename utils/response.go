package utils

import "github.com/gin-gonic/gin"

type Response struct {
    Code    int         `json:"code"`
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Meta    interface{} `json:"meta,omitempty"` 
}

func Success(c *gin.Context, message string, data interface{}, meta interface{}) {
    c.JSON(200, Response{
        Code:    200,
        Status:  "success",
        Message: message,
        Data:    data,
        Meta:    meta,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Code:    code,
        Status:  "error",
        Message: message,
    })
}
