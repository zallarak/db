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

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: colors.Gray("Project management commands"),
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: colors.Gray("List projects in the current organization"),
	RunE:  runProjectList,
}

var projectCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: colors.Gray("Create a new project"),
	Args:  cobra.ExactArgs(1),
	RunE:  runProjectCreate,
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectCreateCmd)
	
	// Silence usage on errors for clean error messages
	projectCmd.SilenceUsage = true
	projectListCmd.SilenceUsage = true
	projectCreateCmd.SilenceUsage = true
}

func runProjectList(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Not logged in. Run ") + colors.Cyan("dbx auth login") + colors.White(" first"))
	}

	orgID := viper.GetString("default-org")
	if orgID == "" {
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("No default organization selected. Run ") + colors.Cyan("dbx org select <org-id>") + colors.White(" first"))
	}

	apiURL := viper.GetString("api-url")
	url := fmt.Sprintf("%s/v1/orgs/%s/projects", apiURL, orgID)
	
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
		Projects []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			CreatedAt string `json:"created_at"`
		} `json:"projects"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	outputFormat := viper.GetString("output")
	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(response.Projects)
	}

	// Clean table output
	if len(response.Projects) == 0 {
		fmt.Println(colors.Gray("No projects found"))
		return nil
	}

	fmt.Printf("%s   %s   %s\n", 
		colors.TableHeader("id"),
		colors.TableHeader("name"), 
		colors.TableHeader("created"))
	
	for _, project := range response.Projects {
		fmt.Printf("%s   %s   %s\n",
			colors.Cyan(project.ID[:8]),
			colors.White(project.Name),
			colors.Gray(project.CreatedAt[:10]))
	}
	return nil
}

func runProjectCreate(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Not logged in. Run ") + colors.Cyan("dbx auth login") + colors.White(" first"))
	}

	orgID := viper.GetString("default-org")
	if orgID == "" {
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("No default organization selected. Run ") + colors.Cyan("dbx org select <org-id>") + colors.White(" first"))
	}

	projectName := args[0]

	createReq := map[string]string{
		"name": projectName,
	}

	reqBody, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := viper.GetString("api-url")
	url := fmt.Sprintf("%s/v1/orgs/%s/projects", apiURL, orgID)
	
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

	var project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf(colors.SuccessIcon() + " " + colors.White("Created project: ") + colors.Cyan(project.Name) + colors.Gray(" (") + colors.Cyan(project.ID[:8]) + colors.Gray(")") + "\n")
	return nil
}