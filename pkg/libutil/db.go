package libutil

// sqlite db connector
import (
	"database/sql"
	"fmt"
	"os"

	nblogger "github.com/banaconda/nb-logger"
)

// sqlite db connector
type SqliteConnector struct {
	logger nblogger.Logger
	db     *sql.DB
}

// open sqlite db
func (s *SqliteConnector) Open(dbPath string, logger nblogger.Logger) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	s.db = db

	s.logger = logger

	return nil
}

// close sqlite db
func (s *SqliteConnector) Close() error {
	return s.db.Close()
}

// query sqlite db
func (s *SqliteConnector) Query(query string, args ...any) (*sql.Rows, error) {
	s.logger.Info("query: %s", fmt.Sprintf(query, args...))
	return s.db.Query(fmt.Sprintf(query, args...))
}

// exec sqlite db
func (s *SqliteConnector) Exec(query string, args ...any) (sql.Result, error) {
	s.logger.Info("query: %s", fmt.Sprintf(query, args...))
	return s.db.Exec(fmt.Sprintf(query, args...))
}
