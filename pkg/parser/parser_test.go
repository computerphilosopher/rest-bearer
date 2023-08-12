package parser_test

import (
	"testing"

	"github.com/computerphilosopher/rest-bearer/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestMetaInfoParser(t *testing.T) {
	assert := assert.New(t)
	line := "MEDIUM: Insufficiently random value detected. [CWE-330]"
	metaInfo, err := parser.ParseMetaInfo(line)

	assert.Nil(err)
	assert.Equal(parser.Medium, metaInfo.Severity)
	assert.Equal("CWE-330", metaInfo.Code)
	assert.Equal("Insufficiently random value detected.", metaInfo.Description)
}

func TestLocationParser(t *testing.T) {
	line := "File: src/main/java/net/spy/memcached/ArcusClientPool.java:76"
	assert := assert.New(t)
	location, err := parser.ParseLocation(line)

	assert.Nil(err)
	assert.Equal("src/main/java/net/spy/memcached/ArcusClientPool.java", location.Path)
	assert.Equal(76, location.Line)
}

func TestSnippetParser(t *testing.T) {
	line := "137     md5.update(KeyUtil.getKeyBytes(k));"
	assert := assert.New(t)

	snippet, err := parser.ParseSnippet(line)
	assert.Nil(err)
	assert.Equal(snippet, "md5.update(KeyUtil.getKeyBytes(k));")
}
