package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zallarak/db/cli/internal/colors"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: colors.Gray("User account management commands"),
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: colors.Gray("Get current user information"),
	RunE:  runUserMe,
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userMeCmd)
	
	// Silence usage on errors for clean error messages
	userCmd.SilenceUsage = true
	userMeCmd.SilenceUsage = true
}

func runUserMe(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Not logged in. Run ") + colors.Cyan("dbx auth login") + colors.White(" first"))
	}

	apiURL := viper.GetString("api-url")
	url := fmt.Sprintf("%s/v1/users/me", apiURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var response struct {
		User struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	outputFormat := viper.GetString("output")
	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(response.User)
	}

	// Clean field output
	fmt.Printf("%s %s\n", colors.FieldLabel("ID"), colors.Cyan(response.User.ID[:8]))
	fmt.Printf("%s %s\n", colors.FieldLabel("Email"), colors.White(response.User.Email))
	fmt.Printf("%s %s\n", colors.FieldLabel("Created"), colors.Gray(response.User.CreatedAt[:10]))

	return nil
}