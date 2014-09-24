package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

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

func main() {
	m := pat.New()
	m.Get(prefix+"/community", serveFile("public/community.json"))

	m.Get(prefix+"/members", mware.GetAll(&model.Member{}))
	m.Get(prefix+"/members/:id", mware.GetByID(&model.Member{}))

	m.Get(prefix+"/categories", mware.GetAll(&model.Category{}))
	m.Get(prefix+"/categories/:id", mware.GetByID(&model.Category{}))

	m.Get(prefix+"/stories", mware.GetAll(&model.Story{}))
	m.Get(prefix+"/stories/:id", mware.GetByID(&model.Story{}))

	http.Handle("/", m)

	mware.SetGetDBConnectionFunc(db)

	// listen to feeds in a background thread
	go func() {
		for {
			if err := listenToFeeds(); err != nil {
				log.Println("feed listening error: ", err)
			}
			time.Sleep(time.Minute * 5)
		}
	}()

	log.Println("Listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func listenToFeeds() error {
	dbmap, err := db()
	if err != nil {
		return err
	}

	feeds := []*model.Feed{}
	query := squirrel.Select("*").From(model.TableNameFeed)
	if err := sqlutil.Select(dbmap, query, &feeds); err != nil {
		return err
	}

	for _, feed := range feeds {
		feed.UpdateStories(dbmap)
	}

	return nil
}

func db() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
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

func serveFile(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}
