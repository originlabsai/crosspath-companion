package main

import (
	"fmt"
	"os"

	"github.com/crosspath/mcp-client/internal/commands"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via -ldflags "-X main.Version=X.Y.Z"
	Version = "0.0.0-dev"
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "cross_companion",
	Short: "CrossPath Companion - Connect local tools to cloud",
	Long: `CrossPath Companion connects local MCP servers to your CrossPath cloud chat.

Quick Start:
  cross_companion                 Launch interactive dashboard (default)
  cross_companion login           Authenticate (first time)
  cross_companion start           Start headless (for scripts/automation)

Commands:
  login                      Authenticate with CrossPath
  logout                     Log out and revoke this device
  start                      Start headless mode (no UI)
  status                     Show connection and auth status
  list                       List configured MCP servers
  add <name>                 Add a new MCP server
  remove <name>              Remove an MCP server
  service install/uninstall  Manage background service
  devices list/revoke        Manage connected devices

Examples:
  cross_companion                                 # Launch dashboard
  cross_companion login                           # First-time setup
  cross_companion add browser --command npx --args @browsermcp/mcp@latest
  cross_companion service install                 # Auto-start on login

Config: ~/.crosspath/mcp-config.yaml
Logs:   ~/.crosspath/logs/cross_companion.log`,
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		// When no subcommand is specified, try to launch TUI
		forceSetup, _ := cmd.Flags().GetBool("setup")
		return commands.RunTUIDefaultWithSetup(forceSetup)
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.Flags().Bool("setup", false, "Re-run the first-time setup wizard")

	// Add all commands
	rootCmd.AddCommand(commands.TUICmd)
	rootCmd.AddCommand(commands.LoginCmd)
	rootCmd.AddCommand(commands.StartCmd)
	rootCmd.AddCommand(commands.AddCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.RemoveCmd)
	rootCmd.AddCommand(commands.StatusCmd)
	rootCmd.AddCommand(commands.ServiceCmd)
	rootCmd.AddCommand(commands.DevicesCmd)
	rootCmd.AddCommand(commands.LogoutCmd)
	rootCmd.AddCommand(commands.DaemonCmd)
}

func main() {
	commands.AppVersion = Version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
