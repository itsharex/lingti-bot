package cmd

import (
	"fmt"
	"os"

	"github.com/pltanton/lingti-bot/internal/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage the lingti-bot service",
	Long:  `Install, uninstall, start, stop, or check the status of the lingti-bot service.`,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install lingti-bot as a system service",
	Long:  `Install lingti-bot as a system service (requires root/admin privileges).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current executable path
		execPath, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Installing lingti-bot service...")
		if err := service.Install(execPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service installed successfully!")
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the lingti-bot service",
	Long:  `Uninstall the lingti-bot service (requires root/admin privileges).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling lingti-bot service...")
		if err := service.Uninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "Error uninstalling service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service uninstalled successfully!")
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the lingti-bot service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service started!")
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the lingti-bot service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.Stop(); err != nil {
			fmt.Fprintf(os.Stderr, "Error stopping service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service stopped!")
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the lingti-bot service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.Restart(); err != nil {
			fmt.Fprintf(os.Stderr, "Error restarting service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Service restarted!")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the lingti-bot service",
	Run: func(cmd *cobra.Command, args []string) {
		installed := service.IsInstalled()
		running := service.IsRunning()

		binaryPath, configPath := service.Paths()

		fmt.Println("=== lingti-bot Service Status ===")
		fmt.Println()
		fmt.Printf("Installed: %v\n", installed)
		fmt.Printf("Running:   %v\n", running)
		fmt.Println()
		fmt.Printf("Binary:    %s\n", binaryPath)
		fmt.Printf("Config:    %s\n", configPath)
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(installCmd)
	serviceCmd.AddCommand(uninstallCmd)
	serviceCmd.AddCommand(startCmd)
	serviceCmd.AddCommand(stopCmd)
	serviceCmd.AddCommand(restartCmd)
	serviceCmd.AddCommand(statusCmd)
}
