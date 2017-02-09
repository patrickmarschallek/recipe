package recipe

import (
	"database/sql"
	"log"
	"os"
	"recipe/config"

	"github.com/uber-go/zap"
)

// type Env struct {
//     db *sql.DB
//     logger *log.Logger
//     templates *template.Template
// }

// type DB struct {
//     *sql.DB
// }

// func NewDB(dataSourceName string) (*DB, error) {
//     db, err := sql.Open("postgres", dataSourceName)
//     if err != nil {
//         return nil, err
//     }
//     if err = db.Ping(); err != nil {
//         return nil, err
//     }
//     return &DB{db}, nil
// }

// logger internal logger represantation.
var logger zap.Logger

// Logger default framework logger.
var Logger zap.Logger

// DBCon database connection
var DBCon *sql.DB

func init() {
	logger = zap.New(
		zap.NewJSONEncoder(
			zap.RFC3339Formatter("@timestamp"), // human-readable timestamps
			zap.MessageKey("@message"),         // customize the message key
			zap.LevelString("@level"),          // stringify the log level
		),
		zap.WarnLevel,
		zap.Fields(
			zap.Int("pid", os.Getpid()),
		),
	)
	Logger = logger
}

// NewLogger adds the logger.
func NewLogger(log zap.Logger) {
	logger = log
	Logger = log
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
