package reporter_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/computerphilosopher/rest-bearer/pkg/gittypes"
	"github.com/computerphilosopher/rest-bearer/pkg/reporter"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(filename)

	reportDir := filepath.Join(pkgDir, "test")
	repoDir := filepath.Join(pkgDir, "../../rest-bearer-test")
	reporter := reporter.NewDefaultReporter(reportDir, repoDir)
	//repo := reporter.Repository

	//read from file

	commit := gittypes.Commit{
		Repository: gittypes.Repository{
			Owner: "tester",
			Name:  "js-repo",
		},
		Id: "6be94cd16db8a555f88290881daece882a4b677e",
	}

	result, err := reporter.Read(commit)
	assert.Nil(err)
	assert.Contains(result, "WARNING")

	//cache hit
	result, err = reporter.Read(commit)
	assert.Nil(err)
	assert.Contains(result, "WARNING")

	err = reporter.Create(commit)
	assert.Nil(err)
	reportPath := filepath.Join(repoDir, commit.Repository.Owner, commit.Repository.Name, commit.Id)
	assert.FileExists(reportPath)

	err = os.Remove(reportPath)
	assert.Nil(err)
}
