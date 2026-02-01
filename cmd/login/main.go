package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"zee-ubot/internal/config"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

func prompt(msg string) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		fmt.Print(msg)
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(text), nil
	}
}

type TermAuth struct {
	PhoneFn func(context.Context) (string, error)
	CodeFn  func(context.Context) (string, error)
	PassFn  func(context.Context) (string, error)
}

func (t TermAuth) Phone(ctx context.Context) (string, error)                    { return t.PhoneFn(ctx) }
func (t TermAuth) Password(ctx context.Context) (string, error)                 { return t.PassFn(ctx) }
func (t TermAuth) Code(ctx context.Context, _ *tg.AuthSentCode) (string, error) { return t.CodeFn(ctx) }
func (t TermAuth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("signup not supported")
}
func (t TermAuth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return nil
}

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("      ZEE-UBOT LOGIN GENERATOR          ")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	sessionDir := "session"
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		fmt.Printf("Failed to create session dir: %v\n", err)
		return
	}
	sessionPath := filepath.Join(sessionDir, "session.json")

	_ = os.Remove(sessionPath)

	storage := &telegram.FileSessionStorage{Path: sessionPath}

	client := telegram.NewClient(cfg.AppID, cfg.AppHash, telegram.Options{
		SessionStorage: storage,
	})

	flow := auth.NewFlow(
		TermAuth{
			PhoneFn: prompt("Enter Phone Number: "),
			CodeFn:  prompt("Enter Code: "),
			PassFn:  prompt("Enter Password (if known): "),
		},
		auth.SendCodeOptions{},
	)

	ctx := context.Background()

	err = client.Run(ctx, func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		self, err := client.Self(ctx)
		if err != nil {
			return err
		}

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("✅ Login Successful!\n")
		fmt.Printf("👤 User: %s (@%s)\n", self.FirstName, self.Username)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(sessionPath); err == nil {
		fmt.Println("✅ Session file generated successfully.")
	} else {
		fmt.Println("❌ Session file was not generated.")
		os.Exit(1)
	}
}
