package gittypes

import (
	"context"
	"fmt"
	"os/exec"
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
	Reset(ctx context.Context, baseDir string, commit Commit) error
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

func (manager *singletoneRepoManager) loadMutex(key string) (*sync.RWMutex, error) {
	raw, exist := manager.mutexes.Load(key)
	if !exist {
		mutex := &sync.RWMutex{}
		manager.mutexes.Store(key, mutex)
		return mutex, nil
	}
	mutex, ok := raw.(*sync.RWMutex)
	if !ok {
		return nil,
			fmt.Errorf("%s's lock type is wrong expected: mutex actual: %T",
				key, raw)
	}
	manager.mutexes.Store(key, mutex)

	return mutex, nil
}

func (manager *singletoneRepoManager) Reset(ctx context.Context, baseDir string, commit Commit) error {
	key := repoPath(commit.Repository)

	mutex, err := manager.loadMutex(key)
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()

	cmd := exec.CommandContext(ctx, "git", "reset", "--hard", commit.Id)
	cmd.Dir = baseDir

	return cmd.Run()
}
