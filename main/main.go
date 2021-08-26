package main

import (
	"fmt"

	"github.org/smartybryan/callorg/classes"
)

const (
	RawCallingDataFilePath = "/Users/bryan/callorg/rawcallings.txt"
	CallingDataFilePath = "/Users/bryan/callorg/callings.csv"
	RawMembersDataFilePath = "/Users/bryan/callorg/rawmembership.txt"
	MembersDataFilePath = "/Users/bryan/callorg/callings.csv"

	MaxCallings = 100
	MaxMembers = 500
)


func main() {
	wardCallings := classes.NewCallings(MaxCallings)

	err := wardCallings.ParseCallingsFromRawData(RawCallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	err = wardCallings.SaveCallings(CallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	err = wardCallings.LoadCallings(CallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	totalCallings := 0
	for _, organization := range wardCallings.OrganizationOrder {
		fmt.Printf("%s\n", organization)
		for _, calling := range wardCallings.CallingMap[organization] {
			fmt.Printf("\t%s\t%s\t%s\t%t\n", calling.Name, calling.Holder, calling.Sustained, calling.CustomCalling)
			totalCallings++
		}
	}
	fmt.Printf("Total callings: %d\n", totalCallings)

	membership := classes.NewMembers(MaxMembers)
	err = membership.ParseMembersFromRawData(RawMembersDataFilePath)
	if err != nil {
		fmt.Println(err)
	}



	//fmt.Println()
	//
	//for _, name := range membership.SortedKeys() {
	//	memberRecord := membership[name]
	//	fmt.Printf("%s %s %s (%d) %t\n", memberRecord.Name, memberRecord.Gender, memberRecord.Birthday, memberRecord.Age(), memberRecord.Unbaptized)
	//}

	//TODO: save/load membership

}
