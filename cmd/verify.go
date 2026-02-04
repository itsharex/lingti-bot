package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pltanton/lingti-bot/internal/platforms/relay"
	"github.com/spf13/cobra"
)

var (
	verifyPlatform   string
	verifyServerURL  string
	verifyTimeout    int
	// WeCom credentials
	verifyWeComCorpID  string
	verifyWeComAgentID string
	verifyWeComSecret  string
	verifyWeComToken   string
	verifyWeComAESKey  string
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify callback URL for messaging platforms",
	Long: `Connect to the cloud relay service to verify callback URL configuration.

This command is used to complete the initial callback URL verification when
setting up a messaging platform integration. It sends your credentials to the
cloud relay server so it can respond to the platform's verification request.

Supported platforms:
  - wecom: Enterprise WeChat (企业微信)

Usage flow:
  1. Create your app in the platform's admin console
  2. Run this verify command with your credentials
  3. Configure callback URL in the platform (e.g., https://bot.lingti.com/wecom)
  4. Save the configuration - the platform will send a verification request
  5. Once verified, you can stop this command and use 'relay' for normal operation

Example for WeCom:
  lingti-bot verify \
    --platform wecom \
    --wecom-corp-id YOUR_CORP_ID \
    --wecom-agent-id YOUR_AGENT_ID \
    --wecom-secret YOUR_SECRET \
    --wecom-token YOUR_TOKEN \
    --wecom-aes-key YOUR_AES_KEY

After verification succeeds, use 'lingti-bot relay' to start processing messages.`,
	Run: runVerify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVar(&verifyPlatform, "platform", "", "Platform: wecom (required)")
	verifyCmd.Flags().StringVar(&verifyServerURL, "server", "", "WebSocket URL (default: wss://bot.lingti.com/ws)")
	verifyCmd.Flags().IntVar(&verifyTimeout, "timeout", 300, "Timeout in seconds (default: 300)")

	// WeCom credentials
	verifyCmd.Flags().StringVar(&verifyWeComCorpID, "wecom-corp-id", "", "WeCom Corp ID (or WECOM_CORP_ID env)")
	verifyCmd.Flags().StringVar(&verifyWeComAgentID, "wecom-agent-id", "", "WeCom Agent ID (or WECOM_AGENT_ID env)")
	verifyCmd.Flags().StringVar(&verifyWeComSecret, "wecom-secret", "", "WeCom Secret (or WECOM_SECRET env)")
	verifyCmd.Flags().StringVar(&verifyWeComToken, "wecom-token", "", "WeCom Callback Token (or WECOM_TOKEN env)")
	verifyCmd.Flags().StringVar(&verifyWeComAESKey, "wecom-aes-key", "", "WeCom Encoding AES Key (or WECOM_AES_KEY env)")
}

func runVerify(cmd *cobra.Command, args []string) {
	// Get values from flags or environment
	if verifyPlatform == "" {
		verifyPlatform = os.Getenv("RELAY_PLATFORM")
	}
	if verifyServerURL == "" {
		verifyServerURL = os.Getenv("RELAY_SERVER_URL")
	}

	// Get WeCom credentials from flags or environment
	if verifyWeComCorpID == "" {
		verifyWeComCorpID = os.Getenv("WECOM_CORP_ID")
	}
	if verifyWeComAgentID == "" {
		verifyWeComAgentID = os.Getenv("WECOM_AGENT_ID")
	}
	if verifyWeComSecret == "" {
		verifyWeComSecret = os.Getenv("WECOM_SECRET")
	}
	if verifyWeComToken == "" {
		verifyWeComToken = os.Getenv("WECOM_TOKEN")
	}
	if verifyWeComAESKey == "" {
		verifyWeComAESKey = os.Getenv("WECOM_AES_KEY")
	}

	// Validate platform
	if verifyPlatform == "" {
		fmt.Fprintln(os.Stderr, "Error: --platform is required (currently supported: wecom)")
		os.Exit(1)
	}
	if verifyPlatform != "wecom" {
		fmt.Fprintln(os.Stderr, "Error: only 'wecom' platform is currently supported for verification")
		os.Exit(1)
	}

	// Validate credentials based on platform
	switch verifyPlatform {
	case "wecom":
		missing := []string{}
		if verifyWeComCorpID == "" {
			missing = append(missing, "--wecom-corp-id")
		}
		if verifyWeComAgentID == "" {
			missing = append(missing, "--wecom-agent-id")
		}
		if verifyWeComSecret == "" {
			missing = append(missing, "--wecom-secret")
		}
		if verifyWeComToken == "" {
			missing = append(missing, "--wecom-token")
		}
		if verifyWeComAESKey == "" {
			missing = append(missing, "--wecom-aes-key")
		}
		if len(missing) > 0 {
			fmt.Fprintf(os.Stderr, "Error: WeCom credentials required: %v\n", missing)
			os.Exit(1)
		}
	}

	// Generate a temporary user ID for verification
	verifyUserID := fmt.Sprintf("verify-%s-%d", verifyPlatform, time.Now().Unix())

	// Create verify-only relay connection
	verifyPlatformInstance, err := relay.New(relay.Config{
		UserID:       verifyUserID,
		Platform:     verifyPlatform,
		ServerURL:    verifyServerURL,
		AIProvider:   "verify", // Special marker for verification mode
		AIModel:      "verify",
		WeComCorpID:  verifyWeComCorpID,
		WeComAgentID: verifyWeComAgentID,
		WeComSecret:  verifyWeComSecret,
		WeComToken:   verifyWeComToken,
		WeComAESKey:  verifyWeComAESKey,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating verify connection: %v\n", err)
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(verifyTimeout)*time.Second)
	defer cancel()

	// Start the platform (connects to relay server)
	if err := verifyPlatformInstance.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to relay server: %v\n", err)
		os.Exit(1)
	}

	log.Println("")
	log.Println("=== Callback URL Verification Mode ===")
	log.Println("")
	log.Printf("  Platform: %s", verifyPlatform)
	log.Printf("  Corp ID:  %s", verifyWeComCorpID)
	log.Println("")
	log.Println("Your credentials have been sent to the cloud relay server.")
	log.Println("")
	log.Println("Now go to your platform's admin console and configure the callback URL:")
	log.Println("")
	switch verifyPlatform {
	case "wecom":
		log.Println("    https://bot.lingti.com/wecom")
	}
	log.Println("")
	log.Println("When you save the configuration, the platform will send a")
	log.Println("verification request which will be handled automatically.")
	log.Println("")
	log.Println("Press Ctrl+C to exit after verification succeeds.")
	log.Println("")

	// Wait for shutdown signal or timeout
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Println("\nShutting down...")
	case <-ctx.Done():
		log.Println("\nTimeout reached. If verification hasn't completed, run again.")
	}

	verifyPlatformInstance.Stop()

	log.Println("")
	log.Println("Next step: Use 'lingti-bot relay' to start processing messages:")
	log.Println("")
	log.Printf("  lingti-bot relay --user-id YOUR_ID --platform %s \\\n", verifyPlatform)
	log.Println("    --provider deepseek --api-key YOUR_API_KEY \\")
	log.Printf("    --wecom-corp-id %s ...\n", verifyWeComCorpID)
}
