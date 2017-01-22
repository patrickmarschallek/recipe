package recipe

import (
	"database/sql"
	"log"
	"recipe/config"

	"github.com/uber-go/zap"
)

// Logger default framework logger.
var Logger zap.Logger

// DBCon database connection
var DBCon *sql.DB

// initialize all framework dependent things.
func init() {
	Logger = zap.New(zap.NewJSONEncoder(),
		zap.DebugLevel,
		zap.AddCaller(),
	)
}

// GetDBConn return the database connection.
func GetDBConn(conf *config.Config) *sql.DB {
	// extract database connection...
	DBCon, err := conf.Persistence.DBConnection()
	if err != nil {
		log.Fatal("error establishing a database connection", zap.Error(err))
	}
	return DBCon
}
