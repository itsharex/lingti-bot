package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pltanton/lingti-bot/internal/agent"
	"github.com/pltanton/lingti-bot/internal/platforms/discord"
	"github.com/pltanton/lingti-bot/internal/platforms/feishu"
	"github.com/pltanton/lingti-bot/internal/platforms/slack"
	"github.com/pltanton/lingti-bot/internal/platforms/telegram"
	"github.com/pltanton/lingti-bot/internal/router"
	"github.com/spf13/cobra"
)

var (
	slackBotToken   string
	slackAppToken   string
	feishuAppID     string
	feishuAppSecret string
	telegramToken   string
	discordToken    string
	aiProvider      string
	aiAPIKey        string
	aiBaseURL       string
	aiModel         string
)

var routerCmd = &cobra.Command{
	Use:   "router",
	Short: "Start the message router",
	Long: `Start the message router to receive messages from various platforms
(Slack, Telegram, Discord, Feishu) and respond using AI.

Supported platforms:
  - Slack: SLACK_BOT_TOKEN + SLACK_APP_TOKEN
  - Telegram: TELEGRAM_BOT_TOKEN
  - Discord: DISCORD_BOT_TOKEN
  - Feishu: FEISHU_APP_ID + FEISHU_APP_SECRET

Required environment variables or flags:
  - AI_PROVIDER: AI provider (claude, deepseek, kimi) default: claude
  - AI_API_KEY: API Key for the AI provider
  - AI_BASE_URL: Custom API base URL (optional)
  - AI_MODEL: Model name (optional)`,
	Run: runRouter,
}

func init() {
	rootCmd.AddCommand(routerCmd)

	routerCmd.Flags().StringVar(&slackBotToken, "slack-bot-token", "", "Slack Bot Token (or SLACK_BOT_TOKEN env)")
	routerCmd.Flags().StringVar(&slackAppToken, "slack-app-token", "", "Slack App Token (or SLACK_APP_TOKEN env)")
	routerCmd.Flags().StringVar(&feishuAppID, "feishu-app-id", "", "Feishu App ID (or FEISHU_APP_ID env)")
	routerCmd.Flags().StringVar(&feishuAppSecret, "feishu-app-secret", "", "Feishu App Secret (or FEISHU_APP_SECRET env)")
	routerCmd.Flags().StringVar(&telegramToken, "telegram-token", "", "Telegram Bot Token (or TELEGRAM_BOT_TOKEN env)")
	routerCmd.Flags().StringVar(&discordToken, "discord-token", "", "Discord Bot Token (or DISCORD_BOT_TOKEN env)")
	routerCmd.Flags().StringVar(&aiProvider, "provider", "", "AI provider: claude or deepseek (or AI_PROVIDER env)")
	routerCmd.Flags().StringVar(&aiAPIKey, "api-key", "", "AI API Key (or AI_API_KEY env)")
	routerCmd.Flags().StringVar(&aiBaseURL, "base-url", "", "Custom API base URL (or AI_BASE_URL env)")
	routerCmd.Flags().StringVar(&aiModel, "model", "", "Model name (or AI_MODEL env)")
}

func runRouter(cmd *cobra.Command, args []string) {
	// Get tokens from flags or environment
	if slackBotToken == "" {
		slackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	}
	if slackAppToken == "" {
		slackAppToken = os.Getenv("SLACK_APP_TOKEN")
	}
	if feishuAppID == "" {
		feishuAppID = os.Getenv("FEISHU_APP_ID")
	}
	if feishuAppSecret == "" {
		feishuAppSecret = os.Getenv("FEISHU_APP_SECRET")
	}
	if telegramToken == "" {
		telegramToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	if discordToken == "" {
		discordToken = os.Getenv("DISCORD_BOT_TOKEN")
	}
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiAPIKey == "" {
		aiAPIKey = os.Getenv("AI_API_KEY")
		// Fallback to legacy env var
		if aiAPIKey == "" {
			aiAPIKey = os.Getenv("ANTHROPIC_API_KEY")
		}
	}
	if aiBaseURL == "" {
		aiBaseURL = os.Getenv("AI_BASE_URL")
		if aiBaseURL == "" {
			aiBaseURL = os.Getenv("ANTHROPIC_BASE_URL")
		}
	}
	if aiModel == "" {
		aiModel = os.Getenv("AI_MODEL")
		if aiModel == "" {
			aiModel = os.Getenv("ANTHROPIC_MODEL")
		}
	}

	// Validate required tokens
	if aiAPIKey == "" {
		fmt.Fprintln(os.Stderr, "Error: AI_API_KEY is required")
		os.Exit(1)
	}

	// Create the AI agent
	aiAgent, err := agent.New(agent.Config{
		Provider: aiProvider,
		APIKey:   aiAPIKey,
		BaseURL:  aiBaseURL,
		Model:    aiModel,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating agent: %v\n", err)
		os.Exit(1)
	}

	// Create the router with the agent as message handler
	r := router.New(aiAgent.HandleMessage)

	// Register Slack if tokens are provided
	if slackBotToken != "" && slackAppToken != "" {
		slackPlatform, err := slack.New(slack.Config{
			BotToken: slackBotToken,
			AppToken: slackAppToken,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Slack platform: %v\n", err)
			os.Exit(1)
		}
		r.Register(slackPlatform)
	} else {
		log.Println("Slack tokens not provided, skipping Slack integration")
	}

	// Register Feishu if tokens are provided
	if feishuAppID != "" && feishuAppSecret != "" {
		feishuPlatform, err := feishu.New(feishu.Config{
			AppID:     feishuAppID,
			AppSecret: feishuAppSecret,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Feishu platform: %v\n", err)
			os.Exit(1)
		}
		r.Register(feishuPlatform)
	} else {
		log.Println("Feishu tokens not provided, skipping Feishu integration")
	}

	// Register Telegram if token is provided
	if telegramToken != "" {
		telegramPlatform, err := telegram.New(telegram.Config{
			Token: telegramToken,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Telegram platform: %v\n", err)
			os.Exit(1)
		}
		r.Register(telegramPlatform)
	} else {
		log.Println("Telegram token not provided, skipping Telegram integration")
	}

	// Register Discord if token is provided
	if discordToken != "" {
		discordPlatform, err := discord.New(discord.Config{
			Token: discordToken,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Discord platform: %v\n", err)
			os.Exit(1)
		}
		r.Register(discordPlatform)
	} else {
		log.Println("Discord token not provided, skipping Discord integration")
	}

	// Start the router
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := r.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting router: %v\n", err)
		os.Exit(1)
	}

	providerName := aiProvider
	if providerName == "" {
		providerName = "claude"
	}
	modelName := aiModel
	if modelName == "" {
		switch providerName {
		case "deepseek":
			modelName = "deepseek-chat"
		case "kimi", "moonshot":
			modelName = "moonshot-v1-8k"
		default:
			modelName = "claude-sonnet-4-20250514"
		}
	}
	log.Printf("Router started. AI Provider: %s, Model: %s", providerName, modelName)
	log.Println("Press Ctrl+C to stop.")

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
	r.Stop()
}
