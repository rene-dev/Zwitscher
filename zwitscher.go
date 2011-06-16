package main

import ()

var accounts Accounts

func main() {
	accounts = Connect()
	accounts.Maxreadid = 0
	Gui()
}

