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

func main() {
	initConfig()
	initSQL()

	mware.SetGetDBConnectionFunc(db)

	m := pat.New()

	// no auth routes
	m.Get(prefix+"/community", mware.CommunityHandler())

	m.Post(prefix+"/login", mware.LoginHandler())
	m.Post(prefix+"/logout", mware.LogoutHandler())
	m.Post(prefix+"/signup", mware.SignupHandler())
	m.Post(prefix+"/reset-password", mware.ResetPasswordHandler())
	// m.Post(prefix+"/request-invite", mware.RequestInviteHandler())

	m.Get(prefix+"/members", mware.GetAll(&model.Member{}))
	m.Get(prefix+"/members/:id", mware.GetByID(&model.Member{}))

	m.Get(prefix+"/feeds", mware.GetAll(&model.Feed{}))
	m.Get(prefix+"/feeds/:id", mware.GetByID(&model.Feed{}))

	m.Get(prefix+"/categories", mware.GetAll(&model.Category{}))
	m.Get(prefix+"/categories/:id", mware.GetByID(&model.Category{}))

	m.Get(prefix+"/top-stories", mware.TopStoriesHandler())
	m.Get(prefix+"/stories", mware.GetAll(&model.Story{}))
	m.Get(prefix+"/stories/:id", mware.GetByID(&model.Story{}))

	// auth routes
	m.Post(prefix+"/invite", mware.Auth(mware.InviteHandler()))
	m.Post(prefix+"/change-password", mware.Auth(mware.ChangePasswordHandler()))

	m.Post(prefix+"/members", mware.Auth(mware.Create(&model.Member{})))
	m.Put(prefix+"/members/:id", mware.Auth(mware.UpdateByID(&model.Member{})))
	m.Del(prefix+"/members/:id", mware.Auth(mware.DeleteByID(&model.Member{})))

	m.Post(prefix+"/feeds", mware.Auth(mware.Create(&model.Feed{})))
	m.Put(prefix+"/feeds/:id", mware.Auth(mware.UpdateByID(&model.Feed{})))
	m.Del(prefix+"/feeds/:id", mware.Auth(mware.DeleteByID(&model.Feed{})))

	m.Post(prefix+"/categories", mware.Auth(mware.Create(&model.Category{})))
	m.Put(prefix+"/categories/:id", mware.Auth(mware.UpdateByID(&model.Category{})))
	m.Del(prefix+"/categories/:id", mware.Auth(mware.DeleteByID(&model.Category{})))

	m.Del(prefix+"/stories/:id", mware.Auth(mware.DeleteByID(&model.Story{})))

	// cors
	m.Options(prefix+"/:any", mware.CommunityHandler())
	m.Options(prefix+"/:any1/:any2", mware.CommunityHandler())

	go runInBackground(time.Minute*10, model.ListenToFeeds)
	go runInBackground(time.Minute*5, model.DecayScores)

	http.Handle("/", m)
	log.Println("Listening...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func runInBackground(d time.Duration, f func(s gorp.SqlExecutor) error) {
	for {
		err := func() error {
			dbmap, err := db()
			if err != nil {
				return err
			}
			defer dbmap.Db.Close()

			return f(dbmap)
		}()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		time.Sleep(d)
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
	dbmap.AddTableWithName(model.Community{}, model.TableNameCommunity).SetKeys(true, "ID")
	dbmap.AddTableWithName(model.Token{}, model.TableNameToken).SetKeys(true, "ID")
	dbmap.AddTableWithName(model.CategoryMember{}, model.TableNameCategoryMember)

	return dbmap, nil
}

func initSQL() error {
	db, err := sql.Open("mysql", os.Getenv("mysql"))
	if err != nil {
		return err
	}
	if _, err := db.Exec(sqlCreateCommunity); err != nil {
		return err
	}
	if _, err := db.Exec(sqlCreateMembers); err != nil {
		return err
	}
	if _, err := db.Exec(sqlCreateCategories); err != nil {
		return err
	}
	if _, err := db.Exec(sqlCreateFeeds); err != nil {
		return err
	}
	if _, err := db.Exec(sqlCreateStories); err != nil {
		return err
	}
	if _, err := db.Exec(sqlCreateTokens); err != nil {
		return err
	}
	if _, err := db.Exec(sqlCreateCategoryMembers); err != nil {
		return err
	}
	return nil
}

const (
	sqlCreateCommunity = `
	CREATE TABLE communities(
		ID bigint(20) NOT NULL AUTO_INCREMENT,
		Created bigint(20) NOT NULL,
		Updated bigint(20) NOT NULL,
		Deleted tinyint(1) NOT NULL,

		Name varchar(255) NOT Null,
		Latitude double Not Null,
		Longitude double Not Null,
		Description text Not Null,
		PRIMARY KEY (ID)
	);`

	sqlCreateMembers = `
	CREATE TABLE members(
		ID bigint(20) NOT NULL AUTO_INCREMENT,
		Created bigint(20) NOT NULL,
		Updated bigint(20) NOT NULL,
		Deleted tinyint(1) NOT NULL,
		
		Email varchar(255) NOT Null,
		Organizer tinyint(1) NOT NULL,
		PasswordHash varchar(255) NOT NULL,

		Name varchar(255) NOT Null,
		Address varchar(255) NOT Null,
		Phone varchar(255) NOT Null,
		Description text NOT Null,
		Icon varchar(255) NOT Null,
		Website varchar(255) NOT Null,
		Latitude double Not Null,
		Longitude double Not Null,
		ImagesRaw text Not Null,
		HashtagsRaw text Not Null,

		PRIMARY KEY (ID),
		UNIQUE (Email)
	);`

	sqlCreateCategories = `
	CREATE TABLE categories(
		ID bigint(20) NOT NULL AUTO_INCREMENT,
		Created bigint(20) NOT NULL,
		Updated bigint(20) NOT NULL,
		Deleted tinyint(1) NOT NULL,

		Name varchar(255) NOT Null,
		PRIMARY KEY (ID),
		UNIQUE (Name)
	);`

	sqlCreateFeeds = `
	CREATE TABLE feeds(
		ID bigint(20) NOT NULL AUTO_INCREMENT,
		Created bigint(20) NOT NULL,
		Updated bigint(20) NOT NULL,
		Deleted tinyint(1) NOT NULL,
		
		MemberID bigint(20) NOT Null,
		Type varchar(255) NOT Null,
		Identifier varchar(255) NOT Null,
		LastRetrieved bigint(20) NOT NULL,

		PRIMARY KEY (ID),
		FOREIGN KEY (MemberID) REFERENCES members(ID),
		UNIQUE (Type, Identifier)
	);`

	sqlCreateStories = `
	CREATE TABLE stories(
		ID bigint(20) NOT NULL AUTO_INCREMENT,
		Created bigint(20) NOT NULL,
		Updated bigint(20) NOT NULL,
		Deleted tinyint(1) NOT NULL,
		
		MemberID bigint(20) NOT Null,
		MemberName varchar(255) NOT Null,
		FeedID bigint(20) NOT Null,
		FeedIdentifier varchar(255) NOT Null,
		Timestamp bigint(20) NOT NULL,
		FeedType varchar(255) NOT Null,
		Body text NOT Null,
		SourceURL varchar(255) NOT Null,
		SourceID varchar(255) NOT Null,
		Score double Not Null,
		Latitude double Not Null,
		Longitude double Not Null,
		LinksRaw text NOT Null,
		ImagesRaw text NOT Null,
		HashtagsRaw text NOT Null,
		LastDecayTimestamp bigint(20) NOT NULL,

		PRIMARY KEY (ID),
		FOREIGN KEY (MemberID) REFERENCES members(ID),
		FOREIGN KEY (FeedID) REFERENCES feeds(ID),
		UNIQUE (Timestamp),
		UNIQUE (SourceID)
	);`

	sqlCreateTokens = `
	CREATE TABLE tokens(
		ID bigint(20) NOT NULL AUTO_INCREMENT,
		Created bigint(20) NOT NULL,
		Updated bigint(20) NOT NULL,
		Deleted tinyint(1) NOT NULL,
		
		MemberID bigint(20) NOT Null,
		Value varchar(255) NOT Null,
		Expiration bigint(20) NOT NULL,

		PRIMARY KEY (ID),
		FOREIGN KEY (MemberID) REFERENCES members(ID)
	);`

	sqlCreateCategoryMembers = `
	CREATE TABLE category_members(
		CategoryID bigint(20) NOT NULL,
		MemberID bigint(20) NOT NULL,
		
		FOREIGN KEY (CategoryID) REFERENCES categories(ID),
		FOREIGN KEY (MemberID) REFERENCES members(ID)
	);`
)
