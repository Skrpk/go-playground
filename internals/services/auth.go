package services

type AuthClient interface {
	SignUp(email string, password string) (error, string)
	SignIn(email string, password string) (error, *struct {
		AccessToken  string
		RefreshToken string
	})
	Refresh(refreshToken string) (error, *struct {
		AccessToken string
	})
}

type SignInOutput struct {
	AccessToken string
}

type AuthService struct {
	Client AuthClient
}

func NewAuthService(client AuthClient) *AuthService {
	return &AuthService{
		Client: client,
	}
}

func (aS *AuthService) SignIn(email string, password string) (error, *struct {
	AccessToken  string
	RefreshToken string
}) {
	return aS.Client.SignIn(email, password)
}

func (aS *AuthService) Refresh(refreshToken string) (error, *struct {
	AccessToken string
}) {
	return aS.Client.Refresh(refreshToken)
}

func (aS *AuthService) SignUp(email string, password string) (error, string) {
	return aS.Client.SignUp(email, password)
}
