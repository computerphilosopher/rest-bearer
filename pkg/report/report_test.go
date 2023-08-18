package report_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/computerphilosopher/rest-bearer/pkg/report"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(filename)

	repoDir := filepath.Join(pkgDir, "test")
	reporter := report.NewDefaultReporter(repoDir)
	//repo := reporter.Repository

	//read from file
	result, err := reporter.Read("tester", "js-repo", "6be94cd16db8a555f88290881daece882a4b677e")
	assert.Nil(err)
	assert.Contains(result, "WARNING")

	//cache hit
	result, err = reporter.Read("tester", "js-repo", "6be94cd16db8a555f88290881daece882a4b677e")
	assert.Nil(err)
	assert.Contains(result, "WARNING")

	err = reporter.Write("tester", "js-repo", "abcde", "report")
	assert.Nil(err)
	reportPath := filepath.Join(repoDir, "tester", "js-repo", "abcde")
	assert.FileExists(reportPath)

	err = os.Remove(reportPath)
	assert.Nil(err)
}
