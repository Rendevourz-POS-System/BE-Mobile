package entities

type UserToken struct {
	Id     uint   `json:"Id" bson:"_id"`
	UserId uint   `json:"UserId" bson:"UserId"`
	Token  string `json:"Token" bson:"Token"`
	IsUsed bool   `json:"IsUsed" bson:"IsUsed"`
}
