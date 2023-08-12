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
