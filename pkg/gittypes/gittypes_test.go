package gittypes_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/computerphilosopher/rest-bearer/pkg/gittypes"
	"github.com/stretchr/testify/assert"
)

const testHead = "03bc415c4b3187fca6adb1268e6f8e282febf86f"

func TestReset(t *testing.T) {

	assert := assert.New(t)

	manager := gittypes.GetRepoManager()
	assert.NotNil(manager)

	_, filename, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(filename)
	repoDir := filepath.Join(pkgDir, "../../rest-bearer-test")

	err := manager.Reset(
		context.TODO(),
		repoDir,
		gittypes.Commit{
			Repository: gittypes.Repository{
				Owner: "computerphilosopher",
				Name:  "rest-bearer-test",
			},
			Id: "origin/test",
		},
	)

	parseCtx, parseCancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer parseCancel()

	cmdParseHead := exec.CommandContext(parseCtx, "git", "rev-parse", "HEAD")
	cmdParseHead.Dir = repoDir
	outputRaw, err := cmdParseHead.CombinedOutput()
	assert.Nil(err)

	output := strings.TrimSpace(string(outputRaw))
	assert.Equal(testHead, output)

	resetCtx, resetCancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer resetCancel()
	cmdResetMaster := exec.CommandContext(resetCtx, "git", "reset", "--hard", "origin/main")
	cmdResetMaster.Dir = repoDir

	outputRaw, err = cmdResetMaster.CombinedOutput()
	assert.Nil(err)
	// t.Fatal("not implemented")
}

func TestClone(t *testing.T) {

	assert := assert.New(t)

	manager := gittypes.GetRepoManager()
	assert.NotNil(manager)

	_, filename, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(filename)
	repoDir := filepath.Join(pkgDir, "rest-bearer-test")

	assert.False(manager.Exists(repoDir))

	remote := gittypes.Remote{
		BaseUrl: "https://github.com",
		Repository: gittypes.Repository{
			Owner: "computerphilosopher",
			Name:  "rest-bearer-test",
		},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	err := manager.Clone(ctx, pkgDir, remote)

	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(err)

	err = os.RemoveAll(repoDir)
	assert.Nil(err)

}
