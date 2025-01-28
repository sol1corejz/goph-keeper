package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var serverURL, username, password string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate a user on the remote server",
	Run: func(cmd *cobra.Command, args []string) {
		// JSON payload for authentication
		payload := map[string]string{
			"username": username,
			"password": password,
		}
		data, _ := json.Marshal(payload)

		resp, err := http.Post(fmt.Sprintf("%s/login", serverURL), "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Fatalf("Error connecting to server: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Authentication successful!")
			// Save token or session details if required
		} else {
			fmt.Printf("Authentication failed with status: %d\n", resp.StatusCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVar(&serverURL, "server", "http://localhost:8080", "Remote server URL")
	loginCmd.Flags().StringVar(&username, "username", "", "Username for authentication")
	loginCmd.Flags().StringVar(&password, "password", "", "Password for authentication")
}
