package gittypes

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type RepoManager interface {
	Reset(ctx context.Context, commit Commit) error
	Exists(Repository) bool
}

type repomanagerWrapper struct {
	baseDir string
	manager *singletoneRepoManager
}
type singletoneRepoManager struct {
	locks *MapLock
}

var instance *singletoneRepoManager = nil
var once sync.Once

func GetRepoManager(baseDir string) repomanagerWrapper {
	once.Do(func() {
		instance = &singletoneRepoManager{
			locks: &MapLock{},
		}
	})
	return repomanagerWrapper{
		baseDir: baseDir,
		manager: instance,
	}
}

func (wrapper repomanagerWrapper) Reset(ctx context.Context, commit Commit) error {
	return wrapper.manager.reset(ctx, wrapper.baseDir, commit)
}

func (wrapper repomanagerWrapper) Exists(repo Repository) bool {
	repoDir := filepath.Join(wrapper.baseDir, repo.ToString())
	return wrapper.manager.exists(repoDir)
}

func (wrapper repomanagerWrapper) Clone(ctx context.Context, remote Remote) error {
	return wrapper.manager.clone(ctx, wrapper.baseDir, remote)
}

func (manager *singletoneRepoManager) reset(ctx context.Context, baseDir string, commit Commit) error {
	mutex, err := manager.locks.Load(commit.Repository)
	if err != nil {
		return err
	}
	mutex.Lock()
	defer mutex.Unlock()

	cmd := exec.CommandContext(ctx, "git", "reset", "--hard", commit.Id)
	cmd.Dir = baseDir

	return cmd.Run()
}

func (manager *singletoneRepoManager) exists(dir string) bool {
	stat, err := os.Stat(dir + ".git")
	if !os.IsExist(err) {
		return false
	}
	return stat.IsDir()
}

func (manager *singletoneRepoManager) clone(ctx context.Context, baseDir string, remote Remote) error {
	url := remote.ToString()

	cmd := exec.CommandContext(ctx, "git", "clone", url)
	cmd.Dir = baseDir

	mutex, err := manager.locks.Load(remote.Repository)
	if err != nil {
		return err
	}
	mutex.Lock()
	defer mutex.Unlock()

	err = cmd.Run()
	return err
}
