package accfin

type WsAuthorization struct {
	Result              WsAuthResult `json:"result" bson:"result"`
	TargetUrl           string       `json:"targetUrl" bson:"targetUrl"`
	Success             bool         `json:"success" bson:"success"`
	Error               bool         `json:"error" bson:"error"`
	UnAuthorizedRequest bool         `json:"unAuthorizedRequest" bson:"unAuthorizedRequest"`
	Abp                 bool         `json:"_abp" bson:"_abp"`
}

type WsAuthResult struct {
	AccessToken                   string      `json:"accessToken" bson:"accessToken"`
	EncryptedAccessToken          string      `json:"encryptedAccessToken" bson:"encryptedAccessToken"`
	ExpireInSeconds               int         `json:"expireInSeconds" bson:"expireInSeconds"`
	ShouldResetPassword           bool        `json:"shouldResetPassword" bson:"shouldResetPassword"`
	PasswordResetCode             string      `json:"passwordResetCode" bson:"passwordResetCode"`
	UserId                        int         `json:"userId" bson:"userId"`
	RequiresTwoFactorVerification bool        `json:"requiresTwoFactorVerification" bson:"requiresTwoFactorVerification"`
	TwoFactorAuthProviders        interface{} `json:"twoFactorAuthProviders" bson:"twoFactorAuthProviders"`
	TwoFactorRememberClientToken  interface{} `json:"twoFactorRememberClientToken" bson:"twoFactorRememberClientToken"`
	ReturnUrl                     string      `json:"returnUrl" bson:"returnUrl"`
	RefreshToken                  string      `json:"refreshToken" bson:"refreshToken"`
}
