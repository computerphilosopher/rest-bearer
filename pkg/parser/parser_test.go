package parser_test

import (
	"testing"

	"github.com/computerphilosopher/rest-bearer/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	assert := assert.New(t)
	line := "MEDIUM: Insufficiently random value detected. [CWE-330]"
	vulnerability, err := parser.ParseVulnerability(line)

	assert.Nil(err)
	assert.Equal(parser.Medium, vulnerability.Severity)
	assert.Equal("CWE-330", vulnerability.Code)
	assert.Equal("Insufficiently random value detected.", vulnerability.Description)
}
