package donation

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

const (
	maxRetries = 5
)

func NewOmiseClientFromEnv() (*omise.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	omisePublicKey := os.Getenv("OMISE_PUBLIC_KEY")
	omiseSecretKey := os.Getenv("OMISE_SECRET_KEY")

	client, err := omise.NewClient(omisePublicKey, omiseSecretKey)
	if err != nil {
		return nil, fmt.Errorf("error creating Omise client: %w", err)
	}

	return client, nil
}

type DonationRecord struct {
	Name       string
	Amount     int64
	CardNumber string
	CVV        string
	ExpMonth   time.Month
	ExpYear    int
}

func (dr *DonationRecord) Clear() {
	dr.Name = ""
	dr.Amount = 0
	dr.CardNumber = ""
	dr.CVV = ""
	dr.ExpMonth = 0
	dr.ExpYear = 0
}

func (dr DonationRecord) Charge(client *omise.Client) error {
	var lastErr error
	backoff := time.Second

	for range maxRetries {
		result := &omise.Card{}
		if err := client.Do(result, &operations.CreateToken{
			Name:            dr.Name,
			Number:          dr.CardNumber,
			ExpirationMonth: dr.ExpMonth,
			ExpirationYear:  dr.ExpYear,
			SecurityCode:    dr.CVV,
		}); err != nil {
			if isRateLimitError(err) {
				time.Sleep(backoff)
				backoff *= 2
				lastErr = err
				continue
			}
			return err
		}
		token := result.Base.ID

		createCharge := &operations.CreateCharge{
			Amount:   dr.Amount,
			Currency: "thb",
			Card:     token,
		}
		charge := &omise.Charge{}
		err := client.Do(charge, createCharge)
		if err != nil {
			if isRateLimitError(err) {
				time.Sleep(backoff)
				backoff *= 2
				lastErr = err
				continue
			}
			return err
		}

		if charge.Status != "successful" {
			return fmt.Errorf("charge failed: %s", *charge.FailureMessage)
		}

		return nil
	}
	return fmt.Errorf("rate limit: %w", lastErr)
}

func isRateLimitError(err error) bool {
	var omErr *omise.Error
	if errors.As(err, &omErr) {
		return omErr.Code == "too_many_requests"
	}
	return false
}
