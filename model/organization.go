package model

type Organization struct {
	OrganizationName string `json:"organizationname,omitempty" bson:"organizationname,omitempty"`
	OrganizationType string `json:"organizationtype,omitempty" bson:"organizationtype,omitempty"`
	OrganizationURL  string `json:"organizationurl,omitempty" bson:"organizationurl,omitempty"`
}
