package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"upsilon_garden_go/config"

	"github.com/lib/pq" // needed for postgres driver
)

// Handler Contains DB related informations
type Handler struct {
	db   *sql.DB
	open bool
}

// New Create a new handler for database, ensure database is created
func New() *Handler {
	handler := new(Handler)
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_HOST)
	db, err := sql.Open("postgres", dbinfo)

	if err, ok := err.(*pq.Error); ok {
		log.Printf("DB: Database failed to be connected: %s", err)
		dbinfo := fmt.Sprintf("user=%s password=%s sslmode=disable host=%s",
			config.DB_USER, config.DB_PASSWORD, config.DB_HOST)
		db, err := sql.Open("postgres", dbinfo)
		if err, ok := err.(*pq.Error); ok {
			log.Fatalf("DB: Can't initiate database creation, aborting : %s", err)
		} else {
			handler.db = db
			log.Print("DB: Creating Database !")

			handler.Exec("CREATE DATABSE %s", config.DB_NAME)

			log.Print("DB: Seeding Database !")
			data, err := ioutil.ReadFile(config.DB_SEED)

			if err != nil {
				log.Fatalf("DB: Failed to read seed file at %s : %s", config.DB_SEED, err)
			}

			// crude way of doing it ... can't have ';' anywhere else than end of query.
			requests := strings.Split(string(data), ";")

			for _, request := range requests {
				handler.db.Exec(request)
			}

		}
	} else {
		log.Printf("DB: Successfully connected to : %s %s", config.DB_HOST, config.DB_NAME)
	}

	handler.db = db
	handler.open = true
	return handler
}

// Exec executes provided query and check if it's correctly executed or not.
// Abort app if not.
func (dbh *Handler) Exec(format string, a ...interface{}) (result *sql.Rows) {
	dbh.CheckState()
	query := fmt.Sprintf(format, a)
	result, err := dbh.db.Query(query)
	errorCheck(query, err)
	return result
}

// Query Just like Exec but uses Postgres formater.
func (dbh *Handler) Query(format string, a ...interface{}) (result *sql.Rows) {
	dbh.CheckState()
	result, err := dbh.db.Query(format, a)
	errorCheck(format, err)
	return result
}

// CheckState assert that connection to DB is still alive. or break
func (dbh *Handler) CheckState() {
	if !dbh.open {
		log.Fatal("DB: Can't use this connection, it's been closed")
	}
	err := dbh.db.Ping()
	if err != nil {
		log.Fatalf("DB: Can't use this connection, an error occured: %s", err)
	}
}

// Close frees db ressource
func (dbh *Handler) Close() {
	if dbh.open {
		dbh.open = false
		defer dbh.db.Close()
	} else {
		log.Print("DB: Already Closed")
	}
}

// Drop database and close connection
func (dbh *Handler) Drop() {
	query := fmt.Sprintf("DROP DATABASE %s IF EXISTS", config.DB_NAME)
	dbh.Exec(query)
	defer dbh.Close()
}

// ErrorCheck checks if query result has an error or not
func errorCheck(query string, err error) bool {
	if err != nil {
		log.Printf("DB: Failed to execute query: %s", query)

		// fatal aborts app
		log.Fatalf("DB: Aborting: %s", err)

		return true
	}

	return false
}
