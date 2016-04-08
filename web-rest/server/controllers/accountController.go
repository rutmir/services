package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rutmir/services/core/log"
	models "github.com/rutmir/services/entities/models/v2"
	"github.com/rutmir/services/web-rest/server/responses"
	"gopkg.in/mgo.v2/bson"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	data, err := validateClientRequest(r)
	if err != nil {
		writeAuthError(w, err.Error(), "Basic realm=\"Users\", error:invalid_realm")
		return
	}
	if data == nil {
		writeAuthError(w, "Client application not exist", "Basic realm=\"Users\"")
		return
	}

	ua, err := populateUserAccount(r)
	if err != nil {
		log.Err(err)
		writeErrorJSON(w, 400, responses.Error{Error: "Bad Request", Description: err.Error()})
		return
	}

	ua.CreatedBy = data.Username
	ua.UpdatedBy = data.Username

	ok, err := checkUserName(ua.NameIdentity)
	if err != nil {
		log.Err(err)
		writeErrorJSON(w, 400, responses.Error{Error: "Bad Request", Description: err.Error()})
		return
	}
	if !ok {
		writeErrorJSON(w, 400, responses.Error{Error: "Bad Request", Description: fmt.Sprintf("Username: %s already exist", ua.NameIdentity)})
		return
	}

	if len(ua.ProfileID) > 0 {
		ok, err := checkExistLocalAccountForProfile(ua.ProfileID)
		if err != nil {
			log.Err(err)
			writeErrorJSON(w, 400, responses.Error{Error: "Bad Request", Description: err.Error()})
			return
		}
		if !ok {
			writeErrorJSON(w, 400, responses.Error{Error: "Bad Request", Description: fmt.Sprintf("User profile: %s already have local account", ua.ProfileID)})
			return
		}
	}

	if len(ua.ProfileID) == 0 {
		up := new(models.UserProfile)
		up.DisplayName = ua.NameIdentity
		up.Communications = make([]models.CommunicationItem, 1)
		up.Communications = append(up.Communications, *new(models.CommunicationItem))
		up.Communications[0].Type = "main"
		up.Communications[0].Value = ua.Email
		up.CreatedBy = data.Username
		up.CreatedDate = ua.CreatedDate
		up.UpdatedBy = data.Username
		up.UpdatedDate = ua.UpdatedDate
		if err = createUserProfile(up); err != nil {
			log.Err(err)
			writeErrorJSON(w, 500, responses.Error{Error: "Internal Server Error", Description: err.Error()})
			return
		}
		log.Info("ProfileID %v", up.ID)
		ua.ProfileID = up.ID
	}

	ua.Password = models.HashPassword(ua.Password)
	err = createLocalAccount(ua)
	if err != nil {
		log.Warn(err)
		writeErrorJSON(w, 500, responses.Error{Error: "Internal Server Error", Description: err.Error()})
		return
	}

	returnJson(w, ua)
}

func populateUserAccount(r *http.Request) (*models.UserAccount, error) {
	now := time.Now()
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	email := strings.TrimSpace(r.FormValue("email"))
	profileIDRaw := strings.TrimSpace(r.FormValue("profileID"))

	var profileID bson.ObjectId

	if len(profileIDRaw) > 0 {
		profileID = bson.ObjectIdHex(profileIDRaw)
	}

	ua := new(models.UserAccount)

	ua.IdentityProvider = "local"
	ua.Email = email
	ua.NameIdentity = strings.ToLower(username)
	ua.Password = password
	ua.ProfileID = profileID
	ua.LoginAttempts = 0
	ua.IsLockedOut = false
	ua.IsApproved = false
	ua.CreatedDate = now
	ua.UpdatedDate = now

	if len(ua.NameIdentity) < 3 {
		return nil, fmt.Errorf("invalid username")
	}

	if len(ua.Password) < 3 {
		return nil, fmt.Errorf("invalid password")
	}

	return ua, nil
}
