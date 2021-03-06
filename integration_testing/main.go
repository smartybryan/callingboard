package main

import (
	"fmt"
	"io/ioutil"

	"github.org/smartybryan/callingboard/engine"
	"github.org/smartybryan/callingboard/util"
)

const (
	RawCallingDataFilePath = "/Users/bryan/callingboard/rawcallings.txt"
	CallingDataFilePath    = "/Users/bryan/callingboard/callings.csv"
	RawMembersDataFilePath = "/Users/bryan/callingboard/rawmembers.txt"
	MembersDataFilePath    = "/Users/bryan/callingboard/members.csv"

	MaxCallings = 300
	MaxMembers  = 500
)

func main() {
	parseAndPrintCallings()
	parseAndPrintMembers()
}

func parseAndPrintCallings() {
	wardCallings := engine.NewCallings(MaxCallings, CallingDataFilePath)
	data, err := ioutil.ReadFile(RawCallingDataFilePath)
	if err != nil {
		fmt.Println(err)
	}
	wardCallings.ParseCallingsFromRawData(data)

	_, err = wardCallings.Save()
	if err != nil {
		fmt.Println(err)
	}

	err = wardCallings.Load()
	if err != nil {
		fmt.Println(err)
	}

	totalCallings := 0
	for _, organization := range wardCallings.OrganizationOrder {
		fmt.Printf("%s\n", organization)
		for _, calling := range wardCallings.CallingMap[organization] {
			fmt.Printf("\t%s\t%s\t%s\t%t\n",
				calling.Name, calling.Holder, util.PrintableDate(calling.Sustained), calling.CustomCalling)
			totalCallings++
		}
	}
	fmt.Printf("Total callings: %d\n", totalCallings)
}

func parseAndPrintMembers() {
	membership := engine.NewMembers(MaxMembers, MembersDataFilePath)
	data, err := ioutil.ReadFile(RawMembersDataFilePath)
	if err != nil {
		fmt.Println(err)
	}
	membership.ParseMembersFromRawData(data)

	_, err = membership.Save()
	if err != nil {
		fmt.Println(err)
	}

	err = membership.Load()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()

	for _, name := range membership.GetMembers(engine.AllMembers) {
		memberRecord := membership.GetMemberRecord(name)
		fmt.Printf("%s %d\n",
			memberRecord.Name, memberRecord.Eligibility)
	}
}
