package gotter

import (
	"bufio"
	"exec"
	"github.com/garyburd/twister/oauth"
	"github.com/garyburd/twister/web"
	"http"
	"iconv"
	"io/ioutil"
	"json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type Tweet struct {
	Text       string
	Identifier string "id_str"
	Source     string
	CreatedAt  string "created_at"
	User       struct {
		Name            string
		ScreenName      string "screen_name"
		FollowersCount  int    "followers_count"
		ProfileImageURL string "profile_image_url"
	}
	Place *struct {
		Id       string
		FullName string "full_name"
	}
	Entities struct {
		HashTags []struct {
			Indices [2]int
			Text    string
		}
		UserMentions []struct {
			Indices    [2]int
			ScreenName string "screen_name"
		}    "user_mentions"
		Urls []struct {
			Indices [2]int
			Url     string
		}
	}
}

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

func ClientAuth(requestToken *oauth.Credentials) (*oauth.Credentials, os.Error) {
	cmd := "xdg-open"
	url := oauthClient.AuthorizationURL(requestToken)

	args := []string{cmd, url}
	if syscall.OS == "windows" {
		cmd = "rundll32.exe"
		args = []string{cmd, "url.dll,FileProtocolHandler", url}
	} else if syscall.OS == "darwin" {
		cmd = "open"
		args = []string{cmd, url}
	}
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal("command not found:", err)
	}
	p, err := os.StartProcess(cmd, args, &os.ProcAttr{Dir: "", Files: []*os.File{nil, nil, os.Stderr}})
	if err != nil {
		log.Fatal("failed to start command:", err)
	}
	defer p.Release()

	print("PIN: ")
	stdin := bufio.NewReader(os.Stdin)
	b, err := stdin.ReadBytes('\n')
	if err != nil {
		log.Fatal("canceled")
	}

	if b[len(b)-2] == '\r' {
		b = b[0:len(b)-2]
	} else {
		b = b[0:len(b)-1]
	}
	accessToken, _, err := oauthClient.RequestToken(requestToken, string(b))
	if err != nil {
		log.Fatal("failed to request token:", err)
	}
	return accessToken, nil
}

func GetAccessToken(config map[string]string) (*oauth.Credentials, bool, os.Error) {
	oauthClient.Credentials.Token = config["ClientToken"]
	oauthClient.Credentials.Secret = config["ClientSecret"]

	authorized := false
	var token *oauth.Credentials
	accessToken, foundToken := config["AccessToken"]
	accessSecert, foundSecret := config["AccessSecret"]
	if foundToken && foundSecret {
		token = &oauth.Credentials{accessToken, accessSecert}
	} else {
		requestToken, err := oauthClient.RequestTemporaryCredentials("")
		if err != nil {
			log.Print("failed to request temporary credentials:", err)
			return nil, false, err
		}
		token, err = ClientAuth(requestToken)
		if err != nil {
			log.Print("failed to request temporary credentials:", err)
			return nil, false, err
		}

		config["AccessToken"] = token.Token
		config["AccessSecret"] = token.Secret
		authorized = true
	}
	return token, authorized, nil
}

func GetTweets(token *oauth.Credentials, url string, opt map[string]string) ([]Tweet, os.Error) {
	param := make(web.ParamMap)
	for k, v := range opt {
		param.Set(k, v)
	}
	oauthClient.SignParam(token, "GET", url, param)
	url = url + "?" + param.FormEncodedString()
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	var tweets []Tweet
	err = json.NewDecoder(res.Body).Decode(&tweets)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

func convert_utf8(s string) string {
	ic, err := iconv.Open("char", "UTF-8")
	if err != nil {
		return s
	}
	defer ic.Close()
	ret, _ := ic.Conv(s)
	return ret
}

func ShowTweets(tweets []Tweet, verbose bool) {
	if verbose {
		for i := len(tweets) - 1; i >= 0; i-- {
			name := convert_utf8(tweets[i].User.Name)
			user := convert_utf8(tweets[i].User.ScreenName)
			text := tweets[i].Text
			text = strings.Replace(text, "\r", "", -1)
			text = strings.Replace(text, "\n", " ", -1)
			text = strings.Replace(text, "\t", " ", -1)
			text = convert_utf8(text)
			println(user + ": " + name)
			println("  " + text)
			println("  " + tweets[i].Identifier)
			println("  " + tweets[i].CreatedAt)
			println()
		}
	} else {
		for i := len(tweets) - 1; i >= 0; i-- {
			user := convert_utf8(tweets[i].User.ScreenName)
			text := convert_utf8(tweets[i].Text)
			println(user + ": " + text)
		}
	}
}

func PostTweet(token *oauth.Credentials, url string, opt map[string]string) os.Error {
	param := make(web.ParamMap)
	for k, v := range opt {
		param.Set(k, v)
	}
	oauthClient.SignParam(token, "POST", url, param)
	res, err := http.PostForm(url, param.StringMap())
	if err != nil {
		log.Println("failed to post tweet:", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("failed to get timeline:", err)
		return err
	}
	return nil
}

func GetConfig() (string, map[string]string) {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, ".config")
	if syscall.OS == "windows" {
		home = os.Getenv("USERPROFILE")
		dir = filepath.Join(home, "Application Data")
	}
	_, err := os.Stat(dir)
	if err != nil {
		if os.Mkdir(dir, 0700) != nil {
			log.Fatal("failed to create directory:", err)
		}
	}
	dir = filepath.Join(dir, "twty")
	_, err = os.Stat(dir)
	if err != nil {
		if os.Mkdir(dir, 0700) != nil {
			log.Fatal("failed to create directory:", err)
		}
	}
	file := filepath.Join(dir, "settings.json")
	config := map[string]string{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		config["ClientToken"] = "lhCgJRAE1ECQzwVXfs5NQ"
		config["ClientSecret"] = "qk9i30vuzWHspsRttKsYrnoKSw9XBmWHdsis76z4"
	} else {
		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Fatal("could not unmarhal settings.json:", err)
		}
	}
	return file, config
}
