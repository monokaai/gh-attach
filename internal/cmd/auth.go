package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sudosubin/gh-attach/internal/app"
	"github.com/sudosubin/gh-attach/internal/cookies"
)

var captureBrowserSession = func(browser, profile, cookieStorePath string) (string, error) {
	return app.NewService(os.Stderr).CaptureBrowserSessionToken(context.Background(), "github.com", cookies.ResolveInput{
		Browser: browser, Profile: profile, CookieStorePath: cookieStorePath,
	}, false)
}
var storeEvidenceSession = app.StoreSystemEvidenceSession

func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{Use: "auth", Short: "Manage the dedicated headless evidence session"}
	cmd.AddCommand(NewCmdAuthCapture())
	return cmd
}

func NewCmdAuthCapture() *cobra.Command {
	var browser, profile, cookieStorePath string
	cmd := &cobra.Command{
		Use:   "capture",
		Short: "Capture a browser session into the dedicated macOS Keychain item",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := captureBrowserSession(browser, profile, cookieStorePath)
			if err != nil {
				return fmt.Errorf("capture browser session: %w", err)
			}
			if err := storeEvidenceSession(token); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Stored the headless evidence session in macOS Keychain.")
			return nil
		},
	}
	cmd.Flags().StringVar(&browser, "browser", "", "Browser to capture from ("+cookies.BrowserChoices()+")")
	cmd.Flags().StringVar(&profile, "profile", "", "Browser profile name")
	cmd.Flags().StringVar(&cookieStorePath, "cookie-store-path", "", "Cookie store file path")
	_ = cmd.MarkFlagRequired("browser")
	return cmd
}
