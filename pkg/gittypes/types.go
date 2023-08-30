package gittypes

import (
	"fmt"
	"sync"
)

type Repository struct {
	Owner string
	Name  string
}

type Commit struct {
	Repository Repository
	Id         string
}

type RepoManager interface {
	Reset(commit Commit) error
}
type singletoneRepoManager struct {
	mutexes sync.Map
}

var instance *singletoneRepoManager = nil
var once sync.Once

func GetRepoManager() *singletoneRepoManager {
	once.Do(func() {
		instance = &singletoneRepoManager{
			mutexes: sync.Map{},
		}
	})
	return instance
}

func repoPath(repo Repository) string {
	return repo.Owner + "/" + repo.Name
}

func (manager *singletoneRepoManager) Reset(baseDir string, commit Commit) error {
	key := repoPath(commit.Repository)
	mutexRaw, exist := manager.mutexes.Load(key)
	if !exist {
		newLock := &sync.RWMutex{}
		manager.mutexes.Store(key, newLock)
		mutexRaw = newLock
	}

	mutex, ok := mutexRaw.(*sync.RWMutex)
	if !ok {
		return fmt.Errorf("%s's lock type is wrong expected: mutex actual: %T",
			key, mutexRaw)
	}
	mutex.Lock()
	defer mutex.Unlock()

	return nil
}
