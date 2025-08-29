package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"github.com/zallarak/db/cli/internal/colors"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: colors.Gray("Authentication commands"),
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: colors.Gray("Login to db.xyz"),
	Long:  colors.Gray("Login to your ") + colors.Cyan("db.xyz") + colors.Gray(" account and save authentication token"),
	RunE:  runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: colors.Gray("Logout from db.xyz"),
	Long:  colors.Gray("Remove stored authentication token"),
	RunE:  runLogout,
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: colors.Gray("Register a new account"),
	Long:  colors.Gray("Create a new ") + colors.Cyan("db.xyz") + colors.Gray(" account"),
	RunE:  runRegister,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(registerCmd)
	
	// Silence usage on errors for clean error messages
	authCmd.SilenceUsage = true
	loginCmd.SilenceUsage = true
	logoutCmd.SilenceUsage = true
	registerCmd.SilenceUsage = true
}

func runLogin(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(colors.Gray("Email: "))
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}
	email = email[:len(email)-1] // remove newline

	fmt.Print(colors.Gray("Password: "))
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // newline after password input

	password := string(passwordBytes)

	// Login request
	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := viper.GetString("api-url")
	resp, err := http.Post(apiURL+"/v1/auth/login", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.Unmarshal(body, &errorResp)
		if msg, ok := errorResp["error"].(string); ok {
			return fmt.Errorf(colors.Red("✗") + " " + colors.White("Login failed: ") + msg)
		}
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Login failed with status %d"), resp.StatusCode)
	}

	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Save token to config
	viper.Set("token", loginResp.Token)
	viper.Set("user.id", loginResp.User.ID)
	viper.Set("user.email", loginResp.User.Email)

	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		home, _ := os.UserHomeDir()
		configPath = home + "/.dbx.yaml"
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf(colors.SuccessIcon() + " " + colors.White("Logged in as ") + colors.Cyan(loginResp.User.Email) + "\n")
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	viper.Set("token", "")
	viper.Set("user.id", "")
	viper.Set("user.email", "")

	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		home, _ := os.UserHomeDir()
		configPath = home + "/.dbx.yaml"
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf(colors.SuccessIcon() + " " + colors.Gray("Logged out") + "\n")
	return nil
}

func runRegister(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(colors.Gray("Email: "))
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}
	email = email[:len(email)-1] // remove newline

	fmt.Print(colors.Gray("Password: "))
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // newline after password input

	password := string(passwordBytes)

	// Registration request
	registerReq := map[string]string{
		"email":    email,
		"password": password,
	}

	reqBody, err := json.Marshal(registerReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := viper.GetString("api-url")
	resp, err := http.Post(apiURL+"/v1/auth/register", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("registration request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.Unmarshal(body, &errorResp)
		if msg, ok := errorResp["error"].(string); ok {
			return fmt.Errorf(colors.Red("✗") + " " + colors.White("Registration failed: ") + msg)
		}
		return fmt.Errorf(colors.Red("✗") + " " + colors.White("Registration failed with status %d"), resp.StatusCode)
	}

	var registerResp struct {
		User struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.Unmarshal(body, &registerResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf(colors.SuccessIcon() + " " + colors.White("Account created for ") + colors.Cyan(registerResp.User.Email) + "\n")
	fmt.Printf(colors.Gray("Run ") + colors.Cyan("dbx auth login") + colors.Gray(" to sign in") + "\n")
	return nil
}