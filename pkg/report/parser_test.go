package report_test

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/computerphilosopher/rest-bearer/pkg/report"
	"github.com/stretchr/testify/assert"
)

func TestMetaInfoParser(t *testing.T) {
	assert := assert.New(t)

	line := "MEDIUM: Insufficiently random value detected. [CWE-330]"
	metaInfo, err := report.ParseMetaInfo(line)

	assert.Nil(err)
	assert.Equal("MEDIUM", metaInfo.Severity)
	assert.Equal("CWE-330", metaInfo.Code)
	assert.Equal("Insufficiently random value detected.", metaInfo.Description)
}

func TestLocationParser(t *testing.T) {
	assert := assert.New(t)

	line := "File: src/main/java/net/spy/memcached/ArcusClientPool.java:76"
	location, err := report.ParseLocation(line)

	assert.Nil(err)
	assert.Equal("src/main/java/net/spy/memcached/ArcusClientPool.java", location.Path)
	assert.Equal(76, location.Line)
}

func TestSnippetParser(t *testing.T) {
	assert := assert.New(t)

	line := "137     md5.update(KeyUtil.getKeyBytes(k));"
	snippet, err := report.ParseSnippet(line)

	assert.Nil(err)
	assert.Equal(snippet, "md5.update(KeyUtil.getKeyBytes(k));")
}

func TestVulnerabilityParser(t *testing.T) {

	assert := assert.New(t)
	lines := []string{
		"MEDIUM: Insufficiently random value detected. [CWE-330]",
		"https://docs.bearer.com/reference/rules/java_lang_insufficiently_random_values",
		"To exclude this finding, use the flag --exclude-fingerprint=05ba7d67a612d13d1c75f296520b8e6f_0",
		"File: src/main/java/net/spy/memcached/ArcusClientPool.java:76",
		" 76     return client[rand.nextInt(poolSize)];",
	}

	expected := report.Vulnerability{
		MetaInfo: report.MetaInfo{
			Severity:    "MEDIUM",
			Code:        "CWE-330",
			Description: "Insufficiently random value detected.",
		},
		Location: report.Location{
			Path: "src/main/java/net/spy/memcached/ArcusClientPool.java",
			Line: 76,
		},
		Reference: "https://docs.bearer.com/reference/rules/java_lang_insufficiently_random_values",
		Snippet:   "return client[rand.nextInt(poolSize)];",
	}

	actual, next, err := report.ParseVulnerability(lines, 0)
	assert.Nil(err)
	assert.Equal(expected, actual)
	assert.Equal(next, 5)

}

func fileToString(relativePath string) (string, error) {

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	testFilePath := filepath.Join(dir, relativePath)

	bytes, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TestParseReport(t *testing.T) {
	assert := assert.New(t)

	content, err := fileToString("test.txt")
	assert.Nil(err)

	rpt, err := report.ParseReport(content)
	assert.Nil(err)
	assert.Equal(4, len(rpt))

	content, err = fileToString("success.txt")
	assert.Nil(err)

	rpt, err = report.ParseReport(content)
	assert.Nil(err)
	assert.Zero(len(rpt))
}
