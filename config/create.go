package config

import (
	"fmt"
	"log"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

func NewConfigFromYAML() Config {

	cfg := Config{}

	err := yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

}
