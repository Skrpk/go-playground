package cognito

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	c "gitlab.com/saas-template1/go-monolith/config"
	"net/http"
)

type AwsCognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientId   string
}

func NewCognitoClient(conf c.AwsConfig) *AwsCognitoClient {
	config := &aws.Config{
		Region: aws.String(conf.AwsRegion),
	}

	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", conf.AwsRegion, conf.UserPoolId)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	fmt.Println(">>>>>>>>", url)
	fmt.Println(">>>>>>>>", resp.Body)

	sess, err := session.NewSession(config)
	client := cognito.New(sess)

	if err != nil {
		panic(err)
	}

	return &AwsCognitoClient{
		cognitoClient: client,
		appClientId:   conf.AppClientId,
	}
}

func (c *AwsCognitoClient) SignUp(email string, password string) (error, string) {
	user := &cognito.SignUpInput{
		Username: aws.String(email),
		Password: aws.String(password),
		ClientId: aws.String(c.appClientId),
	}

	result, err := c.cognitoClient.SignUp(user)

	if err != nil {
		return err, ""
	}
	return nil, result.String()
}

func (c *AwsCognitoClient) SignIn(email string, password string) (error, *struct {
	AccessToken  string
	RefreshToken string
}) {
	initiateAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		}),
		ClientId: aws.String(c.appClientId),
	}
	result, err := c.cognitoClient.InitiateAuth(initiateAuthInput)

	if err != nil {
		return err, nil
	}

	return nil, &struct {
		AccessToken  string
		RefreshToken string
	}{
		AccessToken:  *result.AuthenticationResult.AccessToken,
		RefreshToken: *result.AuthenticationResult.RefreshToken,
	}
}

func (c *AwsCognitoClient) Refresh(refreshToken string) (error, *struct {
	AccessToken string
}) {
	initiateAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"REFRESH_TOKEN": refreshToken,
		}),
		ClientId: aws.String(c.appClientId),
	}
	result, err := c.cognitoClient.InitiateAuth(initiateAuthInput)

	if err != nil {
		return err, nil
	}

	return nil, &struct {
		AccessToken string
	}{
		AccessToken: *result.AuthenticationResult.AccessToken,
	}
}
