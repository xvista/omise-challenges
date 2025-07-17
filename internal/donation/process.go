package donation

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/omise/omise-go"
)

const (
	MaxTopDonors     = 3
	ThrottleDuration = 100 * time.Millisecond
)

type Stats struct {
	TotalDonations      float64
	SuccessfullyDonated float64
	FaultyDonations     float64
	AverageDonation     float64
	TopDonors           []string
}

func ProcessDonations(
	csvReader *csv.Reader,
	client *omise.Client,
) (Stats, error) {
	err := godotenv.Load()
	if err != nil {
		return Stats{}, fmt.Errorf("error loading .env file: %w", err)
	}

	workerCount := min(runtime.NumCPU(), 10)
	jobs := make(chan DonationRecord, workerCount)
	var wg sync.WaitGroup

	donationRecords := map[string]int64{}
	var mu sync.Mutex

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range jobs {
				if err := d.Charge(client); err != nil {
					log.Printf("Error charging donation: %v", err)
				} else {
					mu.Lock()
					donationRecords[d.Name] += d.Amount
					mu.Unlock()
				}
				d.Clear()
				time.Sleep(ThrottleDuration)
				log.Printf("Throttling for %v", ThrottleDuration)
			}
		}()
	}

	totalDonations := float64(0)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record: %v", err)
			continue
		}
		amount, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Println("Error parsing amount:", err)
			continue
		}
		expMonth, err := strconv.Atoi(record[4])
		if err != nil {
			log.Println("Error parsing expiration month:", err)
			continue
		}
		expYear, err := strconv.Atoi(record[5])
		if err != nil {
			log.Println("Error parsing expiration year:", err)
			continue
		}
		d := DonationRecord{
			Name:       record[0],
			Amount:     amount,
			CardNumber: record[2],
			CVV:        record[3],
			ExpMonth:   time.Month(expMonth),
			ExpYear:    expYear,
		}
		totalDonations += float64(amount)
		jobs <- d
	}
	close(jobs)
	wg.Wait()

	donors := make([]string, 0, len(donationRecords))
	successfullyDonated := float64(0)
	for donor, amount := range donationRecords {
		donors = append(donors, donor)
		successfullyDonated += float64(amount)
	}

	sort.SliceStable(donors, func(i, j int) bool {
		return donationRecords[donors[i]] > donationRecords[donors[j]]
	})

	faultyDonations := totalDonations - successfullyDonated
	averageDonation := 0.0
	if len(donationRecords) > 0 {
		averageDonation = successfullyDonated / float64(len(donationRecords))
	}
	topDonors := donors
	if len(topDonors) > MaxTopDonors {
		topDonors = topDonors[:MaxTopDonors]
	}

	stats := Stats{
		TotalDonations:      totalDonations,
		SuccessfullyDonated: successfullyDonated,
		FaultyDonations:     faultyDonations,
		AverageDonation:     averageDonation,
		TopDonors:           topDonors,
	}

	return stats, nil
}
