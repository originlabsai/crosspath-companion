package commands

import (
	"fmt"

	"github.com/crosspath/mcp-client/internal/config"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show client status and configuration",
	Long:  `Display the current status of the MCP client, including authentication and server configuration.`,
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("📊 CrossPath Companion Status")
	fmt.Println()

	// Authentication status
	if cfg.AuthToken != "" {
		fmt.Println("🔐 Authentication: ✅ Logged in")
		if cfg.UserID != "" {
			fmt.Printf("   User ID: %s\n", cfg.UserID)
		}
	} else {
		fmt.Println("🔐 Authentication: ❌ Not logged in")
		fmt.Println("   Run 'cross_companion login' to authenticate")
	}
	fmt.Println()

	// Backend configuration
	fmt.Printf("🌐 Backend: %s\n", cfg.BackendURL)
	fmt.Println()

	// Server configuration
	enabledServers := cfg.GetEnabledServers()
	fmt.Printf("📦 MCP Servers: %d configured (%d enabled)\n", len(cfg.MCPServers), len(enabledServers))

	if len(enabledServers) > 0 {
		fmt.Println()
		fmt.Println("Enabled servers:")
		for _, server := range enabledServers {
			fmt.Printf("  • %s (%s)\n", server.Name, server.Type)
		}
	}
	fmt.Println()

	// Configuration file
	fmt.Printf("📁 Config file: %s\n", config.GetConfigPath())

	// Next steps
	if cfg.AuthToken == "" {
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("1. Login: cross_companion login")
	} else if len(enabledServers) == 0 {
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("1. Add servers: cross_companion add <name> --path <server-path>")
	} else {
		fmt.Println()
		fmt.Println("✅ Ready to start!")
		fmt.Println("   Run: cross_companion start")
	}

	return nil
}
