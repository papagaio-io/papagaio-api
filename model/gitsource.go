package model

import "wecode.sorint.it/opensource/papagaio-api/types"

type GitSource struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	GitType     types.GitType `json:"gitType"`
	GitAPIURL   string        `json:"gitApiUrl"`
	GitClientID string        `json:"gitClientId"`
	GitSecret   string        `json:"gitSecret"`

	//Campi creati da app
	AgolaRemoteSource string `json:"agolaRemoteSource"`
	//AgolaClientID     string `json:"agolaClientId"`
	//AgolaSecret       string `json:"agolaSecret"`

	//Note gitsource: il remotesource viene creato, in automatico, quando creo il gitsource, utilizzando l'account dell'utente loggato.
	//Ad Agola vengono poi passati clienteID+secret+remotesource name

	//Note: nella organizaion model lascio la mail dell'utente che ha creato la organization(è owner!), creo il token e lo setto sulla organization
	//L'utente che crea la organization deve essere registrato in Agola con quel remotesource
	//Gestire nelle chiamate verso agola la rigenerazione del token e il salvataggio nella organization model
}
