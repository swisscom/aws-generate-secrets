package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	generatesecrets "github.com/swisscom/aws-generate-secrets/pkg"
)

var args struct {
	AwsProfile    string `arg:"env:AWS_PROFILE"`
	SecretsConfig string `arg:"positional,required"`
}

var logger = logrus.New()

func main() {
	arg.MustParse(&args)
	err := generatesecrets.Do(args.SecretsConfig, args.AwsProfile)
	if err != nil {
		logger.Fatalf("unable to generate secrets: %v", err)
	}
}
