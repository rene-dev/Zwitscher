package main

import (
	"github.com/garyburd/twister/oauth"
	)

type Accounts struct {
	Name        string
	Credentials *oauth.Credentials
}

var accounts Accounts

func main() {
	accounts = Connect()
	Gui()
}

