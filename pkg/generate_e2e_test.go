//go:build e2e

package generatesecrets

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestDo(t *testing.T) {
	generatesecrets.SetLogLevel(logrus.DebugLevel)
	err := generatesecrets.Do("../resources/examples/postgres.yaml", os.Getenv("AWS_PROFILE"))
	assert.Nil(t, err)
}
