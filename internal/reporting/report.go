package reporting

import (
	"fmt"
	"xvista/omise-challenge/internal/donation"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func formatAmount(amount float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", amount/100.0)
}

func PrintSummary(stats donation.Stats) {
	maxLen := len(formatAmount(stats.TotalDonations))

	fmt.Printf("%21s: THB %*s\n", "total received", maxLen, formatAmount(stats.TotalDonations))
	fmt.Printf("%21s: THB %*s\n", "successfully donated", maxLen, formatAmount(stats.SuccessfullyDonated))
	fmt.Printf("%21s: THB %*s\n", "faulty donation", maxLen, formatAmount(stats.FaultyDonations))
	fmt.Println()
	fmt.Printf("%21s: THB %*s\n", "average per person", maxLen, formatAmount(stats.AverageDonation))

	for i, donor := range stats.TopDonors {
		if i == 0 {
			fmt.Printf("%23s", "top donors: ")
		} else {
			fmt.Printf("%23s", "")
		}
		fmt.Printf("%s\n", donor)
	}
}
