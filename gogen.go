package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gogen/data"
	"gogen/exporter"
	"gogen/test_fixtures"
	"gogen/utilities"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

const VERSION = "0.2.11"

var defaultOpts struct{}

type runOpts struct {
	OutputFolder       string `long:"outputs" description:"The folder in which to place result files"`
	DOJFiles           string `long:"input-doj" description:"The files containing criminal histories from CA DOJ"`
	County             string `long:"county" short:"c" description:"The county for which eligibility will be computed"`
	ComputeAt          string `long:"compute-at" description:"The date for which eligibility will be evaluated, ex: 2020-10-31"`
	EligibilityOptions string `long:"eligibility-options" description:"File containing options for which eligibility logic to apply"`
	FileNameSuffix     string `long:"file-name-suffix" hidden:"true" description:"string to append to file names"`
}

type exportTestCSVOpts struct {
	ExcelFixturePath string `long:"excel-fixture-path" short:"e" description:"Path to a county's excel fixture file to generate test CSVs"`
	OutputFolder     string `long:"outputs" short:"o" description:"The folder in which to place result files"`
}

type versionOpts struct{}

var opts struct {
	Version   versionOpts       `command:"version" description:"Print the version"`
	Run       runOpts           `command:"run" description:"Process an input DOJ file and produce an annotated DOJ data file"`
	ExportCSV exportTestCSVOpts `command:"export-test-csv" description:"Export example data files from excel fixtures"`
}

func (r runOpts) Execute(args []string) error {

	var processingStartTime time.Time

	utilities.SetErrorFileName(utilities.GenerateFileName(r.OutputFolder, "gogen%s.err", r.FileNameSuffix))

	if r.OutputFolder == "" || r.DOJFiles == "" || r.County == "" || r.EligibilityOptions == "" {
		utilities.ExitWithError(errors.New("missing required field: Run gogen --help for more info"))
	}

	inputFiles := strings.Split(r.DOJFiles, ",")

	computeAtDate := time.Now()

	if r.ComputeAt != "" {
		computeAtOption, err := time.Parse("2006-01-02", r.ComputeAt)
		if err != nil {
			utilities.ExitWithError(errors.New("invalid --compute-at date: Must be a valid date in the format YYYY-MM-DD"))
		} else {
			computeAtDate = computeAtOption
		}
	}

	var configurableEligibilityFlow data.ConfigurableEligibilityFlow

	var options data.EligibilityOptions
	optionsFile, err := os.Open(r.EligibilityOptions)
	if err != nil {
		utilities.ExitWithError(err)
	}
	defer optionsFile.Close()

	optionsBytes, err := ioutil.ReadAll(optionsFile)
	if err != nil {
		utilities.ExitWithError(err)
	}

	err = json.Unmarshal(optionsBytes, &options)
	if err != nil {
		utilities.ExitWithError(err)
	}
	configurableEligibilityFlow, err = data.NewConfigurableEligibilityFlow(options, r.County)
	if err != nil {
		utilities.ExitWithError(err)
	}

	runErrors := make(map[string]utilities.GogenError)
	var runSummary exporter.Summary
	outputJsonFilePath := utilities.GenerateFileName(r.OutputFolder, "gogen%s.json", r.FileNameSuffix)

	err = os.MkdirAll(r.OutputFolder, os.ModePerm)
	if err != nil {
		utilities.ExitWithError(err)
	}

	for fileIndex, inputFile := range inputFiles {
		processingStartTime = time.Now()
		fileIndex = fileIndex + 1
		dojInformation, gogenErr := data.NewDOJInformation(inputFile, computeAtDate, configurableEligibilityFlow)
		if gogenErr.ErrorType != "" {
			runErrors[inputFile] = gogenErr
			continue
		}
		countyEligibilities := dojInformation.DetermineEligibility(r.County, configurableEligibilityFlow)

		dismissAllProp64Eligibilities := dojInformation.DetermineEligibility(r.County, data.EligibilityFlows["DISMISS ALL PROP 64"])
		dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility(r.County, data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

		dojFilePath := utilities.GenerateIndexedFileName(r.OutputFolder, "All_Results_%d%s.csv", fileIndex, r.FileNameSuffix)
		condensedFilePath := utilities.GenerateIndexedFileName(r.OutputFolder, "All_Results_Condensed_%d%s.csv", fileIndex, r.FileNameSuffix)
		prop64ConvictionsFilePath := utilities.GenerateIndexedFileName(r.OutputFolder, "Prop64_Results_%d%s.csv", fileIndex, r.FileNameSuffix)

		dojWriter, err := exporter.NewDOJWriter(dojFilePath)
		if err != nil {
			runErrors[inputFile] = utilities.GogenError{ErrorType: "OTHER", ErrorMessage: err.Error()}
			continue
		}
		condensedDojWriter, err := exporter.NewCondensedDOJWriter(condensedFilePath)
		if err != nil {
			runErrors[inputFile] = utilities.GogenError{ErrorType: "OTHER", ErrorMessage: err.Error()}
			continue
		}
		prop64ConvictionsDojWriter, err := exporter.NewDOJWriter(prop64ConvictionsFilePath)
		if err != nil {
			runErrors[inputFile] = utilities.GogenError{ErrorType: "OTHER", ErrorMessage: err.Error()}
			continue
		}

		dataExporter := exporter.NewDataExporter(
			dojInformation,
			countyEligibilities,
			dismissAllProp64Eligibilities,
			dismissAllProp64AndRelatedEligibilities,
			dojWriter,
			condensedDojWriter,
			prop64ConvictionsDojWriter)

		fileSummary := dataExporter.Export(r.County, configurableEligibilityFlow)
		runSummary = dataExporter.AccumulateSummaryData(runSummary, fileSummary)
	}

	if encounteredErrors(runErrors) {
		utilities.ExitWithErrors(runErrors)
	}

	ExportSummary(runSummary, processingStartTime, outputJsonFilePath)
	return nil
}

func encounteredErrors(runErrors map[string]utilities.GogenError) bool {
	for _, value := range runErrors {
		if value.ErrorType != "" {
			return true
		}
	}
	return false

}

func ExportSummary(summary exporter.Summary, startTime time.Time, filePath string) {
	summary.ProcessingTimeInSeconds = time.Since(startTime).Seconds()

	s, err := json.Marshal(summary)
	if err != nil {
		utilities.ExitWithError(err)
	}
	err = ioutil.WriteFile(filePath, s, 0644)
	if err != nil {
		utilities.ExitWithError(err)
	}
}

func (e exportTestCSVOpts) Execute(args []string) error {
	if e.ExcelFixturePath != "" {
		inputCSV, expectedResultsCSV, err := test_fixtures.ExportFullCSVFixtures(e.ExcelFixturePath, e.OutputFolder)
		if err != nil {
			fmt.Println("Extracting test CSVs failed")
			os.Exit(1)
		}

		fmt.Println("Wrote input CSV at: " + inputCSV)
		fmt.Println("Wrote expected results CSV at: " + expectedResultsCSV)

		openCommand := exec.Command("open", e.OutputFolder)
		err = openCommand.Run()
		if err != nil {
			panic(err)
		}
	} else {
		return errors.New("something went wrong")
	}

	return nil
}

func (v versionOpts) Execute(args []string) error {
	fmt.Println(VERSION)
	return nil
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
