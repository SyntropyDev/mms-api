package mware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/mms-api/model"
	"github.com/SyntropyDev/sqlutil"
)

func StoryQueryHandler(h httperr.Handler) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		categoryID := r.URL.Query().Get("categoryID")
		if categoryID != "" {
			catMems := []*model.CategoryMember{}
			query := "select * from category_members where categoryID = ?"
			dbmap.Select(&catMems, query, categoryID)

			memIDs := []string{}
			for _, catMem := range catMems {
				memIDs = append(memIDs, fmt.Sprint(catMem.MemberID))
			}
			r.URL.Query().Del("memberId")
			r.URL.Query().Set("memberId", "in-"+strings.Join(memIDs, ","))
		}

		return h(w, r)
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
