package models

type Secret struct {
	Filter string
	Name   string
	Value  any
}

type SecretJson struct {
	Filter string `json:"filter"`
	Name   string `json:"name"`
	Value  any    `json:"value"`
}
