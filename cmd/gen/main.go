package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"zee-ubot/internal/config"

	"github.com/gotd/td/session"
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

func updateEnvFile(sessionString string) error {
	content, err := os.ReadFile(".env")
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(".env", []byte("SESSION_STRING="+sessionString+"\n"), 0644)
		}
		return err
	}

	lines := strings.Split(string(content), "\n")
	found := false
	var newLines []string

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "SESSION_STRING=") {
			newLines = append(newLines, "SESSION_STRING="+sessionString)
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}

	if !found {
		newLines = append(newLines, "SESSION_STRING="+sessionString)
	}

	output := strings.Join(newLines, "\n")
	return os.WriteFile(".env", []byte(output), 0644)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config Error: %v\n", err)
		return
	}

	for {
		var sessionData []byte
		customStorage := &CustomStorage{
			OnStore: func(data []byte) {
				sessionData = data
			},
		}

		client := telegram.NewClient(cfg.AppID, cfg.AppHash, telegram.Options{
			SessionStorage: customStorage,
		})

		flow := auth.NewFlow(
			TermAuth{
				PhoneFn: prompt("Enter Phone Number (e.g. +62...): "),
				CodeFn:  prompt("Enter Code: "),
				PassFn:  prompt("Enter 2FA Password: "),
			},
			auth.SendCodeOptions{},
		)

		fmt.Println("\n🚀 Starting Login Helper...")

		err = client.Run(context.Background(), func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				return err
			}

			user, err := client.Self(ctx)
			if err != nil {
				return err
			}

			fmt.Printf("\n✅ Successfully Logged In as: %s (@%s)\n\n", user.FirstName, user.Username)
			return nil
		})

		if err != nil {
			fmt.Printf("\n❌ Login Failed: %v\n", err)

			fmt.Print("🔄 Try again? (Y/n): ")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(strings.ToLower(text))
			if text == "n" || text == "no" {
				return
			}
			continue
		}

		sessionString := base64.StdEncoding.EncodeToString(sessionData)
		fmt.Println("🔑 Generated Session String.")

		if err := updateEnvFile(sessionString); err != nil {
			fmt.Printf("⚠️ Failed to update .env automatically: %v\n", err)
			fmt.Println("Please manually add this to your .env file:")
			fmt.Println("SESSION_STRING=" + sessionString)
		} else {
			fmt.Println("✅ Updated .env file automatically!")
		}

		break
	}
}

type CustomStorage struct {
	OnStore func(data []byte)
}

func (s *CustomStorage) LoadSession(ctx context.Context) ([]byte, error) {
	return nil, session.ErrNotFound
}

func (s *CustomStorage) StoreSession(ctx context.Context, data []byte) error {
	if s.OnStore != nil {
		s.OnStore(data)
	}
	return nil
}
