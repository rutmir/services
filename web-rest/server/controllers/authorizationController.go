package controllers

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/golang/protobuf/proto"
	"github.com/rutmir/services/core/log"
	"github.com/rutmir/services/core/memcache"
	dto "github.com/rutmir/services/entities/dto/v2"
	models "github.com/rutmir/services/entities/models/v2"
	"github.com/rutmir/services/web-rest/server/dal"
	"github.com/rutmir/services/web-rest/server/responses"
)

const (
	atPrefix = "at_"
	rtPrefix = "rt_"
)

// Authentication http handler for Authentication request
func Authentication(w http.ResponseWriter, r *http.Request) {
	data, err := validateClientRequest(r)
	if err != nil {
		writeAuthError(w, err.Error(), "Basic realm=\"Users\", error:invalid_realm")
		return
	}

	grant, err := parseGrant(r)
	if err != nil {
		writeAuthErrorJSON(w, "Basic realm=\"Users\", error:invalid_grant", responses.Error{Error: "invalid_grant", Description: err.Error()})
		return
	}

	switch grant.GrantType {
	case grantTypeRefreshToken:
		rtItem, err := memCtrl.Get(rtPrefix + grant.RefreshToken)
		if err != nil {
			log.Warn(err)
			writeAuthErrorJSON(w, "Basic realm=\"Users\"", responses.Error{Error: "Error find user account", Description: err.Error()})
			return
		}

		mem := new(dto.AuthTokenMem)
		if err = proto.Unmarshal(rtItem.Value, mem); err != nil {
			log.Warn(err)
			writeErrorJSON(w, 400, responses.Error{Error: "Bad Request", Description: err.Error()})
			return
		}

		if mem.ClientID != data.Username {
			log.Err("Not valid ClientID on refresh_token exchange")
			writeErrorJSON(w, 400, responses.Error{Error: "Bad Request", Description: "Not valid ClientID on refresh_tocen exchange"})
			return
		}

		clearMem(mem.AccessToken, mem.RefreshToken)
		finalizeAuthentication(w, mem.AccountID, mem.ProfileID, mem.ClientID)
		return
	case grantTypePassword:
		account := models.UserAccount{}
		collection := dal.Session.DB("test").C("/v2/useraccounts")

		err = collection.Find(bson.M{"identityProvider": "local", "nameIdentity": grant.Username, "password": models.HashPassword(grant.Password)}).One(&account)
		if err != nil {
			log.Warn(err)
			writeAuthErrorJSON(w, "Basic realm=\"Users\"", responses.Error{Error: "Error find user account", Description: err.Error()})
			return
		}

		finalizeAuthentication(w, account.ID, account.ProfileID, data.Username)
		return
	}

	writeErrorJSON(w, 501, responses.Error{Error: "Not Implemented", Description: "Required grant type not inplemented"})
}

func finalizeAuthentication(w http.ResponseWriter, accID, profID bson.ObjectId, clientID string) {
	accessToken, err := generateToken()
	if err != nil {
		log.Warn(err)
		writeErrorJSON(w, 500, responses.Error{Error: "Internal server error", Description: err.Error()})
		return
	}

	refreshToken, err := generateToken()
	if err != nil {
		log.Warn(err)
		writeErrorJSON(w, 500, responses.Error{Error: "Internal server error", Description: err.Error()})
		return
	}

	auth := new(responses.Authentication)
	auth.AccessToken = accessToken
	auth.RefreshToken = refreshToken
	auth.TokenType = "Bearer"
	auth.ExpiresIn = 3600
	auth.AccountID = accID
	auth.ProfileID = profID

	_, err = populateMem(accID, profID, clientID, accessToken, refreshToken)
	if err != nil {
		log.Err("set to memcache: ", err)
		writeErrorJSON(w, 500, responses.Error{Error: "Internal server error", Description: err.Error()})
		return
	}

	returnJson(w, auth)
}

func clearMem(accessToken, refreshToken string) (bool, error) {
	if err := memCtrl.Delete(atPrefix + accessToken); err != nil {
		return false, err
	}

	if err := memCtrl.Delete(rtPrefix + refreshToken); err != nil {
		return false, err
	}

	return true, nil
}

func populateMem(accID, profID bson.ObjectId, clientID, accessToken, refreshToken string) (bool, error) {
	mem := new(dto.AuthTokenMem)
	mem.ClientID = clientID
	mem.ProfileID = profID
	mem.AccountID = accID
	mem.AccessToken = accessToken
	mem.RefreshToken = refreshToken

	dataMem, err := proto.Marshal(mem)
	if err != nil {
		return false, err
	}

	dataPkg := &memcache.Item{
		Key:        atPrefix + accessToken,
		Value:      dataMem,
		Expiration: 3600}
	if err := memCtrl.Set(dataPkg); err != nil {
		return false, err
	}

	dataPkg.Key = rtPrefix + refreshToken
	dataPkg.Expiration = 3600 * 24 * 7

	if err := memCtrl.Set(dataPkg); err != nil {
		return false, err
	}

	return true, nil
}
