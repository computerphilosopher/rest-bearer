package reporter

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type Reporter interface {
	Read(owner, repoName, commitHash string) (string, error)
	Write(owner, repoName, commitHash, report string) error
}

var reportCache *cache.Cache
var once sync.Once

// TODO: DBReporter
type FileReporter struct {
	BaseDir string
}

func NewDefaultReporter(baseDir string) Reporter {
	once.Do(func() {
		reportCache = cache.New(5*time.Minute, 10*time.Minute)
	})
	return FileReporter{
		BaseDir: baseDir,
	}
}

func cacheKey(owner, repoName, commitHash string) string {
	return fmt.Sprintf("%s/%s/%s", owner, repoName, commitHash)
}

func (reporter FileReporter) readFromFile(owner, repoName, commitHash string) (string, error) {
	reportPath := filepath.Join(reporter.BaseDir, owner, repoName, commitHash)
	bytes, err := os.ReadFile(reportPath)
	if err != nil {
		return "",
			fmt.Errorf("can't read report for %s/%s:%s %w",
				owner, repoName, commitHash, err)
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

func (reporter FileReporter) Read(owner, repoName, commitHash string) (string, error) {
	key := cacheKey(owner, repoName, commitHash)
	cached, found := reportCache.Get(key)
	if found {
		return validateCacheData(cached)
	}

	report, err := reporter.readFromFile(owner, repoName, commitHash)
	if err != nil {
		return "", err
	}
	reportCache.Set(key, report, cache.DefaultExpiration)

	return report, err
}

func (reporter FileReporter) Write(owner,
	repoName, commitHash string, report string) error {

	reportPath := filepath.Join(reporter.BaseDir, owner, repoName, commitHash)
	return os.WriteFile(reportPath, []byte(report), 0644)
}
