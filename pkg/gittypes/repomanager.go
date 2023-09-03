package gittypes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

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

func (manager *singletoneRepoManager) loadMutex(repo Repository) (*sync.RWMutex, error) {
	key := repoPath(repo)
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
	mutex, err := manager.loadMutex(commit.Repository)
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()

	cmd := exec.CommandContext(ctx, "git", "reset", "--hard", commit.Id)
	cmd.Dir = baseDir

	return cmd.Run()
}

func (manager *singletoneRepoManager) Exists(dir string) bool {
	stat, err := os.Stat(dir + ".git")
	if !os.IsExist(err) {
		return false
	}
	return stat.IsDir()
}

func (manager *singletoneRepoManager) Clone(ctx context.Context, baseDir string, remote Remote) error {
	url := remote.ToString()

	cmd := exec.CommandContext(ctx, "git", "clone", url)
	cmd.Dir = baseDir

	return cmd.Run()
}
