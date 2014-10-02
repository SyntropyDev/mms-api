package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/SyntropyDev/mms-api/model"
	"github.com/SyntropyDev/mms-api/mware"
	"github.com/SyntropyDev/sqlutil"
	"github.com/bmizerany/pat"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lann/squirrel"
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

	mware.SetGetDBConnectionFunc(db)

	http.Handle("/", m)
	log.Println("Listening...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func listenToFeeds() error {
	dbmap, err := db()
	if err != nil {
		return err
	}
	defer dbmap.Db.Close()

	feeds := []*model.Feed{}
	query := squirrel.Select("*").From(model.TableNameFeed)
	if err := sqlutil.Select(dbmap, query, &feeds); err != nil {
		return err
	}

	for _, feed := range feeds {
		if err := feed.UpdateStories(dbmap); err != nil {
			return err
		}
	}

	return nil
}

func db() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", "root@cloudsql(mms-api:db)/mms")
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

	if err := dbmap.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	return dbmap, nil
}
