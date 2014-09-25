package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/SyntropyDev/mms-api/model"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

const (
	colName = iota
	colAddress
	colLatitude
	colLongitude
	colPhone
	colWebsite
	colEmail
	colFacebook
	colTwitter
	colRSS
	colDescription
)

func main() {
	dbmap, err := db()
	if err != nil {
		panic(err)
	}
	trans, err := dbmap.Begin()
	if err != nil {
		panic(err)
	}

	csvFile, err := os.Open("tucker.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		lat, err := strconv.ParseFloat(record[colLatitude], 64)
		if err != nil {
			panic(err)
		}
		lng, err := strconv.ParseFloat(record[colLongitude], 64)
		if err != nil {
			panic(err)
		}

		member := &model.Member{
			Name:        record[colName],
			Address:     record[colAddress],
			Phone:       record[colPhone],
			Description: record[colDescription],
			Website:     record[colWebsite],
			Latitude:    lat,
			Longitude:   lng,
		}
		if err := trans.Insert(member); err != nil {
			panic(err)
		}
		fmt.Println(member.Name, "Inserted")

		if record[colFacebook] != "" {
			feed := &model.Feed{
				MemberID:   member.ID,
				Type:       string(model.FeedTypeFacebook),
				Identifier: record[colFacebook],
			}
			if err := trans.Insert(feed); err != nil {
				panic(err)
			}
			fmt.Println(feed.Type, feed.Identifier, "Inserted")
		}
		if record[colTwitter] != "" {
			feed := &model.Feed{
				MemberID:   member.ID,
				Type:       string(model.FeedTypeTwitter),
				Identifier: record[colTwitter],
			}
			if err := trans.Insert(feed); err != nil {
				panic(err)
			}
			fmt.Println(feed.Type, feed.Identifier, "Inserted")
		}
		if record[colRSS] != "" {
			feed := &model.Feed{
				MemberID:   member.ID,
				Type:       string(model.FeedTypeRSS),
				Identifier: record[colRSS],
			}
			if err := trans.Insert(feed); err != nil {
				panic(err)
			}
			fmt.Println(feed.Type, feed.Identifier, "Inserted")
		}
		// out the csv content
	}

	if err := trans.Commit(); err != nil {
		panic(err)
	}
}

func db() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", "root:gDwrFNfDYLzkQTugaoGZEh7pYnjZAXr8@tcp(173.194.255.242:3306)/mms")
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
