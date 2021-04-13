package utils

import "sync"

type CommonMutex struct {
	OrganizationMutex map[string]*OrganizationMutex //mutex per ogni organization
	GeneralMutex      *sync.Mutex
}

type OrganizationMutex struct {
	Mutex *sync.Mutex
	Count uint //numero di thread in attesa su OrganizationMutex dell'evento
}

func NewEventMutex() CommonMutex {
	return CommonMutex{OrganizationMutex: make(map[string]*OrganizationMutex), GeneralMutex: &sync.Mutex{}}
}

func ReserveOrganizationMutex(organizationName string, mutex *CommonMutex) *sync.Mutex {
	var eventMutex *sync.Mutex

	mutex.GeneralMutex.Lock()

	if _, ok := mutex.OrganizationMutex[organizationName]; !ok {
		mutex.OrganizationMutex[organizationName] = &OrganizationMutex{Mutex: &sync.Mutex{}, Count: 1}
	} else {
		mutex.OrganizationMutex[organizationName].Count++
	}

	eventMutex = mutex.OrganizationMutex[organizationName].Mutex

	mutex.GeneralMutex.Unlock()

	return eventMutex
}

func ReleaseOrganizationMutex(organizationName string, mutex *CommonMutex) {
	mutex.GeneralMutex.Lock()

	if _, ok := mutex.OrganizationMutex[organizationName]; ok {
		mutex.OrganizationMutex[organizationName].Count--

		if mutex.OrganizationMutex[organizationName].Count == 0 {
			delete(mutex.OrganizationMutex, organizationName)
		}
	}

	mutex.GeneralMutex.Unlock()
}

func ReleaseOrganizationMutexDefer(uid string, commonMutex *CommonMutex, mutex *sync.Mutex, locked *bool) {
	if *locked {
		mutex.Unlock()
		ReleaseOrganizationMutex(uid, commonMutex)
	}
}
