package results

import (
	"context"
	"net/http"
	"time"

	"github.com/WeatherGod3218/mlc-project-template/internal/firebase"
	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SubmitResults(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*30)
	defer cancel()

	var results []map[string]any

	if err := c.ShouldBindJSON(&results); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "SubmitResults"}).Warn("error casting")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if err := firebase.PushToDatabase(ctx, results); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "SubmitResults"}).Warn("error updating database!")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving results.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully saved results!"})
}

func Routes(r *gin.RouterGroup) {
	data := r.Group("/results")

	data.PUT("", SubmitResults)
}
