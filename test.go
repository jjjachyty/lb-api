package main

import "fmt"

type People interface {
	Speak(string) string
}

type Stduent struct{}

func (stu *Stduent) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}
func main() {
	s := "我是好人"
	ss := []rune(s)
	ss[0] = rune('她')
	fmt.Println(string(ss))
}
