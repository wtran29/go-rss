package driver

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

func ConnectPostgres(dsn string) (*DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	dbConn.SQL = conn

	err = testDB(err, conn)

	return dbConn, err
}

// testDB pings database
func testDB(err error, d *sql.DB) error {
	err = d.Ping()
	if err != nil {
		fmt.Println("Error!", err)
	} else {
		log.Println("*** Pinged database successfully! ***")
	}
	return err
}
