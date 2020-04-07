package postgresql

import (
	"fmt"

	config "github.com/duyhtq/incognito-data-sync/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PGStoreAbs is PDE State postgresql store abstraction
type PGStoreAbs struct {
	DB *sqlx.DB
}

func Init(config *config.Config) (*sqlx.DB, error) {
	pgConn, err := sqlx.Open("postgres", config.Db)
	if err != nil {
		fmt.Println("An error occured while openning connection to pg")
		return nil, err
	}
	// fmt.Sprintf(`%s %s %s`, 1, 2, 2)
	return pgConn, nil
}
