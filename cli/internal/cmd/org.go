package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zallarak/db/cli/internal/colors"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: colors.Gray("Organization management commands"),
}

var orgListCmd = &cobra.Command{
	Use:   "list",
	Short: colors.Gray("List organizations"),
	RunE:  runOrgList,
}

var orgSelectCmd = &cobra.Command{
	Use:   "select [org-id]",
	Short: colors.Gray("Select default organization"),
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgSelect,
}

var orgCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: colors.Gray("Create a new organization"),
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgCreate,
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(orgListCmd)
	orgCmd.AddCommand(orgSelectCmd)
	orgCmd.AddCommand(orgCreateCmd)
	
	// Silence usage on errors for clean error messages
	orgCmd.SilenceUsage = true
	orgListCmd.SilenceUsage = true
	orgSelectCmd.SilenceUsage = true
	orgCreateCmd.SilenceUsage = true
}

func runOrgList(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Not logged in. Run ") + colors.Cyan("dbx auth login") + colors.White(" first"))
	}

	apiURL := viper.GetString("api-url")
	req, err := http.NewRequest("GET", apiURL+"/v1/orgs", nil)
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
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Request failed with status %d"), resp.StatusCode)
	}

	var response struct {
		Orgs []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Role      string `json:"role"`
			CreatedAt string `json:"created_at"`
		} `json:"orgs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	outputFormat := viper.GetString("output")
	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(response.Orgs)
	}

	// Clean table output
	if len(response.Orgs) == 0 {
		fmt.Println(colors.Gray("No organizations found"))
		return nil
	}

	fmt.Printf("%s   %s   %s   %s\n", 
		colors.TableHeader("id"),
		colors.TableHeader("name"), 
		colors.TableHeader("role"),
		colors.TableHeader("created"))
	
	for _, org := range response.Orgs {
		fmt.Printf("%s   %s   %s   %s\n",
			colors.Cyan(org.ID[:8]),
			colors.White(org.Name),
			colors.Gray(org.Role),
			colors.Gray(org.CreatedAt[:10]))
	}
	return nil
}

func runOrgSelect(cmd *cobra.Command, args []string) error {
	orgID := args[0]
	
	viper.Set("default-org", orgID)
	
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		home, _ := os.UserHomeDir()
		configPath = home + "/.dbx.yaml"
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf(colors.SuccessIcon() + " " + colors.White("Selected organization: ") + colors.Cyan(orgID) + "\n")
	return nil
}

func runOrgCreate(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Not logged in. Run ") + colors.Cyan("dbx auth login") + colors.White(" first"))
	}

	orgName := args[0]

	createReq := map[string]string{
		"name": orgName,
	}

	reqBody, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := viper.GetString("api-url")
	url := fmt.Sprintf("%s/v1/orgs", apiURL)
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
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

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var response struct {
		Org struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"org"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf(colors.SuccessIcon() + " " + colors.White("Created organization: ") + colors.Cyan(response.Org.Name) + colors.Gray(" (") + colors.Cyan(response.Org.ID[:8]) + colors.Gray(")") + "\n")
	return nil
}