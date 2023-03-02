package generatesecrets_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	generatesecrets "github.com/swisscom/aws-generate-secrets/pkg"
	"strings"
	"testing"
)

func TestGenerateSecretString(t *testing.T) {
	theSecret := generatesecrets.GenerateSecretString(generatesecrets.GeneratorConfig{
		Length:  10,
		Charset: "special",
	})
	fmt.Printf("generated=%v\n", theSecret)
	assert.Len(t, theSecret, 10)

	theSecret2 := generatesecrets.GenerateSecretString(generatesecrets.GeneratorConfig{
		Length:  8,
		Charset: "standard",
	})
	fmt.Printf("generated2=%v\n", theSecret2)
	assert.Len(t, theSecret2, 8)
}

func TestGenerateLowercaseSecretString(t *testing.T) {
	theSecret := generatesecrets.GenerateSecretString(generatesecrets.GeneratorConfig{
		Length:  10,
		Charset: "lowercase",
	})
	fmt.Printf("generated=%v\n", theSecret)
	assert.Len(t, theSecret, 10)
	assert.Equal(t, strings.ToLower(theSecret), theSecret)
}
