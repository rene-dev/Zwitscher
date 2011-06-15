package main

import (
	"github.com/garyburd/twister/oauth"
	)

type Accounts struct {
	Name        string
	Credentials *oauth.Credentials
	Maxreadid	int64
}

var accounts Accounts

func main() {
	accounts = Connect()
	accounts.Maxreadid = 0
	Gui()
}

