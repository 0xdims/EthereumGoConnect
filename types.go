package main

type request struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

type EtherscanApi struct {
	Status  string
	Message string
	Result  string
}