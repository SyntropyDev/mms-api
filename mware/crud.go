package mware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/merge"
	"github.com/SyntropyDev/querystr"
	"github.com/SyntropyDev/sqlutil"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"github.com/lann/squirrel"
)

const (
	KeyFields = "q-fields"
)

type CrudResource interface {
	TableName() string
	TableId() int64
	Delete()
}

type GetDBConnectionFunc func() (*gorp.DbMap, error)

var (
	getDB GetDBConnectionFunc = func() (*gorp.DbMap, error) {
		panic("failed to register db w/ crud")
	}
)

func SetGetDBConnectionFunc(f GetDBConnectionFunc) {
	getDB = f
}

func GetAll(m CrudResource) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		values := r.URL.Query()
		sql, args, err := querystr.Query(m, m.TableName(), values)
		if err != nil {
			return clientError(err)
		}

		models, err := dbmap.Select(m, sql, args...)
		if err != nil {
			return clientError(err)
		}

		return getAllWriteJSON(w, values, models)
	}
}

func GetByID(m CrudResource) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		params := mux.Vars(r)
		mCopy := copyResource(m)
		if err := GetID(dbmap, mCopy, params["id"]); err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(mCopy)
	}
}

func Create(m CrudResource) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		trans, err := dbmap.Begin()
		if err != nil {
			return err
		}

		mCopy := copyResource(m)
		if err := json.NewDecoder(r.Body).Decode(mCopy); err != nil {
			return clientError(err)
		}

		if err := trans.Insert(mCopy); err != nil {
			message := fmt.Sprintf("%s did not pass validation.", m.TableName())
			return httperr.New(http.StatusBadRequest, message, err)
		}

		if err := trans.Commit(); err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(mCopy)
	}
}

func UpdateByID(m CrudResource) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		params := mux.Vars(r)
		mCopy := copyResource(m)
		if err := GetID(dbmap, mCopy, params["id"]); err != nil {
			return err
		}

		updateCopy := copyResource(m)
		if err := json.NewDecoder(r.Body).Decode(updateCopy); err != nil {
			message := fmt.Sprintf("%s's json could not parsed.", m.TableName())
			return httperr.New(http.StatusBadRequest, message, err)
		}

		merge.TagWl(updateCopy, mCopy)

		if _, err := dbmap.Update(mCopy); err != nil {
			message := fmt.Sprintf("%s did not pass validation.", m.TableName())
			return httperr.New(http.StatusBadRequest, message, err)
		}

		return json.NewEncoder(w).Encode(mCopy)
	}
}

func DeleteByID(m CrudResource) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		params := mux.Vars(r)
		mCopy := copyResource(m)
		if err := GetID(dbmap, mCopy, params["id"]); err != nil {
			return err
		}

		mCopy.Delete()

		if _, err := dbmap.Update(mCopy); nil != err {
			return err
		}

		return json.NewEncoder(w).Encode(mCopy)
	}
}

func GetID(dbmap *gorp.DbMap, m CrudResource, id interface{}) error {
	query := squirrel.Select("*").
		From(m.TableName()).
		Where(squirrel.Eq{"ID": id})
	if err := sqlutil.SelectOne(dbmap, query, m); err != nil {
		message := fmt.Sprintf("Could not find %s.", m.TableName())
		return httperr.New(http.StatusNotFound, message, err)
	}
	return nil
}

func clientError(err error) error {
	message := "Problem performing request.  Please alert the Account owner if the problem continues."
	return httperr.New(http.StatusBadRequest, message, err)
}

func copyResource(m CrudResource) CrudResource {
	ptr := reflect.New(reflect.TypeOf(m).Elem())
	iFace := ptr.Interface().(CrudResource)
	return iFace
}

func getJsonKeyFromTag(tag string) string {
	elems := strings.Split(tag, ",")
	for _, s := range elems {
		value := strings.TrimSpace(s)
		if value != "omitempty" {
			return value
		}
	}
	return "-"
}

func getAllWriteJSON(w http.ResponseWriter, values url.Values, models []interface{}) error {
	if len(models) == 0 || values.Get(KeyFields) == "" {
		return json.NewEncoder(w).Encode(models)
	}

	paramFields := strings.Split(values.Get(KeyFields), ",")

	fields := map[string]string{"ID": "id"}
	model := models[0]
	objT := reflect.TypeOf(model).Elem()

	// mapping struct field names to json keys
	for i := 0; i < objT.NumField(); i++ {
		field := objT.Field(i)
		jsonTag := getJsonKeyFromTag(field.Tag.Get("json"))
		if jsonTag != "-" {
			for _, sqlField := range paramFields {
				if sqlField == jsonTag {
					fields[field.Name] = jsonTag
				}
			}
		}
	}

	// build map containing only the returned fields
	slimModels := []map[string]interface{}{}
	for _, m := range models {
		mV := reflect.ValueOf(m).Elem()
		modelMap := map[string]interface{}{}
		for k, v := range fields {
			i := mV.FieldByName(k).Interface()
			modelMap[v] = i
		}
		slimModels = append(slimModels, modelMap)
	}

	return json.NewEncoder(w).Encode(slimModels)
}
