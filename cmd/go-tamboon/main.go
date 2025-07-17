package main

import (
	"fmt"
	"log"
	"os"

	"xvista/omise-challenge/internal/donation"
	"xvista/omise-challenge/internal/fileio"
	"xvista/omise-challenge/internal/reporting"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the CSV file as an argument.")
	}
	filePath := os.Args[1]

	client, err := donation.NewOmiseClientFromEnv()
	if err != nil {
		log.Fatalf("Error creating Omise client: %v", err)
	}

	csv, file, err := fileio.OpenEncryptedCSV(filePath)
	if err != nil {
		log.Fatalf("Error opening encrypted CSV file: %v", err)
	}
	defer file.Close()

	fmt.Println("performing donations...")

	stats, err := donation.ProcessDonations(csv, client)
	if err != nil {
		log.Fatalf("Error processing donations: %v", err)
	}

	fmt.Println("done.")
	fmt.Println()

	reporting.PrintSummary(stats)
}
