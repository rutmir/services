package responses

import "gopkg.in/mgo.v2/bson"

type Authentication struct {
	TokenType    string        `json:"token_type,required"`
	AccessToken  string        `json:"access_token,required"`
	RefreshToken string        `json:"refresh_token,required"`
	ExpiresIn    int32         `json:"expires_in,required"`
	AccountID    bson.ObjectId `json:"accountID,required"`
	ProfileID    bson.ObjectId `json:"profileID,required"`
}
