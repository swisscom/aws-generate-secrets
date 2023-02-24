package generatesecrets_test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	generatesecrets "github.com/swisscom/aws-generate-secrets/pkg"
	"os"
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

func TestDo(t *testing.T) {
	generatesecrets.SetLogLevel(logrus.DebugLevel)
	err := generatesecrets.Do("../resources/examples/postgres.yaml", os.Getenv("AWS_PROFILE"))
	assert.Nil(t, err)
}
