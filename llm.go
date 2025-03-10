package main

type ILlm interface {
	Query(input string) (string, error)
}
