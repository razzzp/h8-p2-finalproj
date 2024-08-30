package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type InvoiceService struct {
	client   *http.Client
	hostname string
}

func NewInvoiceService() *InvoiceService {
	return &InvoiceService{
		client:   &http.Client{},
		hostname: "https://api.xendit.co/v2/invoices",
	}
}

func (is *InvoiceService) BuildSuccessUrl() string {
	return os.Getenv("XENDIT_INVOICE_CALLBACK")
}

func (is *InvoiceService) BuildFailureUrl() string {
	return os.Getenv("XENDIT_INVOICE_CALLBACK")
}

func (is *InvoiceService) GenerateInvoice(id uint,
	amount float64,
	desc string,
	email string,
) (string, error) {
	data := map[string]any{
		"external_id": fmt.Sprintf("%d", id),
		"amount":      amount,
		"description": desc,
		"customer": map[string]any{
			"email": email,
		},
		"success_redirect_url": is.BuildSuccessUrl(),
		"failure_redirect_url": is.BuildFailureUrl(),
		"payment_methods": []string{
			"CREDIT_CARD", "BCA", "BNI", "BSI", "BRI", "MANDIRI", "PERMATA",
		},
		"currency": "IDR",
	}
	body, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", is.hostname, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("XENDIT_API_KEY"), "")

	resp, err := is.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		bytes, err := io.ReadAll(resp.Body)
		if err == nil {
			fmt.Printf("%s\n", string(string(bytes)))
		}
		fmt.Printf("Here\n")
		return "", fmt.Errorf("xendit status code %d", resp.StatusCode)
	}

	var respBody map[string]any
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", nil
	}

	if url, ok := respBody["invoice_url"].(string); ok {
		return url, nil
	} else {
		return "", errors.New("failed to get invoice url from resp")
	}
}
