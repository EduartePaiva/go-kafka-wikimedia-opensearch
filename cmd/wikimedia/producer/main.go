package main

import "github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal"

const (
	TOPIC = "wikimedia.recentchange"
)

func main() {
	internal.TestWikimedia()
}
