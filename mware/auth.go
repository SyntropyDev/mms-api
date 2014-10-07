package mware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/mms-api/model"
)

const (
	authEmailKey = "auth-email"
	authTokenKey = "auth-token"
)

type ResetPasswordReq struct {
	Email string
}

func Auth(h httperr.Handler) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		v := r.URL.Query()
		email := v.Get(authEmailKey)
		token := v.Get(authTokenKey)

		errResp := httperr.New(http.StatusUnauthorized, "not authorized", errors.New("not authorized"))
		member, err := model.FindMember(dbmap, email)
		if err != nil {
			return errResp
		} else if err := model.ValidateToken(dbmap, member.ID, token); err != nil {
			return errResp
		}
		return h(w, r)
	}
}

func ResetPasswordHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := &ResetPasswordReq{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		// return 200 no matter what
		member, err := model.FindMember(dbmap, req.Email)
		if err != nil {
			return nil
		}
		member.ResetPassword()
		dbmap.Update(member)
		return nil
	}
}
