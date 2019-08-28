package main

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega/gstruct"
	"gogen/exporter"
	"gogen/utilities"
	"io/ioutil"
	"os/exec"
	path "path/filepath"
	"time"

	. "gogen/test_fixtures"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func GetOutputSummary(filePath string) exporter.Summary {
	bytes, _ := ioutil.ReadFile(filePath)
	var summary exporter.Summary
	json.Unmarshal(bytes, &summary)
	return summary
}

func GetErrors(filePath string) map[string]utilities.GogenError {
	bytes, _ := ioutil.ReadFile(filePath)
	var errors map[string]utilities.GogenError
	json.Unmarshal(bytes, &errors)
	return errors
}

var _ = Describe("gogen", func() {
	var (
		outputDir string
		pathToDOJ string
		err       error
	)
	It("can handle a csv with extra comma at the end of headers", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).ToNot(gbytes.Say("required"))

		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary.LineCount).To(Equal(38))
	})

	It("can handle an input file without headers", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "no_headers.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary.LineCount).To(Equal(38))
	})

	It("can accept a compute-at option for determining eligibility", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).ToNot(gbytes.Say("required"))
		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary.LineCount).To(Equal(38))
	})

	It("can accept a suffix for the output file names", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())
		dateSuffix := "Feb_8_2019_3.32.43.PM"

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", dateSuffix)
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		expectedDojResultsFileName := fmt.Sprintf("%v/All_Results_1_%s.csv", outputDir, dateSuffix)
		expectedCondensedFileName := fmt.Sprintf("%v/All_Results_Condensed_1_%s.csv", outputDir, dateSuffix)
		expectedConvictionsFileName := fmt.Sprintf("%v/Prop64_Results_1_%s.csv", outputDir, dateSuffix)
		expectedJsonOutputFileName := fmt.Sprintf("%v/gogen_%s.json", outputDir, dateSuffix)

		Ω(expectedDojResultsFileName).Should(BeAnExistingFile())
		Ω(expectedCondensedFileName).Should(BeAnExistingFile())
		Ω(expectedConvictionsFileName).Should(BeAnExistingFile())
		Ω(expectedJsonOutputFileName).Should(BeAnExistingFile())
	})

	It("validates required options", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(utilities.ERROR_EXIT))
		Eventually(session.Err).Should(gbytes.Say("missing required field: Run gogen --help for more info"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen%s.err", outputDir, "")

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		errors := GetErrors(expectedErrorFileName)
		Expect(errors).To(gstruct.MatchAllKeys(gstruct.Keys{
			"": gstruct.MatchAllFields(gstruct.Fields{
				"ErrorType":    Equal("OTHER"),
				"ErrorMessage": Equal("missing required field: Run gogen --help for more info"),
			}),
		}))
	})

	It("fails and reports errors for missing input files", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "missing.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())
		filenameSuffix := "a_suffix"

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(utilities.ERROR_EXIT))
		Eventually(session.Err).Should(gbytes.Say("open .*missing.csv: no such file or directory"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen_%s.err", outputDir, filenameSuffix)

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		errors := GetErrors(expectedErrorFileName)
		Expect(errors).To(gstruct.MatchAllKeys(gstruct.Keys{
			pathToDOJ: gstruct.MatchAllFields(gstruct.Fields{
				"ErrorType":    Equal("OTHER"),
				"ErrorMessage": Equal(fmt.Sprintf("open %s: no such file or directory", pathToDOJ)),
			}),
		}))
	})

	It("fails and reports errors for invalid input files", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "bad.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())
		filenameSuffix := "a_suffix"

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(utilities.ERROR_EXIT))
		Eventually(session.Err).Should(gbytes.Say("record on line 1: wrong number of fields"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen_%s.err", outputDir, filenameSuffix)

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		errors := GetErrors(expectedErrorFileName)
		Expect(errors).To(gstruct.MatchAllKeys(gstruct.Keys{
			pathToDOJ: gstruct.MatchAllFields(gstruct.Fields{
				"ErrorType":    Equal("PARSING"),
				"ErrorMessage": Equal("record on line 1: wrong number of fields"),
			}),
		}))
	})

	It("can accept path to eligibility options file", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
			"County":                  Equal("SACRAMENTO"),
			"LineCount":               Equal(38),
			"ProcessingTimeInSeconds": BeNumerically(">", 0),
			"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
			"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(5),
				"CountSubjectsNoConviction":           Equal(4),
				"CountSubjectsNoConvictionLast7Years": Equal(2),
			}),
			"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(5),
				"CountSubjectsNoConviction":           Equal(4),
				"CountSubjectsNoConvictionLast7Years": Equal(2),
			}),
			"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11357": Equal(3),
				"11358": Equal(7),
				"11359": Equal(8),
			}),
			"SubjectsWithProp64ConvictionCountInCounty": Equal(12),
			"Prop64FelonyConvictionsCountInCounty":      Equal(15),
			"Prop64NonFelonyConvictionsCountInCounty":   Equal(3),
			"SubjectsWithSomeReliefCount":               Equal(12),
			"ConvictionDismissalCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11357": Equal(2),
				"11358": Equal(6),
			}),
			"ConvictionReductionCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11359": Equal(1),
				"11360": Equal(0),
			}),
			"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
				"21 years or younger":                      Equal(1),
				"57 years or older":                        Equal(2),
				"Conviction occurred 10 or more years ago": Equal(1),
				"Individual is deceased":                   Equal(1),
				"Only has 11357-60 charges":                Equal(1),
			}),
		}))
	})

	Describe("Processing multiple input files", func() {
		It("nests and indexes the names of the results files for each input file", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
			inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

			pathToGogen, err := gexec.Build("gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV+","+inputCSV)
			countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
			computeAtFlag := "--compute-at=2019-11-11"
			eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			expectedDojResultsFile1Name := fmt.Sprintf("%v/All_Results_1.csv", outputDir)
			expectedDojResultsFile2Name := fmt.Sprintf("%v/All_Results_2.csv", outputDir)
			expectedCondensedFile1Name := fmt.Sprintf("%v/All_Results_Condensed_1.csv", outputDir)
			expectedCondensedFile2Name := fmt.Sprintf("%v/All_Results_Condensed_2.csv", outputDir)
			expectedConvictionsFile1Name := fmt.Sprintf("%v/Prop64_Results_1.csv", outputDir)
			expectedConvictionsFile2Name := fmt.Sprintf("%v/Prop64_Results_2.csv", outputDir)
			expectedJsonOutputFileName := fmt.Sprintf("%v/gogen.json", outputDir)

			Ω(expectedDojResultsFile1Name).Should(BeAnExistingFile())
			Ω(expectedDojResultsFile2Name).Should(BeAnExistingFile())
			Ω(expectedCondensedFile1Name).Should(BeAnExistingFile())
			Ω(expectedCondensedFile2Name).Should(BeAnExistingFile())
			Ω(expectedConvictionsFile1Name).Should(BeAnExistingFile())
			Ω(expectedConvictionsFile2Name).Should(BeAnExistingFile())
			Ω(expectedJsonOutputFileName).Should(BeAnExistingFile())
		})

		It("can aggregate statistics for multiple input files", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
			inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

			pathToGogen, err := gexec.Build("gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV+","+inputCSV)
			countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
			computeAtFlag := "--compute-at=2019-11-11"
			eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
			Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
				"County":                  Equal("SACRAMENTO"),
				"LineCount":               Equal(76),
				"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
				"ProcessingTimeInSeconds": BeNumerically(">", 0),
				"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(10),
					"CountSubjectsNoConviction":           Equal(8),
					"CountSubjectsNoConvictionLast7Years": Equal(4),
				}),
				"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(10),
					"CountSubjectsNoConviction":           Equal(8),
					"CountSubjectsNoConvictionLast7Years": Equal(4),
				}),
				"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357": Equal(6),
					"11358": Equal(14),
					"11359": Equal(16),
				}),
				"SubjectsWithProp64ConvictionCountInCounty": Equal(24),
				"Prop64FelonyConvictionsCountInCounty":      Equal(30),
				"Prop64NonFelonyConvictionsCountInCounty":   Equal(6),
				"SubjectsWithSomeReliefCount":               Equal(24),
				"ConvictionDismissalCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357": Equal(4),
					"11358": Equal(12),
				}),
				"ConvictionReductionCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11359": Equal(2),
					"11360": Equal(0),
				}),
				"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
					"21 years or younger":                      Equal(2),
					"57 years or older":                        Equal(4),
					"Conviction occurred 10 or more years ago": Equal(2),
					"Individual is deceased":                   Equal(2),
					"Only has 11357-60 charges":                Equal(2),
				}),
			}))
		})

		It("can return errors for multiple input files", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
			pathToValidDOJ, _, _ := ExtractFullCSVFixtures(pathToInputExcel)
			pathToBadDOJ, err := path.Abs(path.Join("test_fixtures", "bad.csv"))
			pathToMissingDOJ, err := path.Abs(path.Join("test_fixtures", "missing.csv"))

			pathToGogen, err := gexec.Build("gogen")
			Expect(err).ToNot(HaveOccurred())
			filenameSuffix := "a_suffix"

			pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", pathToValidDOJ+","+pathToBadDOJ+","+pathToMissingDOJ)

			countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
			computeAtFlag := "--compute-at=2019-11-11"
			dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)
			eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(utilities.ERROR_EXIT))
			Eventually(session.Err).Should(gbytes.Say("record on line 1: wrong number of fields"))
			Eventually(session.Err).Should(gbytes.Say("open .*missing.csv: no such file or directory"))

			expectedErrorFileName := fmt.Sprintf("%v/gogen_%s.err", outputDir, filenameSuffix)

			Ω(expectedErrorFileName).Should(BeAnExistingFile())
			errors := GetErrors(expectedErrorFileName)
			Expect(errors).To(gstruct.MatchAllKeys(gstruct.Keys{
				pathToMissingDOJ: gstruct.MatchAllFields(gstruct.Fields{
					"ErrorType":    Equal("OTHER"),
					"ErrorMessage": Equal(fmt.Sprintf("open %s: no such file or directory", pathToMissingDOJ)),
				}),
				pathToBadDOJ: gstruct.MatchAllFields(gstruct.Fields{
					"ErrorType":    Equal("PARSING"),
					"ErrorMessage": Equal("record on line 1: wrong number of fields"),
				}),
			}))
		})
	})
})
