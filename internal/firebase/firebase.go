package firebase

import (
	"context"
	"fmt"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"

	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var database *db.Client

func InitFirebase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	options := option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")))

	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL: os.Getenv("FIREBASE_DATABASE_URL"),
	}, options)

	if err != nil {
		return err
	}

	database, err = app.Database(ctx)
	if err != nil {
		return err
	}

	logging.Logger.WithFields(logrus.Fields{"module": "firebase", "method": "initFirebase"}).Info("connected to firebase!")
	return nil
}

func PushToDatabase(ctx context.Context, cleanedData []map[string]any) error {
	id, err := uuid.NewV7()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "firebase", "method": "PushToDatabase"}).Warn("error generating the UUID7")
	}
	err = database.NewRef(fmt.Sprintf("users/new/%s", id.String())).Set(ctx, cleanedData)

	return err
}
