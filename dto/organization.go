package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrganizationName string             `json:"organizationname,omitempty" bson:"organizationname,omitempty"`
	OrganizationType string             `json:"organizationtype,omitempty" bson:"organizationtype,omitempty"`
	OrganizationURL  string             `json:"organizationurl,omitempty" bson:"organizationurl,omitempty"`
}
