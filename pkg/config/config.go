package config

import (
	"flag"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	apiKey           string
	apiKeyHeaderName string
	bearerToken      string
	serverUrl        string
	username         string
	password         string
	authType         string
	ikey             string
	sKey             string
	accessToken      string
	refreshToken     string
	expiresAt        string
}

func Get() *Config {
	conf := &Config{}
	// creating a new flagset everytime the get function is called allows for different flagsets to exist
	// rather than a conflict to be created when generating a new config object (such as for tests)
	flags := flag.NewFlagSet("Passthrough Flag Set", flag.PanicOnError)
	flags.StringVar(&conf.username, "username", os.Getenv("USERNAME"), "Basic Auth username")
	flags.StringVar(&conf.password, "password", os.Getenv("PASSWORD"), "Basic Auth password")
	flags.StringVar(&conf.authType, "authType", os.Getenv("AUTH_TYPE"), "Auth Type")
	flags.StringVar(&conf.apiKeyHeaderName, "apiKeyHeaderName", os.Getenv("API_KEY_HEADER_NAME"), "API Key Header Name")
	flags.StringVar(&conf.apiKey, "apiKey", os.Getenv("API_KEY"), "API Key")
	flags.StringVar(&conf.bearerToken, "bearerToken", os.Getenv("BEARER_TOKEN"), "Bearer Token")
	flags.StringVar(&conf.ikey, "ikey", os.Getenv("IKEY"), "Duo Security IKey")
	flags.StringVar(&conf.sKey, "skey", os.Getenv("SKEY"), "Duo Security SKey")
	flags.StringVar(&conf.serverUrl, "serverUrl", os.Getenv("SERVER_URL"), "Server Url")
	flags.StringVar(&conf.accessToken, "accessToken", os.Getenv("ACCESS_TOKEN"), "Oauth2 Access Token")
	flags.StringVar(&conf.refreshToken, "refreshToken", os.Getenv("REFRESH_TOKEN"), "Oauth2 Refresh Token")
	flags.StringVar(&conf.expiresAt, "expiresAt", os.Getenv("EXPIRES_AT"), "Oauth2 Expires At")

	var arguments []string
	arguments = append(arguments, "os.Environ")
	flags.Parse(arguments)

	return conf
}

func (c *Config) GetApiKey() string {
	return c.apiKey
}

func (c *Config) GetApiKeyHeaderName() string {
	return c.apiKeyHeaderName
}

func (c *Config) GetBearerToken() string {
	return c.bearerToken
}

// GetAuthType returns the auth type accepted by the server
// Possible values include: API_KEY, BASIC_AUTH, HMAC
func (c *Config) GetAuthType() string {
	// convert all characters to upper case
	authType := strings.ToUpper(c.authType)
	// replace space, hyphen with underscore
	authType = strings.ReplaceAll(authType, " ", "_")
	authType = strings.ReplaceAll(authType, "%20", "_")
	authType = strings.ReplaceAll(authType, "-", "_")
	return authType
}

func (c *Config) GetUsernameAndPassword() (string, string) {
	return c.username, c.password
}

func (c *Config) GetAccessToken() string {
	return c.accessToken
}

func (c *Config) GetRefreshToken() string {
	return c.refreshToken
}

func (c *Config) GetExpiresAt() string {
	return c.expiresAt
}

func (c *Config) GetDuoIKeyAndSKey() (string, string) {
	return c.ikey, c.sKey
}

func (c *Config) GetServerURL() string {
	c.serverUrl = strings.TrimSuffix(c.serverUrl, "/")
	u, _ := url.Parse(c.serverUrl)
	if u.Scheme == "" {
		return "https://" + c.serverUrl
	} else {
		return c.serverUrl
	}
}

func (c *Config) GetServerHost() string {
	c.serverUrl = strings.TrimSuffix(c.serverUrl, "/")
	u, _ := url.Parse(c.serverUrl)
	if u.Scheme == "" {
		return u.Host
	} else {
		return u.Host
	}
}

func (c *Config) GetServerPath() string {
	c.serverUrl = strings.TrimSuffix(c.serverUrl, "/")
	u, _ := url.Parse(c.serverUrl)
	if u.Scheme == "" {
		return u.Path
	} else {
		return u.Path
	}
}