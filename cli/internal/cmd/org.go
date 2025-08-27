package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Organization management commands",
}

var orgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List organizations",
	RunE:  runOrgList,
}

var orgSelectCmd = &cobra.Command{
	Use:   "select [org-id]",
	Short: "Select default organization",
	Args:  cobra.ExactArgs(1),
	RunE:  runOrgSelect,
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(orgListCmd)
	orgCmd.AddCommand(orgSelectCmd)
}

func runOrgList(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("not logged in. Run 'dbx login' first")
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
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
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

	// Table output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tROLE\tCREATED")
	for _, org := range response.Orgs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", org.ID[:8], org.Name, org.Role, org.CreatedAt[:10])
	}
	return w.Flush()
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

	fmt.Printf("Selected organization: %s\n", orgID)
	return nil
}