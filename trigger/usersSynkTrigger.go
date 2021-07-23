package trigger

import (
	"errors"
	"fmt"
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/trigger/dto"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func StartSynkUsers(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway, c chan string, rtDto *dto.TriggersRunTimeDto) {
	go synkUsersRun(db, tr, commonMutex, agolaApi, gitGateway, c, rtDto)
	go userRunTimer(tr, c)
}

func userRunTimer(tr utils.ConfigUtils, c chan string) {
	for {
		log.Println("userRunTimer wait for", time.Duration(time.Minute.Nanoseconds()*int64(tr.GetUsersTriggerTime())))
		time.Sleep(time.Duration(time.Minute.Nanoseconds() * int64(tr.GetUsersTriggerTime())))
		c <- "resume"
	}
}

func synkUsersRun(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway, c chan string, rtDto *dto.TriggersRunTimeDto) {
	for {
		log.Println("start users synk")
		rtDto.UsersTriggerLastRun = time.Now()

		usersVerifiedOK := make(map[uint64]*model.User)

		usersId, err := db.GetUsersID()
		if err != nil {
			log.Println("error in GetUsersID:", err)
			continue
		}

		for _, userId := range usersId {
			user, err := db.GetUserByUserId(userId)
			if err != nil || user == nil {
				log.Println("user not found in db")
				continue
			}

			log.Println("synk user:", userId, user.Login, user.ID)

			gitSource, err := db.GetGitSourceByName(user.GitSourceName)
			if err != nil || gitSource == nil {
				log.Println("gitSource", user.GitSourceName, "not found in db:", err)
				continue
			}
			userGitNotFound, err := verifyUserAccount(user, db, gitGateway, agolaApi, gitSource)

			if userGitNotFound {
				log.Println("delete user:", user.UserID)
				err := db.DeleteUser(*user.UserID)

				if err != nil {
					log.Println("error in DeleteUser:", err)
				}
			} else if user != nil {
				_, err = db.SaveUser(user)

				if err != nil {
					log.Println("error in SaveUser:", err)
				}
			}

			if err == nil {
				usersVerifiedOK[userId] = user
			}
		}

		organizationsRef, _ := db.GetOrganizationsRef()
		for _, organizationRef := range organizationsRef {
			mutex := utils.ReserveOrganizationMutex(organizationRef, commonMutex)
			mutex.Lock()

			org, _ := db.GetOrganizationByAgolaRef(organizationRef)
			if org == nil {
				log.Println("synkUsersRun organization ", organizationRef, "not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			user := usersVerifiedOK[org.UserIDConnected]

			gitSource, _ := db.GetGitSourceByName(org.GitSourceName)
			if err != nil || gitSource == nil {
				log.Println("gitSource", org.GitSourceName, "not found:", err)

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			if user != nil {
				isOwner, _ := gitGateway.IsUserOwner(gitSource, user, org.GitPath)
				if isOwner {
					mutex.Unlock()
					utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

					continue
				}
			}

			log.Println("findUserToConnect for organization ", org.GitPath)
			user = findUserToConnect(db, gitGateway, agolaApi, gitSource, org.GitPath, usersVerifiedOK)
			if user != nil {
				log.Println("findUserToConnect result UserID", user.UserID)
				org.UserIDConnected = *user.UserID
				err = db.SaveOrganization(org)

				if err != nil {
					log.Println("error in SaveOrganization:", err)
				}
			}

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
		}

		fmt.Println("synkUsersRun:", <-c)
	}
}

func findUserToConnect(db repository.Database, gitGateway *git.GitGateway, agolaApi agola.AgolaApiInterface, gitSource *model.GitSource, gitOrgRef string, usersVerifiedOK map[uint64]*model.User) *model.User {
	for _, user := range usersVerifiedOK {
		isOwner, _ := gitGateway.IsUserOwner(gitSource, user, gitOrgRef)
		if isOwner {
			return user
		}
	}

	return nil
}

func verifyUserAccount(user *model.User, db repository.Database, gitGateway *git.GitGateway, agolaApi agola.AgolaApiInterface, gitSource *model.GitSource) (bool, error) {
	userGitNotFound, err := verifyUserGiteaAccount(user, gitGateway, gitSource)

	if userGitNotFound || err != nil {
		return userGitNotFound, err
	}

	err = verifyUserAgolaAccount(user, agolaApi, gitSource)
	return false, err
}

func verifyUserGiteaAccount(user *model.User, gitGateway *git.GitGateway, gitSource *model.GitSource) (bool, error) {
	userInfo, err := gitGateway.GetUserByLogin(gitSource, user)

	if userInfo == nil && err == nil {
		return true, errors.New("user not found in git")
	}

	if err != nil {
		return false, errors.New("error in GetUserByLogin or user not found")
	}

	user.Email = userInfo.Email
	user.IsAdmin = userInfo.IsAdmin
	user.Login = userInfo.Login
	user.ID = uint64(userInfo.ID)

	if common.IsAccessTokenExpired(user.Oauth2AccessTokenExpiresAt) {
		token, err := gitGateway.RefreshToken(gitSource, user.Oauth2RefreshToken)
		if err != nil {
			log.Println("error in RefreshToken:", err)
			return false, err
		}

		user.Oauth2AccessToken = token.AccessToken
		user.Oauth2RefreshToken = token.RefreshToken
		user.Oauth2AccessTokenExpiresAt = time.Now().Add(time.Second * time.Duration(token.Expiry))
	}

	return false, nil
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
