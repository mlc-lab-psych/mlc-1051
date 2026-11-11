package data

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/mlc-project-template/internal/airtable"
	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/WeatherGod3218/mlc-project-template/internal/redis"
	"github.com/sirupsen/logrus"
)

// func replaceNullWithString(obj []map[string]interface{}) {
// 	for key, value := range obj {
// 		if value == nil {
// 			obj[key] = "null"
// 		} else if nested, ok := value.(map[string]interface{}); ok {
// 			replaceNullWithString(nested)
// 		}
// 	}
// }

func getAirtableData() (*airtable.SavedData, error) {
	tableName, err := redis.GetNextTable()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtableData"}).Warn("error in firebase deciding next table!")
		tableName = airtable.GetNextAirtable()
	}

	airTable := airtable.GetAirtableData(tableName)

	if airTable == nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtableData"}).Warn("error deciding which airtable to use!")
		return nil, fmt.Errorf("Failed to find the tablename %s", tableName)
	}

	return airTable, nil
}

func GetData(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), time.Second*30)
	defer cancel()

	data, err := getAirtableData()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetData"}).Fatal("error fetching airtable!")
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func Routes(r *gin.RouterGroup) {
	data := r.Group("/data")

	data.GET("", GetData)
}
