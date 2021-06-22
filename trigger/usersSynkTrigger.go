package trigger

import (
	"errors"
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func StartSynkUsers(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	synkUsersRun(db, tr, commonMutex, agolaApi, gitGateway)
}

func synkUsersRun(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	for {
		log.Println("start users synk")

		usersVerifiedOK := make(map[uint64]bool)
		usersVerifiedError := make(map[uint64]bool)

		organizationsRef, _ := db.GetOrganizationsRef()
		for _, organizationRef := range organizationsRef {
			mutex := utils.ReserveOrganizationMutex(organizationRef, commonMutex)
			mutex.Lock()

			locked := true
			defer utils.ReleaseOrganizationMutexDefer(organizationRef, commonMutex, mutex, &locked)

			org, _ := db.GetOrganizationByAgolaRef(organizationRef)
			if org == nil {
				log.Println("synkUsersRun organization ", organizationRef, "not found")
				continue
			}

			userError := false
			if usersVerifiedOK[org.UserIDConnected] {
				continue
			}

			gitSource, err := db.GetGitSourceByName(org.GitSourceName)
			if err != nil || gitSource == nil {
				log.Println("gitSource", org.GitSourceName, "not found:", err)
				continue
			}

			if usersVerifiedError[org.UserIDConnected] {
				userError = true
			} else {
				user, err := verifyUserAccount(org.UserIDConnected, db, gitGateway, agolaApi, gitSource, org.Name)
				if user != nil {
					db.SaveUser(user)
				}

				if err != nil {
					userError = true
					usersVerifiedError[org.UserIDConnected] = true
					log.Println("userError:", err)
				}
			}

			if userError {
				user := findUserToConnect(db, gitGateway, agolaApi, gitSource, org.Name)
				if user != nil {
					org.UserIDConnected = *user.UserID
					db.SaveOrganization(org)
				}
			}

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
			locked = false
		}

		//TODO trigger time value
		//time.Sleep(time.Duration(tr.GetOrganizationsTriggerTime()) * time.Minute)
		time.Sleep(24 * time.Hour)
	}
}

func findUserToConnect(db repository.Database, gitGateway *git.GitGateway, agolaApi agola.AgolaApiInterface, gitSource *model.GitSource, gitOrgRef string) *model.User {
	usersID, _ := db.GetUsersIDByGitSourceName(gitSource.Name)
	for _, id := range usersID {
		user, err := verifyUserAccount(id, db, gitGateway, agolaApi, gitSource, gitOrgRef)
		if user != nil {
			db.SaveUser(user)
		}

		if err == nil {
			return user
		}
	}

	return nil
}

func verifyUserAccount(userID uint64, db repository.Database, gitGateway *git.GitGateway, agolaApi agola.AgolaApiInterface, gitSource *model.GitSource, gitOrgRef string) (*model.User, error) {
	user, err := db.GetUserByUserId(userID)
	if err != nil || user == nil {
		log.Println("user not found")
		return nil, errors.New("user not found in db")
	}

	err = verifyUserGiteaAccount(user, gitGateway, gitSource, gitOrgRef)
	if err != nil {
		return user, err
	}

	err = verifyUserAgolaAccount(user, agolaApi, gitSource)
	return user, err
}

func verifyUserGiteaAccount(user *model.User, gitGateway *git.GitGateway, gitSource *model.GitSource, gitOrgRef string) error {
	//TODO gestire utente cancellato da git

	//TODO gestire la modifica dei campi info utente(es. isAdmin)

	if common.IsAccessTokenExpired(user.Oauth2AccessTokenExpiresAt) {
		token, err := gitGateway.RefreshToken(gitSource, user.Oauth2RefreshToken)
		if err != nil {
			log.Println("error in RefreshToken:", err)
			return err
		}

		user.Oauth2AccessToken = token.AccessToken
		user.Oauth2RefreshToken = token.RefreshToken
		user.Oauth2AccessTokenExpiresAt = time.Now().Add(time.Second * time.Duration(token.Expiry))
	}

	isOwner, err := gitGateway.IsUserOwner(gitSource, user, gitOrgRef)
	if err != nil {
		return err
	}

	if !isOwner {
		return errors.New("user is not owner")
	}

	return nil
}

func verifyUserAgolaAccount(user *model.User, agolaApi agola.AgolaApiInterface, gitSource *model.GitSource) error {
	agolaUserRef := utils.GetAgolaUserRefByGitUsername(agolaApi, gitSource.AgolaRemoteSource, user.Login)
	if agolaUserRef == nil {
		user.AgolaUserRef = nil
		return errors.New("user not present in Agola")
	}
	user.AgolaUserRef = agolaUserRef

	return nil
}
