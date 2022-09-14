package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	PaymentsConf *PaymentsConfig
	AwsConfig    *AwsConfig
}

type PaymentsConfig struct {
	PaymentPlans     []string
	PlansIDs         map[string]string
	StripeKey        string
	WebhookSignature string
	SuccessUrl       string
	CancelUrl        string
}

type AwsConfig struct {
	AwsRegion   string
	AppClientId string
	UserPoolId  string
}

type YamlConfig struct {
	PaymentPlans []string `yaml:"payment_plans"`
}

func NewConfig() *Config {
	filename, _ := filepath.Abs("./config/main.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)

	if err != nil {
		panic(err)
	}

	plansIDs := make(map[string]string)

	for _, key := range yamlConfig.PaymentPlans {
		planEnvKey := fmt.Sprintf("PAYMENT_PLAN_%s", strings.ToUpper(key))
		planId := os.Getenv(planEnvKey)
		if planId == "" {
			panic(errors.New(fmt.Sprintf("Failed to extract plan id from env var %s", planEnvKey)))
		}
		plansIDs[key] = planId
	}

	stripeKey := os.Getenv("STRIPE_KEY")
	webhookSignature := os.Getenv("STRIPE_WEBHOOK_SIGNATURE")
	successUrl := os.Getenv("STRIPE_SUCCESS_URL")
	cancelUrl := os.Getenv("STRIPE_CANCEL_URL")

	if stripeKey == "" {
		panic(errors.New("stripe key is not provided"))
	}

	if webhookSignature == "" {
		panic(errors.New("webhook signature is not provided"))
	}

	if successUrl == "" || cancelUrl == "" {
		panic(errors.New("redirect url data is not provided"))
	}

	awsRegion := os.Getenv("AWS_REGION")
	appClientId := os.Getenv("AWS_APP_CLIENT_ID")
	userPoolId := os.Getenv("AWS_USER_POOL_ID")

	if awsRegion == "" || appClientId == "" || userPoolId == "" {
		panic(errors.New("aws config is not provided"))
	}

	return &Config{
		PaymentsConf: &PaymentsConfig{
			PlansIDs:         plansIDs,
			StripeKey:        stripeKey,
			WebhookSignature: webhookSignature,
			SuccessUrl:       successUrl,
			CancelUrl:        cancelUrl,
		},
		AwsConfig: &AwsConfig{
			AwsRegion:   awsRegion,
			AppClientId: appClientId,
			UserPoolId:  userPoolId,
		},
	}
}
