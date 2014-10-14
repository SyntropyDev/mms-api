package mware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/mms-api/model"
	"github.com/SyntropyDev/sqlutil"
	"github.com/lann/squirrel"
)

func TopStoriesHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		v := r.URL.Query()

		// get members for category if specified
		memIDs := []string{}
		categoryID := v.Get("categoryId")
		if categoryID != "" {
			catMems := []*model.CategoryMember{}
			query := "select * from category_members where categoryID = ?"
			dbmap.Select(&catMems, query, categoryID)

			for _, catMem := range catMems {
				memIDs = append(memIDs, fmt.Sprint(catMem.MemberID))
			}
		}

		limit, err := strconv.ParseUint(v.Get("q-limit"), 10, 64)
		if err != nil {
			limit = 20
		}

		offset, _ := strconv.ParseUint(v.Get("q-offset"), 10, 64)

		query := squirrel.Select("*").From(model.TableNameStory).OrderBy("Score desc, Timestamp desc").Limit(limit).Offset(offset)
		if len(memIDs) > 0 {
			query = query.Where(squirrel.Eq{"memberId": memIDs})
		}

		stories := []*model.Story{}
		if err := sqlutil.Select(dbmap, query, &stories); err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(stories)
	}
}

func CommunityHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		community := &model.Community{}
		if err := sqlutil.SelectOneRelation(dbmap, model.TableNameCommunity, 1, community); err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(community)
	}
}

func ConstantHandler(src interface{}) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return json.NewEncoder(w).Encode(src)
	}
}

func ServeFile(name string) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		http.ServeFile(w, r, name)
		return nil
	}
}
