package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// func ResponseJSON(w http.ResponseWriter, data interface{}) {
// 	w.Header().Add("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(data)
// }

func ResponseJSON(c *gin.Context, data interface{}) {
	c.Writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(data)
}
