package main

import (
	"fmt"

	"github.org/smartybryan/callorg/types"
)

const (
	RawCallingDataFilePath = "/Users/bryan/callorg/rawcallings.txt"
	CallingDataFilePath = "/Users/bryan/callorg/callings.csv"

	MaxCallings = 100
)


func main() {
	wardCallings := types.NewCallings(MaxCallings)
	_ = wardCallings.ParseCallingsFromRawData(RawCallingDataFilePath)

	totalCallings := 0
	for _, organization := range wardCallings.OrganizationOrder {
		fmt.Printf("%s\n", organization)
		for _, calling := range wardCallings.CallingMap[organization] {
			fmt.Printf("\t%s\t%s\t%s\t%t\n", calling.Name, calling.Holder, calling.Sustained, calling.CustomCalling)
			totalCallings++
		}
	}

	fmt.Printf("Total callings: %d\n", totalCallings)
}
