package generatesecrets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

var SpecialChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-=[]{};,./?")

var logger = logrus.New()

func SetLogLevel(level logrus.Level) {
	logger.SetLevel(level)
}

func Do(inputFile string, awsProfileOverride string) error {
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	var cfg Config
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&cfg)
	if err != nil {
		return err
	}

	if awsProfileOverride != "" {
		cfg.Profile = awsProfileOverride
	}

	// Init AWS SDK
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedConfigProfile(cfg.Profile),
	)
	if err != nil {
		return err
	}

	logger.Debugf("using AWS profile %s", cfg.Profile)
	sm := secretsmanager.NewFromConfig(awsCfg)

	// Test AWS Config
	for _, v := range cfg.Secrets {
		createdSecret, err := generateSecret(ctx, sm, v)
		if err != nil {
			return err
		}

		if createdSecret.Name != nil && createdSecret.ARN != nil {
			logger.Infof("Secret created: %v (%v)", *createdSecret.Name, *createdSecret.ARN)
		}
	}

	return nil
}

func generateSecret(
	ctx context.Context,
	client *secretsmanager.Client,
	secret Secret,
) (*secretsmanager.CreateSecretOutput, error) {
	if client == nil {
		return nil, fmt.Errorf("SecretsManager client cannot be nil")
	}

	// Check if secret exists
	existingSecrets, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{Filters: []types.Filter{
		{Key: types.FilterNameStringTypeName,
			Values: []string{
				secret.Name,
			}},
	}})

	if err != nil {
		return nil, fmt.Errorf("cannot check if secret exists: %v", err)
	}

	// Secrets in SecretsManager cannot be replaced (grace period of 7 - 30 days)
	for _, v := range existingSecrets.SecretList {
		if v.Name != nil && *v.Name == secret.Name {
			logger.Warnf("secret %v exists already: %v", *v.Name, *v.ARN)
			return nil, nil
		}
	}

	secretContent := map[string]string{}

	err = generateSecretContent(secretContent, secret)
	if err != nil {
		return nil, fmt.Errorf("unable to generate secret content: %v", err)
	}

	jsonBinary, err := json.Marshal(secretContent)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal JSON: %v", err)
	}

	jsonString := string(jsonBinary)

	createdSecret, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         &secret.Name,
		SecretString: &jsonString,
	})
	if err != nil {
		return nil, err
	}

	return createdSecret, nil
}

func generateSecretContent(m map[string]string, secret Secret) error {
	if m == nil {
		return fmt.Errorf("map cannot be nil")
	}
	for k, v := range secret.Generator {
		m[k] = GenerateSecretString(v)
	}
	return nil
}

func GenerateSecretString(v GeneratorConfig) string {
	if len(v.Chars) == 0 {
		switch v.Charset {
		case CharsetStd:
			return uniuri.NewLenChars(v.Length, uniuri.StdChars)
		case CharsetSpecial:
			return uniuri.NewLenChars(v.Length, SpecialChars)
		}
		if v.Charset != "" {
			logger.Warnf("unknown charset \"%s\", using %s", v.Charset, CharsetStd)
		}
		return uniuri.NewLenChars(v.Length, uniuri.StdChars)
	}

	return uniuri.NewLenChars(v.Length, []byte(v.Chars))
}
