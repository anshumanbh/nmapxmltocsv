package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	gn "github.com/tomsteele/go-nmap"
)

type config struct {
	inputFile  string
	outputFile string
}

type portscan struct {
	Hostname     string
	IPAddress    string
	Port         int
	Protocol     string
	Servicename  string
	Servicestate string
}

var (
	cfg config
)

func loadConfig() {
	inputFile := flag.String("inputFile", "in.xml", "Input File Path")
	outputFile := flag.String("outputFile", "", "Output File Path")

	flag.Parse()

	cfg = config{
		inputFile:  *inputFile,
		outputFile: *outputFile,
	}

	validateParams()
}

func validateParams() {
	var didError = false

	if cfg.inputFile == "" {
		log.Println("Error: inputFile is a required parameter, cannot be blank.")
		didError = true
	}
	if cfg.outputFile == "" {
		log.Println("Error: outputFile is a required parameter, cannot be blank.")
		didError = true
	}

	if didError {
		log.Fatalf("Usage: nmapxmltocsv -inputFile TODO -outputFile TODO")
		os.Exit(1)
	}
}

func exists(path string) (bool, int64, error) {
	fi, err := os.Stat(path)
	if err == nil {
		return true, fi.Size(), nil
	}
	if os.IsNotExist(err) {
		return false, int64(0), nil
	}
	return false, int64(0), err
}

func ensureFilePathExists(filepath string) error {
	value := false
	fsize := int64(0)

	for (value == false) || (fsize == int64(0)) {
		i, s, err := exists(filepath)
		if err != nil {
			log.Println("Failed to determine if the file exists or not..")
		}
		value = i
		fsize = s
	}

	log.Println("File exists:", value)
	log.Println("File size:", fsize)

	return nil
}

func writeResultsToCsv(scanResults []portscan, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Couldn't create the output file: %v", err)
		return err
	}
	defer outputFile.Close()

	err = gocsv.MarshalFile(&scanResults, outputFile)
	if err != nil {
		fmt.Printf("Couldn't marshal the output file: %v", err)
		return err
	}
	return nil
}

func main() {

	loadConfig()

	err := ensureFilePathExists(cfg.inputFile)
	if err != nil {
		log.Fatalf("Couldn't ensure whether the file exists or not: %v", err)
	}

	var allResults []portscan

	fb, err := ioutil.ReadFile(cfg.inputFile)
	if err != nil {
		log.Fatalf("Couldn't open the file: %v", err)
	}

	n, err := gn.Parse(fb)
	if err != nil {
		log.Fatalf("Couldn't parse the file: %v", err)
	}

	for _, host := range n.Hosts {

		for _, ip := range host.Addresses {

			for _, port := range host.Ports {
				scanResult := portscan{
					Hostname:     host.Hostnames[0].Name,
					IPAddress:    ip.Addr,
					Port:         port.PortId,
					Protocol:     port.Protocol,
					Servicename:  port.Service.Name,
					Servicestate: port.State.State,
				}
				allResults = append(allResults, scanResult)
			}
		}

	}

	err = writeResultsToCsv(allResults, cfg.outputFile)
	if err != nil {
		log.Fatalf("Couldn't write to the output file: %v", err)
	}

	fmt.Println("Results saved to: " + cfg.outputFile)
}
