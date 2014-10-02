package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SyntropyDev/mms-api/model"
	"github.com/SyntropyDev/mms-api/mware"
	"github.com/bmizerany/pat"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

const (
	prefix = "/api/v1"
)

var (
	community = map[string]interface{}{
		"object":      "Community",
		"name":        "Test",
		"location":    []float64{39.607672, -79.958496},
		"url":         "",
		"description": "",
	}
)

func main() {
	initConfig()
	dbmap, err := db()
	if err != nil {
		panic(err)
	}
	if err := dbmap.CreateTablesIfNotExists(); err != nil {
		panic(err)
	}
	mware.SetGetDBConnectionFunc(db)

	m := pat.New()
	m.Get(prefix+"/community", mware.ConstantHandler(&community))

	m.Get(prefix+"/members", mware.GetAll(&model.Member{}))
	m.Get(prefix+"/members/:id", mware.GetByID(&model.Member{}))

	m.Get(prefix+"/feeds", mware.GetAll(&model.Feed{}))
	m.Get(prefix+"/feeds/:id", mware.GetByID(&model.Feed{}))

	m.Get(prefix+"/categories", mware.GetAll(&model.Category{}))
	m.Get(prefix+"/categories/:id", mware.GetByID(&model.Category{}))

	m.Get(prefix+"/stories", mware.GetAll(&model.Story{}))
	m.Get(prefix+"/stories/:id", mware.GetByID(&model.Story{}))

	go listenToFeedsInBackground()

	http.Handle("/", m)
	log.Println("Listening...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func listenToFeedsInBackground() {
	for {
		err := func() error {
			dbmap, err := db()
			if err != nil {
				return err
			}
			defer dbmap.Db.Close()

			return model.ListenToFeeds(dbmap)
		}()
		if err != nil {
			fmt.Println("LISTEN TO FEEDS: ", err)
		}
		time.Sleep(time.Minute * 10)
	}
}

func initConfig() error {
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}

	config := map[string]string{}
	if err := json.Unmarshal(b, &config); err != nil {
		return err
	}

	for key, value := range config {
		os.Setenv(key, value)
	}
	return nil
}

func db() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", os.Getenv("mysql"))
	if err != nil {
		return nil, err
	} else {
		log.Println("Database Connection Established")
	}

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"},
	}

	dbmap.AddTableWithName(model.Category{}, model.TableNameCategory).SetKeys(true, "ID")
	dbmap.AddTableWithName(model.Feed{}, model.TableNameFeed).SetKeys(true, "ID")
	dbmap.AddTableWithName(model.Story{}, model.TableNameStory).SetKeys(true, "ID")
	dbmap.AddTableWithName(model.Member{}, model.TableNameMember).SetKeys(true, "ID")

	return dbmap, nil
}
