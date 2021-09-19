package main

import (
	"fmt"

	"github.org/smartybryan/callorg/engine"
	"github.org/smartybryan/callorg/util"
)

const (
	RawCallingDataFilePath = "/Users/bryan/callorg/rawcallings.txt"
	CallingDataFilePath = "/Users/bryan/callorg/callings.csv"
	RawMembersDataFilePath = "/Users/bryan/callorg/rawmembers.txt"
	MembersDataFilePath = "/Users/bryan/callorg/members.csv"

	MaxCallings = 300
	MaxMembers = 500
)

func main() {
	//parseAndPrintCallings()
	//parseAndPrintMembers()

	wardCallings := engine.NewCallings(MaxCallings)
	err := wardCallings.Load(CallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	wardMembers := engine.NewMembers(MaxMembers)
	err = wardMembers.Load(MembersDataFilePath)
	if err != nil {
		fmt.Println(err)
	}


}

func parseAndPrintCallings() {
	wardCallings := engine.NewCallings(MaxCallings)
	err := wardCallings.ParseCallingsFromRawData(RawCallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	err = wardCallings.Save(CallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	err = wardCallings.Load(CallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	totalCallings := 0
	//for _, organization := range wardCallings.OrganizationOrder {
	//	fmt.Printf("%s\n", organization)
	//	for _, calling := range wardCallings.CallingMap[organization] {
	//		fmt.Printf("\t%s\t%s\t%s\t%t\n",
	//			calling.Name, calling.Holder, util.PrintableDate(calling.Sustained), calling.CustomCalling)
	//		totalCallings++
	//	}
	//}
	fmt.Printf("Total callings: %d\n", totalCallings)
}

func parseAndPrintMembers() {
	membership := engine.NewMembers(MaxMembers)
	err := membership.ParseMembersFromRawData(RawMembersDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	err = membership.Save(MembersDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	err = membership.Load(MembersDataFilePath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()

	for _, name := range membership.SortedKeys() {
		memberRecord := membership[name]
		fmt.Printf("%s %s %s (%d) (eoy:%d) %t\n",
			memberRecord.Name, memberRecord.Gender, util.PrintableDate(memberRecord.Birthday),
			memberRecord.Age(), memberRecord.AgeByEndOfYear(), memberRecord.Unbaptized)
	}
}