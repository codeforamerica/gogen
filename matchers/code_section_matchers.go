package matchers

import (
	"regexp"
)

var prop64matcher = regexp.MustCompile(`(11357|11358|11359|11360)`)
var relatedChargeMatcher = regexp.MustCompile(`(32\s*PC|11366\s*HS|11366\.5\s*[^\(]HS|11366\.5\s*\([ABC]\)\s*HS)`)

var Prop64MatchersByCodeSection = map[string]*regexp.Regexp{
	"11357":                 regexp.MustCompile(`11357.*`),
	"11358":                 regexp.MustCompile(`11358.*`),
	"11359":                 regexp.MustCompile(`11359.*`),
	"11360":                 regexp.MustCompile(`11360.*`),
}
var Section11357SubSectionMatcher = regexp.MustCompile(`11357\(([A-D])\)`)

func ExtractProp64Section(codeSection string) (bool, string) {
	if IsProp64Charge(codeSection) {
		return true, prop64matcher.FindStringSubmatch(codeSection)[1]
	} else {
		return false, ""
	}
}

func ExtractRelatedChargeSection(codeSection string) (bool, string) {
	if IsRelatedCharge(codeSection) {
		return true, relatedChargeMatcher.FindStringSubmatch(codeSection)[1]
	} else {
		return false, ""
	}
}

func IsProp64Charge(codeSection string) bool {
	return prop64matcher.Match([]byte(codeSection))
}

func IsRelatedCharge(codeSection string) bool {
	return relatedChargeMatcher.Match([]byte(codeSection))
}

func Extract11357SubSection(codeSection string) (bool, string) {
	result := Section11357SubSectionMatcher.FindStringSubmatch(codeSection)
	if result != nil {
		return true, result[1]
	} else {
		return false, ""
	}
}