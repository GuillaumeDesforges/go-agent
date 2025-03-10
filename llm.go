package main

type ILlm interface {
	UpdateModel(model string) error
	Query(input string) (string, error)
}
