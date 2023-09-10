package reporter

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/computerphilosopher/rest-bearer/pkg/gittypes"
	"github.com/computerphilosopher/rest-bearer/pkg/maplock"
	cache "github.com/patrickmn/go-cache"
)

type Reporter interface {
	Read(commit gittypes.Commit) (string, error)
	Path(commit gittypes.Commit) string
	Create(ctx context.Context, commit gittypes.Commit) error
}

var reportCache *cache.Cache
var once sync.Once

// TODO: dbReporter
type fileReporter struct {
	reportDir string
	gitDir    string
	lock      *maplock.MapLock
}

func NewDefaultReporter(reportDir, repositoryDir string) Reporter {
	once.Do(func() {
		reportCache = cache.New(5*time.Minute, 10*time.Minute)
	})
	return fileReporter{
		reportDir: reportDir,
		gitDir:    repositoryDir,
		lock:      maplock.GetMapLock(),
	}
}

func cacheKey(commit gittypes.Commit) string {
	return fmt.Sprintf("%s/%s/%s",
		commit.Repository.Owner,
		commit.Repository.Name,
		commit.Id,
	)
}

func (reporter fileReporter) Path(commit gittypes.Commit) string {
	return filepath.Join(reporter.reportDir,
		commit.Repository.Owner,
		commit.Repository.Name,
		commit.Id,
	)
}

func (reporter fileReporter) readFromFile(commit gittypes.Commit) (string, error) {
	reportPath := reporter.Path(commit)
	bytes, err := os.ReadFile(reportPath)
	if err != nil {
		return "",
			fmt.Errorf("can't read report for %s/%s:%s %w",
				commit.Repository.Owner,
				commit.Repository.Name,
				commit.Id,
				err,
			)
	}
	return string(bytes), nil
}

func validateCacheData(raw interface{}) (string, error) {
	cached, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("cached data '%v' is not string type", cached)
	}
	return cached, nil
}

func (reporter fileReporter) Read(commit gittypes.Commit) (string, error) {
	key := cacheKey(commit)
	cached, found := reportCache.Get(key)
	if found {
		return validateCacheData(cached)
	}

	report, err := reporter.readFromFile(commit)
	if err != nil {
		return "", err
	}
	reportCache.Set(key, report, cache.DefaultExpiration)

	return report, err
}

func (reporter fileReporter) Create(ctx context.Context, commit gittypes.Commit) error {
	//return os.WriteFile(reporter.Path(commit), []byte(report), 0644)
	cmd := exec.CommandContext(ctx, "bearer", "scan", reporter.gitDir)
	cmd.Dir = reporter.gitDir

	output, err := func() ([]byte, error) {
		output, err := cmd.CombinedOutput()
		if err == nil {
			return output, nil
		}

		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return nil, err
		}
		if exitErr.ExitCode() != 1 {
			return nil, err
		}

		return output, nil
	}()
	if err != nil {
		return err
	}

	path := reporter.Path(commit)
	return os.WriteFile(path, output, fs.FileMode(os.O_RDWR))
}
