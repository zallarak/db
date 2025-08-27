package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Database instance management commands",
}

var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List database instances",
	RunE:  runInstanceList,
}

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new database instance",
	RunE:  runInstanceCreate,
}

var instanceDeleteCmd = &cobra.Command{
	Use:   "delete [instance-id]",
	Short: "Delete a database instance",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstanceDelete,
}

func init() {
	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceDeleteCmd)

	// Instance create flags
	instanceCreateCmd.Flags().String("project", "", "Project ID (required)")
	instanceCreateCmd.Flags().String("name", "", "Instance name (required)")
	instanceCreateCmd.Flags().String("plan", "nano", "Instance plan (nano, lite, pro, pro-heavy)")
	instanceCreateCmd.Flags().Int("pg-version", 16, "PostgreSQL version")
	instanceCreateCmd.Flags().Int("disk", 0, "Disk size in GiB (optional)")

	instanceCreateCmd.MarkFlagRequired("project")
	instanceCreateCmd.MarkFlagRequired("name")

	// Instance delete flags
	instanceDeleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")
}

func runInstanceList(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("not logged in. Run 'dbx login' first")
	}

	projectID, _ := cmd.Flags().GetString("project")
	if projectID == "" {
		return fmt.Errorf("project flag is required")
	}

	apiURL := viper.GetString("api-url")
	url := fmt.Sprintf("%s/v1/projects/%s/instances", apiURL, projectID)
	
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
		Instances []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Plan      string `json:"plan"`
			PgVersion int    `json:"pg_version"`
			Status    string `json:"status"`
			FQDN      string `json:"fqdn"`
			CreatedAt string `json:"created_at"`
		} `json:"instances"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	outputFormat := viper.GetString("output")
	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(response.Instances)
	}

	// Table output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tPLAN\tPG\tSTATUS\tFQDN\tCREATED")
	for _, instance := range response.Instances {
		fqdn := instance.FQDN
		if fqdn == "" {
			fqdn = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\t%s\n", 
			instance.ID[:8], instance.Name, instance.Plan, instance.PgVersion, 
			instance.Status, fqdn, instance.CreatedAt[:10])
	}
	return w.Flush()
}

func runInstanceCreate(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("not logged in. Run 'dbx login' first")
	}

	projectID, _ := cmd.Flags().GetString("project")
	name, _ := cmd.Flags().GetString("name")
	plan, _ := cmd.Flags().GetString("plan")
	pgVersion, _ := cmd.Flags().GetInt("pg-version")
	diskSize, _ := cmd.Flags().GetInt("disk")

	createReq := map[string]interface{}{
		"name":       name,
		"plan":       plan,
		"pg_version": pgVersion,
	}

	if diskSize > 0 {
		createReq["disk_gib"] = diskSize
	}

	reqBody, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := viper.GetString("api-url")
	url := fmt.Sprintf("%s/v1/projects/%s/instances", apiURL, projectID)
	
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

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("Instance creation initiated: %s\n", name)
	if jobID, ok := response["job_id"].(string); ok {
		fmt.Printf("Job ID: %s\n", jobID)
		fmt.Println("Use 'dbx instance list --project', createReq[projectID] to track progress")
	}
	return nil
}

func runInstanceDelete(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("not logged in. Run 'dbx login' first")
	}

	instanceID := args[0]
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		fmt.Printf("Are you sure you want to delete instance %s? (y/N): ", instanceID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled")
			return nil
		}
	}

	apiURL := viper.GetString("api-url")
	url := fmt.Sprintf("%s/v1/instances/%s", apiURL, instanceID)
	
	req, err := http.NewRequest("DELETE", url, nil)
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

	fmt.Printf("Instance %s deletion initiated\n", instanceID)
	return nil
}