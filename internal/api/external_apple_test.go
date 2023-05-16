package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

func (ts *ExternalTestSuite) TestSignupExternalApple() {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/authorize?provider=apple", nil)
	w := httptest.NewRecorder()
	ts.API.handler.ServeHTTP(w, req)
	ts.Require().Equal(http.StatusFound, w.Code)
	u, err := url.Parse(w.Header().Get("Location"))
	ts.Require().NoError(err, "redirect url parse failed")
	q := u.Query()
	ts.Equal(ts.Config.External.Apple.RedirectURI, q.Get("redirect_uri"))
	ts.Equal(ts.Config.External.Apple.ClientID, q.Get("client_id"))
	ts.Equal("code", q.Get("response_type"))
	ts.Equal("email name", q.Get("scope"))

	claims := ExternalProviderClaims{}
	_, err = parseJWTTokenWithClaims(q.Get("state"), ts.Config, &claims)
	ts.Require().NoError(err)

	ts.Equal("apple", claims.Provider)
	ts.Equal(ts.Config.SiteURL, claims.SiteURL)
}
