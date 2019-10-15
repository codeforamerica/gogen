package data

import (
	"fmt"
	"gogen/matchers"
	"time"
)

type dismissAllProp64EligibilityFlow struct {
}

func (ef dismissAllProp64EligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County, county) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, county)
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef dismissAllProp64EligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef dismissAllProp64EligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if matchers.IsProp64Charge(row.CodeSection) {
		info.SetEligibleForDismissal("Dismiss all Prop 64 charges")
	}
}

func (ef dismissAllProp64EligibilityFlow) checkRelevancy(codeSection string, convictionCounty string, flowCounty string) bool {
	return convictionCounty == flowCounty && matchers.IsProp64Charge(codeSection)
}

type dismissAllProp64AndRelatedEligibilityFlow struct {
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County, county) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, county)
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if matchers.IsProp64Charge(row.CodeSection) || matchers.IsRelatedCharge(row.CodeSection) {
		info.SetEligibleForDismissal("Dismiss all Prop 64 and related charges")
	}
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) checkRelevancy(codeSection string, convictionCounty string, flowCounty string) bool {
	return convictionCounty == flowCounty && (matchers.IsProp64Charge(codeSection) || matchers.IsRelatedCharge(codeSection))
}

type findRelatedChargesFlow struct {
}

func (ef findRelatedChargesFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, event := range subject.ArrestsAndConvictions {
		// if event is in county
		// 	if event is arrest with Prop 64 charge
		// check to see if there are any related charge convictions with same cyc_count
		if event.County == county {
			info := NewEligibilityInfo(event, subject, comparisonTime, county)
			ef.BeginEligibilityFlow(info, event, subject)
			infos[event.Index] = info
		}
	}
	return infos
}

func (ef findRelatedChargesFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef findRelatedChargesFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	stpOrder := row.CountOrder[0:3]
	prop64ArrestInSameCycle := subject.CyclesWithProp64Arrest[stpOrder]
	fmt.Println(prop64ArrestInSameCycle)
	if matchers.IsRelatedCharge(row.CodeSection) && prop64ArrestInSameCycle {
		info.SetPotentiallyEligibleRelatedConviction()
	}
}
