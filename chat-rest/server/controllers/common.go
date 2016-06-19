package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/rutmir/services/core/log"
	"github.com/rutmir/services/core/memcache"
	models "github.com/rutmir/services/entities/models/v2"
	"github.com/rutmir/services/web-rest/server/dal"
	"gopkg.in/mgo.v2/bson"
)

const (
	authTypeBasic byte = iota

//authTypeBearer
)

const (
	grantTypePassword byte = iota
	grantTypeRefreshToken

//grantTypeAuthorizationCode
)

type authData struct {
	AuthType byte
	Username string
	Password string
}

func (ad *authData) toString() string {
	return fmt.Sprintf("authData@ AuthType: %v, Username: %s, Password: %s.", ad.AuthType, ad.Username, ad.Password)
}

type grantData struct {
	GrantType    byte
	Username     string
	Password     string
	RefreshToken string
}

func (ad *grantData) toString() string {
	return fmt.Sprintf("grantData@ AuthType: %v, Username: %s, Password: %s, RefreshToken: %s.", ad.GrantType, ad.Username, ad.Password, ad.RefreshToken)
}

var memCtrl memcache.MemCache

func validateClientRequest(r *http.Request) (*authData, error) {
	data, err := parseAuthHeader(r)
	if err != nil {
		return nil, err
	}

	clientCollection := dal.Session.DB("test").C("/v2/clientapps")
	n, err := clientCollection.Find(bson.M{"clientId": data.Username, "clientSecret": data.Password}).Count()
	if err != nil {
		return nil, err
	}

	if n > 0 {
		return data, nil
	}
	return nil, nil
}

func parseAuthHeader(r *http.Request) (*authData, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return nil, fmt.Errorf("Unauthorized")
	}
	h := strings.TrimSpace(authHeader)
	sp := strings.Split(h, " ")

	if len(sp) > 1 {
		result := new(authData)
		switch strings.ToLower(sp[0]) {
		case "basic":
			result.AuthType = authTypeBasic

			str, err := base64.StdEncoding.DecodeString(sp[1])
			if err != nil {
				return nil, fmt.Errorf("Auth credential required.")
			}
			sp = strings.Split(string(str), ":")
			if len(sp) > 1 {
				result.Username = sp[0]
				result.Password = sp[1]
			} else {
				return nil, fmt.Errorf("invalid client id or secret")
			}
			return result, nil
		case "bearer":
			break
		}
	}
	return nil, fmt.Errorf("Auth type not supported.")
}

func parseGrant(r *http.Request) (*grantData, error) {
	grant := strings.TrimSpace(r.FormValue("grant_type"))

	if len(grant) > 1 {
		result := new(grantData)
		switch strings.ToLower(grant) {
		case "password":
			result.GrantType = grantTypePassword
			username := strings.TrimSpace(r.FormValue("username"))
			password := strings.TrimSpace(r.FormValue("password"))
			if len(username) > 0 && len(password) > 0 {
				result.Password = password
				result.Username = username
				return result, nil
			}
			return nil, fmt.Errorf("username and password required")
		case "refresh_token":
			result.GrantType = grantTypeRefreshToken
			refreshToken := strings.TrimSpace(r.FormValue("refresh_token"))
			if len(refreshToken) > 0 {
				result.RefreshToken = refreshToken
				return result, nil
			}
			return nil, fmt.Errorf("refresh_token required")

		}
		return nil, fmt.Errorf("grant_type not supported")
	}
	return nil, fmt.Errorf("grant_type required")
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func writeAuthError(w http.ResponseWriter, error, realm string) {
	w.Header().Add("WWW-Authenticate", realm)
	w.Header().Set("Cache-Control", "no-cache")

	http.Error(w, error, 401)
}

func writeAuthErrorJSON(w http.ResponseWriter, realm string, i interface{}) {
	w.Header().Add("WWW-Authenticate", realm)
	w.Header().Set("Cache-Control", "no-cache")
	buf, err := json.Marshal(i)
	if err != nil {
		http.Error(w, fmt.Sprintf("json.Marshal failed: %v", err), 401)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(buf)
	}
}

func writeErrorJSON(w http.ResponseWriter, code int, i interface{}) {
	w.Header().Set("Cache-Control", "no-cache")
	buf, err := json.Marshal(i)
	if err != nil {
		http.Error(w, fmt.Sprintf("json.Marshal failed: %v", err), code)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(buf)
	}
}

func returnJson(w http.ResponseWriter, obj interface{}) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	buf, err := json.Marshal(obj)
	if err != nil {
		w.Header().Set("Cache-Control", "no-cache")
		http.Error(w, fmt.Sprintf("json.Marshal failed: %v", err), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Cache-Control", "max-age=600, must-revalidate")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	}
}

func checkUserName(username string) (bool, error) {
	collection := dal.Session.DB("test").C("/v2/useraccounts")
	n, err := collection.Find(bson.M{"IdentityProvider": "local", "NameIdentity": username}).Count()
	if err != nil {
		return false, err
	}

	return n == 0, nil
}

func checkExistLocalAccountForProfile(profileID bson.ObjectId) (bool, error) {
	collection := dal.Session.DB("test").C("/v2/useraccounts")
	n, err := collection.Find(bson.M{"IdentityProvider": "local", "profileID": profileID}).Count()
	if err != nil {
		return false, err
	}

	return n == 0, nil
}

func createLocalAccount(newAcc *models.UserAccount) error {
	collection := dal.Session.DB("test").C("/v2/useraccounts")
	newAcc.ID = bson.NewObjectId()
	err := collection.Insert(newAcc)
	return err
}

func createUserProfile(newProf *models.UserProfile) error {
	collection := dal.Session.DB("test").C("/v2/userprofiles")
	newProf.ID = bson.NewObjectId()
	err := collection.Insert(newProf)
	return err
}

func init() {
	var err error
	memCtrl, err = memcache.GetInstance("memcached", "test", "192.168.2.177:11211")
	if err != nil {
		log.Fatal(err)
	}
}
