# Gogen

🚫 **This repository has been archived**

Gogen was an application created by the Clear My Record team at Code for America in 2019 to support the state of California’s implementation of H&S § 11361.9 (marijuana conviction relief).

Since the deadline for the use case of this application by government was July 2020, the Clear My Record team discontinued maintenance of the application in September 2020.

**Please see below for more background on the project, and please reach out to clearmyrecord@codeforamerica.org for questions.**

## The Application

`gogen` is a command-line tool that takes in California Department of Justice (CA DOJ) .dat files containing criminal record data and identifies convictions eligible for relief under California's Proposition 64.  
The output is a bundle of CSV files that contain original data from the CA DOJ as well as eligibility info for relevant convictions.

This tool is intended for use with (and is packaged within) [B.E.A.R.](https://github.com/codeforamerica/bear).

## About

This application was developed by Code for America's [Clear My Record team](https://www.codeforamerica.org/programs/clear-my-record).

For more information about Clear My Record and how you might use this tool, visit our [H&S§11361.9 Implementation Toolkit](https://www.codeforamerica.org/programs/clear-my-record).

## Prerequisites

 - [Golang](https://golang.org/) install with `brew install golang`
 
## Cloning the project

Go projects live in a specific location on your file system under `/Users/[username]/go/src/[project]`.
Be sure to create the directory structure before cloning this project into `../go/src/gogen`

We recommend you add `../go/bin` to your path so you can run certain go tools from the command line 

## Setup

 - Change to project root directory `cd ~/go/src/gogen`
 - Install project dependencies with `go get ./...`
 - Install the Ginkgo test library with `go get github.com/onsi/ginkgo/ginkgo`
 - Install project test dependencies with `go get -t ./...`
 - Verify the tests are passing with `ginkgo -r`
 
## Running locally

This tool requires input files in the CA DOJ research file format. These files are tightly controlled for security and confidentiality purposes. 
We have created test fixture files that mimic the structure of the DOJ files, and you can use these to run the code on your local machine.

To compile and run gogen, run:
```
$ go run gogen run
    --input-doj=/Users/[username]/go/src/gogen/test_fixtures/no_headers.csv
    --county="SAN JOAQUIN"
    --compute-at=2020-07-01
    --eligibility-options=/path/to/bearConfig.json
    --outputs=/path/to/desired/output
```

If you would like to create a compiled artifact of gogen and install it (e.g. for use with BEAR), run the following commands from project root:
```
$ go build .
$ go install -i gogen
$ gogen run 
    --input-doj=/Users/[username]/go/src/gogen/test_fixtures/no_headers.csv
    --county="SAN JOAQUIN"
    --compute-at=2020-07-01
    --eligibility-options=/path/to/bearConfig.json
    --outputs=/path/to/desired/output
```

## Generating test data

We have provided a tool for generating sample data in the CA DOJ research file format for a given county, for use with Gogen.

To generate data, download the `./generate_test_data` script from the releases, make it executable (`chmod +x generate_test_data`) and run the following command:
```
$ ./generate_test_data 
    --county="LOS ANGELES"
    --target-size=50000
    --outputs=/path/to/file
```

County must be capitalized. Target size is the estimated number of rows in the produced file.

**Note:** This data is entirely fabricated and should only be used to understand the functionality of Gogen.
 
## License

MIT. Please see LICENSE and NOTICE.md.
