package utils

import (
	"fmt"
	"net/http"
)

// PingAddress pings the given address using the provided client.
func PingAddress(address string, client *http.Client) error {
	resp, err := client.Get(address)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}
