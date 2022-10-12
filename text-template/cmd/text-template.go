package main

import (
	"os"
	"text/template"
)

type Inventory struct {
	Material string
	Count    int
}

func main() {
	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}\n")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 5; i++ {
		sweaters.Count += i
		err = tmpl.Execute(os.Stdout, sweaters)
		if err != nil {
			panic(err)
		}
	}
}
