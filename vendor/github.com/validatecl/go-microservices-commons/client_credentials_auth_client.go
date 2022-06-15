package commons

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// AuthResponse response
type AuthResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

// MakeClientCredentialsOAUTHClient Crea cliente de autenticacion
func MakeClientCredentialsOAUTHClient(url, clientID, clientSecret string, timeout time.Duration, logger log.Logger) endpoint.Endpoint {

	encodeRequest := MakeEncodeClientCredentialsAuthRequest(clientID, clientSecret)

	return MakeHTTPClientBuilder("POST", url, timeout, encodeRequest, DecodeAuthResponse, logger).Build()
}

// DecodeAuthResponse decode de response de auth
func DecodeAuthResponse(_ context.Context, r *http.Response) (interface{}, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	if r.StatusCode != http.StatusOK {
		responseBody := buf.String()
		return nil,
			&ServiceError{
				Message: fmt.Sprintf("Error de callcenter service status: %d: %v", r.StatusCode, responseBody),
			}
	}

	authResponse := new(AuthResponse)

	if err := json.Unmarshal(buf.Bytes(), authResponse); err != nil {
		return nil, err
	}

	return authResponse, nil
}

// MakeEncodeClientCredentialsAuthRequest crea encode de auth request
func MakeEncodeClientCredentialsAuthRequest(clientID, clientSecret string) kithttp.EncodeRequestFunc {
	clientIDAndSecret := fmt.Sprintf("%s:%s", clientID, clientSecret)
	encodedClientAndSecret := base64.StdEncoding.EncodeToString([]byte(clientIDAndSecret))

	return func(ctx context.Context, r *http.Request, _ interface{}) error {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedClientAndSecret))

		formParams := url.Values{}
		formParams["grant_type"] = []string{"client_credentials"}
		r.Body = ioutil.NopCloser(strings.NewReader(formParams.Encode()))

		return nil
	}
}

//MakeOAuthClientCredentialsRequestEncodeMiddleware  crea middleware para encodeRequest
func MakeOAuthClientCredentialsRequestEncodeMiddleware(authEndpoint endpoint.Endpoint) func(kithttp.EncodeRequestFunc) kithttp.EncodeRequestFunc {
	return func(next kithttp.EncodeRequestFunc) kithttp.EncodeRequestFunc {
		return func(ctx context.Context, r *http.Request, in interface{}) error {
			res, err := authEndpoint(ctx, nil)
			if err != nil {
				return err
			}

			response := res.(*AuthResponse)
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", response.AccessToken))

			return next(ctx, r, in)
		}
	}
}
