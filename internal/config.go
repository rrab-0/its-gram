package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/viper"
)

var (
	region     = "ap-southeast-1"
	secretName = "its-gram"
)

func LoadConfig(cfgType string) error {
	if cfgType == "AWS" {
		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			return err
		}

		svc := secretsmanager.NewFromConfig(config)
		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"),
		}

		result, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			return err
		}

		viper.SetConfigType("json")
		err = viper.ReadConfig(strings.NewReader(*result.SecretString))
		if err != nil {
			return err
		}

		return nil
	}

	if cfgType == "DEV" {
		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cfgType %s not supported", cfgType)
}
