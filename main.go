package main


import (
        "fmt"
)

var BuildVcRef string
var BuildTime string
var BuildVersion string
var BuildSource string

func main() {
		fmt.Printf("Hello world, version: %s\n", BuildVcRef)
		fmt.Printf("Hello world, version: %s\n", BuildTime)
		fmt.Printf("Hello world, version: %s\n", BuildVersion)
		fmt.Printf("Hello world, version: %s\n", BuildSource)
}