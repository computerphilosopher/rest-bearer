package reporter

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/computerphilosopher/rest-bearer/pkg/gittypes"
	"github.com/computerphilosopher/rest-bearer/pkg/maplock"
	cache "github.com/patrickmn/go-cache"
)

type Reporter interface {
	Read(commit gittypes.Commit) (string, error)
	Write(commit gittypes.Commit, report string) error
	Path(commit gittypes.Commit) string
}

var reportCache *cache.Cache
var once sync.Once

// TODO: DBReporter
type FileReporter struct {
	BaseDir string
	lock    *maplock.MapLock
}

func NewDefaultReporter(baseDir string) Reporter {
	once.Do(func() {
		reportCache = cache.New(5*time.Minute, 10*time.Minute)
	})
	return FileReporter{
		BaseDir: baseDir,
		lock:    maplock.GetMapLock(),
	}
}

func cacheKey(commit gittypes.Commit) string {
	return fmt.Sprintf("%s/%s/%s",
		commit.Repository.Owner,
		commit.Repository.Name,
		commit.Id,
	)
}

func (reporter FileReporter) Path(commit gittypes.Commit) string {
	return filepath.Join(reporter.BaseDir,
		commit.Repository.Owner,
		commit.Repository.Name,
		commit.Id,
	)
}

func (reporter FileReporter) readFromFile(commit gittypes.Commit) (string, error) {
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

func (reporter FileReporter) Read(commit gittypes.Commit) (string, error) {
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

func (reporter FileReporter) Write(commit gittypes.Commit, report string) error {
	return os.WriteFile(reporter.Path(commit), []byte(report), 0644)
}
