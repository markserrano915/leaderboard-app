package function

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/openfaas/openfaas-cloud/sdk"
)

var db *sql.DB

func init() {

	password, _ := sdk.ReadSecret("password")
	user, _ := sdk.ReadSecret("username")
	host, _ := sdk.ReadSecret("host")
	dbName := os.Getenv("postgres_db")
	port := os.Getenv("postgres_port")
	sslmode := os.Getenv("postgres_sslmode")

	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=" + sslmode

	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {

	rows, getErr := db.Query(`select * from get_leaderboard();`)

	if getErr != nil {
		log.Printf("get error: %s", getErr.Error())
		return handler.Response{
			Body:       []byte(errors.Wrap(getErr, "unable to get from leaderboard")),
			StatusCode: http.StatusInternalServerError,
		}, updateErr
	}

	results := []Result{}
	defer rows.Close()
	for rows.Next() {
		result := Result{}
		scanErr := rows.Scan(&result.UserID, &result.UserLogin, &result.IssueComments, &result.IssuesCreated)
		if scanErr != nil {
			log.Println("scan err:", scanErr)
		}
		results = append(results, result)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(results)
	w.Write(res)
}

type Result struct {
	UserID    int
	UserLogin string

	IssueComments int
	IssuesCreated int
}
