package repository

import "wecode.sorint.it/opensource/papagaio-be/model"

//TODO
func (db *AppDb) GetGitSources() (*[]model.GitSource, error) {
	var retVal *[]model.GitSource
	var err error
	return retVal, err
}

//TODO
func (db *AppDb) SaveGitSource(gitSource *model.GitSource) error {
	var err error
	return err
}

//TODO
func (db *AppDb) GetGitSourceByID(id string) (*model.GitSource, error) {
	var retVal *model.GitSource
	var err error
	return retVal, err
}

//TODO
func (db *AppDb) DeleteGitSource(id string) error {
	var err error
	return err
}
