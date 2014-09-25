package cloud

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/sqlutil"
	"github.com/bmizerany/pat"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lann/squirrel"
)

const (
	prefix = "/api/v1"
)

func init() {

	SetGetDBConnectionFunc(db)

	m := pat.New()
	m.Get(prefix+"/community", serveFile("config/community.json"))

	m.Get(prefix+"/members", GetAll(&Member{}))
	m.Get(prefix+"/members/:id", GetByID(&Member{}))

	m.Get(prefix+"/feeds", GetAll(&Feed{}))
	m.Get(prefix+"/feeds/:id", GetByID(&Feed{}))

	m.Get(prefix+"/categories", GetAll(&Category{}))
	m.Get(prefix+"/categories/:id", GetByID(&Category{}))

	m.Get(prefix+"/stories", GetAll(&Story{}))
	m.Get(prefix+"/stories/:id", GetByID(&Story{}))

	m.Get("/tasks/feeds", httperr.Handler(func(w http.ResponseWriter, r *http.Request) error {
		return listenToFeeds()
	}))

	http.Handle("/", m)

	// log.Println("Listening...")
	// err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	// if err != nil {
	// 	panic(err)
	// }
}

func listenToFeeds() error {
	dbmap, err := db()
	if err != nil {
		return err
	}
	defer dbmap.Db.Close()

	feeds := []*Feed{}
	query := squirrel.Select("*").From(TableNameFeed)
	if err := sqlutil.Select(dbmap, query, &feeds); err != nil {
		return err
	}

	for _, feed := range feeds {
		return feed.UpdateStories(dbmap)
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

	dbmap.AddTableWithName(Category{}, TableNameCategory).SetKeys(true, "ID")
	dbmap.AddTableWithName(Feed{}, TableNameFeed).SetKeys(true, "ID")
	dbmap.AddTableWithName(Story{}, TableNameStory).SetKeys(true, "ID")
	dbmap.AddTableWithName(Member{}, TableNameMember).SetKeys(true, "ID")

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
