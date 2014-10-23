package mware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SyntropyDev/httperr"
	"github.com/SyntropyDev/mms-api/model"
	"github.com/SyntropyDev/sqlutil"
	"github.com/lann/squirrel"
)

const (
	authEmailKey = "auth-email"
	authTokenKey = "auth-token"
)

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

func LoginHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {

		type loginReq struct {
			Email    string
			Password string
		}

		req := &loginReq{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		member, err := model.AuthenticateMember(dbmap, req.Email, req.Password)
		if err != nil {
			return err
		}

		token := &model.Token{
			MemberID: member.ID,
		}
		if err := dbmap.Insert(token); err != nil {
			return err
		}

		member.Token = token.Value
		return json.NewEncoder(w).Encode(member)
	}
}

func LogoutHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		tokenValue := r.URL.Query().Get(authTokenKey)
		query := squirrel.Select("*").
			From(model.TableNameToken).
			Where(squirrel.Eq{"Value": tokenValue})

		tokens := []*model.Token{}
		if err := sqlutil.Select(dbmap, query, &tokens); err != nil {
			return err
		}
		for _, token := range tokens {
			dbmap.Delete(token)
		}
		return nil
	}
}

// func RequestInviteHandler() httperr.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) error {

// 		type requestInviteReq struct {
// 			Name     string
// 			Email    string
// 			Password string
// 		}

// 		req := &requestInviteReq{}
// 		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 			return httperr.New(http.StatusBadRequest, err.Error(), err)
// 		}

// 		dbmap, err := getDB()
// 		defer dbmap.Db.Close()
// 		if err != nil {
// 			return err
// 		}

// 		coms := []*model.Community{}
// 		if _, err := dbmap.Select(&coms, "select * from communities"); err != nil {
// 			return err
// 		}
// 		community := coms[0]

// 		if community.RegistrationPolicy == model.RegistrationPolicyOpen {
// 			member := &model.Member{
// 				Email:    req.Email,
// 				Password: req.Password,
// 			}
// 		}
// 		return nil
// 	}
// }

func InviteHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		type inviteReq struct {
			MemberID int64
			Email    string
		}

		req := &inviteReq{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		member := &model.Member{}
		if err := sqlutil.SelectOneRelation(dbmap, model.TableNameMember, req.MemberID, member); err != nil {
			return httperr.New(http.StatusBadRequest, "member not found", err)
		}
		member.SetPassword(model.NewAutoPassword())
		member.Email = req.Email
		if _, err := dbmap.Update(member); err != nil {
			return err
		}
		if err := member.Invite(req.Email); err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(member)
	}
}

func SignupHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		// check if the community already exits, if so prevent signup
		coms := []*model.Community{}
		dbmap.Select(&coms, "select * from communities")
		if len(coms) > 0 {
			err := errors.New("community already created")
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		type signupReq struct {
			Email              string
			Password           string
			Name               string
			CommunityName      string
			RegistrationPolicy string
			Location           []float64
			Description        string
		}

		req := &signupReq{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		trans, err := dbmap.Begin()
		if err != nil {
			return err
		}

		member := &model.Member{
			Email:     req.Email,
			Organizer: true,
			Name:      req.Name,
		}
		pword, err := model.NewPassword(req.Password)
		if err != nil {
			return httperr.New(http.StatusBadRequest, "password must be between 7 and 32 characters", err)
		}
		member.SetPassword(pword)

		if err := trans.Insert(member); err != nil {
			return err
		}

		com := &model.Community{
			Name:               req.CommunityName,
			Location:           req.Location,
			Description:        req.Description,
			RegistrationPolicy: req.RegistrationPolicy,
		}
		if err := trans.Insert(com); err != nil {
			return err
		}

		categories := []*model.Category{
			{Name: "Play Local"},
			{Name: "Be Local"},
			{Name: "Eat Local"},
			{Name: "Shop Local"},
		}
		for _, cat := range categories {
			if err := trans.Insert(cat); err != nil {
				return err
			}
		}
		if err := trans.Commit(); err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(member)
	}
}

func ResetPasswordHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {

		type resetPasswordReq struct {
			Email string
		}

		req := &resetPasswordReq{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		member, err := model.FindMember(dbmap, req.Email)
		// if member not found return 200 anyway
		if err != nil {
			return nil
		}
		if err := member.ResetPassword(); err != nil {
			return err
		}
		if _, err := dbmap.Update(member); err != nil {
			return err
		}
		return nil
	}
}

func ChangePasswordHandler() httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {

		type changePasswordReq struct {
			OldPassword string
			NewPassword string
		}

		req := &changePasswordReq{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		dbmap, err := getDB()
		defer dbmap.Db.Close()
		if err != nil {
			return err
		}

		email := r.URL.Query().Get(authEmailKey)
		member, err := model.FindMember(dbmap, email)
		if err != nil {
			return err
		}

		if !member.HasPassword(req.OldPassword) {
			err := errors.New("password invalid")
			return httperr.New(http.StatusBadRequest, err.Error(), err)
		}

		pword, err := model.NewPassword(req.NewPassword)
		if err != nil {
			return httperr.New(http.StatusBadRequest, "password must be between 7 and 32 characters", err)
		}
		member.SetPassword(pword)
		if _, err := dbmap.Update(member); err != nil {
			return err
		}
		return nil
	}
}
