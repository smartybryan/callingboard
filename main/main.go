package main

import "github.org/smartybryan/callorg/data"

const (
	RawCallingDataFilePath = "/Users/bryan/callorg/rawcallings.txt"
	CallingDataFilePath = "/Users/bryan/callorg/callings.csv"

	MaxCallings = 100
)


func main() {
	callings := data.NewCallings(MaxCallings)


}

