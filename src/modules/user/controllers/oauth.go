package user

import (
	"fmt"
	"io"
	"modules/user/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		//TODO: put your project's Client Id here.  To be got from https://code.google.com/apis/console
		ClientID: ucfg.Cfg.OAuth.ClientID,

		//TODO: put your project's Client Secret value here https://code.google.com/apis/console
		ClientSecret: ucfg.Cfg.OAuth.ClientSecret,
		Endpoint:     google.Endpoint,
		//To return your oauth2 code, Google will redirect the browser to this page that you have defined
		//TODO: This exact URL should also be added in your Google API console for this project within "API Access"->"Redirect URIs"
		RedirectURL: "http://home.rubi.gd/api/user/oauth/callback",

		//This is the 'scope' of the data that you are asking the user's permission to access. For getting user's info, this is the url that Google has defined.
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}
}

// reserveUser is the route for reserve a email/phone for registration
// @Route {
// 		url = /authenticate
//		method = get
//		200 = base.NormalResponse
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) oauthInit(ctx *gin.Context) {
	oauthCfg := getConfig()
	url := oauthCfg.AuthCodeURL("ok")
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
	oauthCfg := getConfig()
	//Get the code from the response
	code := ctx.Request.FormValue("code")
	t, err := oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		// TODO : better error
		u.BadResponse(ctx, err)
	}
	client := oauthCfg.Client(oauth2.NoContext, t)
	//now get user data based on the Transport which has the token
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		// TODO : better error
		u.BadResponse(ctx, err)
	}
	io.Copy(ctx.Writer, resp.Body)
}
