package user

import (
	"fmt"
	"modules/user/config"
	"net/http"

	"encoding/json"

	"modules/user/aaa"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		//TODO: put your project's Client Id here.  To be got from https://code.google.com/apis/console
		ClientID: ucfg.Cfg.OAuth.ClientID,

		//TODO: put your project's Client Secret value here https://code.google.com/apis/console
		ClientSecret: ucfg.Cfg.OAuth.ClientSecret,
		Endpoint:     google.Endpoint,
		//To return your oauth2 code, Google will redirect the browser to this page that you have defined
		//TODO: This exact URL should also be added in your Google API console for this project within "API Access"->"Redirect URIs"
		RedirectURL: ucfg.Cfg.OAuth.RedirectURI,

		//This is the 'scope' of the data that you are asking the user's permission to access.
		// For getting user's info, this is the url that Google has defined.
		Scopes: []string{
			//"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}
}

// reserveUser is the route for reserve a email/phone for registration
// @Route {
// 		url = /authenticate/:action
//		method = get
//		200 = base.NormalResponse
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) oauthInit(ctx *gin.Context) {
	// Check the action, valid ones are login/register
	state := ctx.Param("action")
	if state != "login" && state != "register" {
		u.NotFoundResponse(ctx, nil)
		return
	}
	oauthCfg := getConfig()
	url := oauthCfg.AuthCodeURL(state)
	fmt.Print(oauthCfg)
	//redirect user to that page
	ctx.Redirect(http.StatusFound, url)
}

// reserveUser is the route for reserve a email/phone for registration
// @Route {
// 		url = /oauth/callback
//		method = get
//		200 = base.NormalResponse
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) oauthCallback(ctx *gin.Context) {
	state := ctx.Request.FormValue("state")
	var redirect string
	switch state {
	case "login":
		redirect = ucfg.Cfg.OAuth.LoginRedirect
	case "register":
		redirect = ucfg.Cfg.OAuth.RegisterRedirect
	default:
		u.NotFoundResponse(ctx, nil)
		return
	}
	eString := ctx.Request.FormValue("error")
	if eString != "" {
		// the request is canceled
		ctx.Redirect(http.StatusFound, redirect+"?error=unauthorized")
		return
	}

	oauthCfg := getConfig()
	//Get the code from the response
	code := ctx.Request.FormValue("code")
	t, err := oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		ctx.Redirect(http.StatusFound, redirect+"?error="+err.Error())
		return
	}
	client := oauthCfg.Client(oauth2.NoContext, t)
	//now get user data based on the Transport which has the token
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		ctx.Redirect(http.StatusFound, redirect+"?error="+err.Error())
		return
	}
	gp := googleResponse{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&gp)
	if err != nil {
		logrus.Panic(err)
	}

	m := aaa.NewAaaManager()
	var (
		usr *aaa.User
	)
	switch state {
	case "login":
		var token string
		token, usr, err = m.LoginUserByOAuth(gp.Email)
		if err != nil {
			ctx.Redirect(http.StatusFound, redirect+"?error="+err.Error())
			return
		}
		ctx.Redirect(http.StatusFound, redirect+"?token="+token)
	case "register":
		usr, err = m.RegisterUserByContact(gp.Email)
		if err != nil {
			ctx.Redirect(http.StatusFound, redirect+"?error=already_regsitered")
			return
		}
		token := m.GetNewToken(usr.Token)
		ctx.Redirect(http.StatusFound, redirect+"?token="+token)
	}

}
