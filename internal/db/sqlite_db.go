package db

import (
	"database/sql"
	"errors"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/ii/xds-test-harness/internal/parser"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("Record already exists")
	ErrNotExists    = errors.New("Record does not exist")
	ErrUpdateFailed = errors.New("Update failed")
	ErrDeleteFailed = errors.New("Delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (s *SQLiteRepository) Migrate() error {
	_, err := s.db.Exec(MigrateSQL)
	return err
}

func (s *SQLiteRepository) InsertRequest(req *discovery.DiscoveryRequest) error {
	b, err := parser.ProtoJSONMarshal(req)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(InsertRequestSQL, string(b))
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return err
			}
		}
		return err
	}
	return err
}

func (s *SQLiteRepository) InsertResponse(res *discovery.DiscoveryResponse) error {
	b, err := parser.ProtoJSONMarshal(res)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(InsertResponseSQL, string(b))
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return err
			}
		}
		return err
	}
	return err
}
func (s *SQLiteRepository) DeleteAll() error {
	_, err := s.db.Exec(DeleteAllSQL)
	if err != nil {
		return fmt.Errorf("Couldn't delete the records: %v", err)
	}
	return nil
}

