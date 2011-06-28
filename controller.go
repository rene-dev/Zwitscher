package main

import (
	"github.com/mattn/go-gtk/gdkpixbuf"
	"http"
	"json"
	"io"
	"io/ioutil"
	"strings"
	"gotter"
	"github.com/garyburd/twister/oauth"
	"log"
)

type Accounts struct {
	Name        string
	Credentials *oauth.Credentials
	Maxreadid	int64
}

func Connect() Accounts{
	var account Accounts
	file, config := gotter.GetConfig("zwitscher")
	config["ClientToken"] = "lhCgJRAE1ECQzwVXfs5NQ"
	config["ClientSecret"] = "qk9i30vuzWHspsRttKsYrnoKSw9XBmWHdsis76z4"
	token, authorized, err := gotter.GetAccessToken(config)
	if err != nil {
		log.Fatal("faild to get access token:", err)
	}
	if authorized {
		b, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
		err = ioutil.WriteFile(file, b, 0700)
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
	}
	account.Credentials = token
	return account
}

func url2pixbuf(url string) *gdkpixbuf.GdkPixbuf {
	r, err := http.Get(url)
	if err != nil {
		return nil
	}
	t := r.Header.Get("Content-Type")
	b := make([]byte, r.ContentLength)
	if _, err = io.ReadFull(r.Body, b); err != nil {
		return nil
	}
	var loader *gdkpixbuf.GdkPixbufLoader
	if strings.Index(t, "jpeg") >= 0 {
		loader, _ = gdkpixbuf.PixbufLoaderWithMimeType("image/jpeg")
	} else {
		loader, _ = gdkpixbuf.PixbufLoaderWithMimeType("image/png")
	}
	loader.SetSize(24, 24)
	loader.Write(b)
	loader.Close()
	return loader.GetPixbuf()
}

func SendTweet(text string) {
	gotter.PostTweet(accounts.Credentials, "https://api.twitter.com/1/statuses/update.json", map[string]string{"status": text})
}
