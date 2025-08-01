package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"heckel.io/ntfy/v2/server"
	"heckel.io/ntfy/v2/test"
	"testing"
)

func TestCLI_Tier_AddListChangeDelete(t *testing.T) {
	s, conf, port := newTestServerWithAuth(t)
	defer test.StopServer(t, s, port)

	app, _, stdout, _ := newTestApp()
	require.Nil(t, runTierCommand(app, conf, "add", "--name", "Pro", "--message-limit", "1234", "pro"))
	require.Contains(t, stdout.String(), "tier added\n\ntier pro (id: ti_")

	err := runTierCommand(app, conf, "add", "pro")
	require.NotNil(t, err)
	require.Equal(t, "tier pro already exists", err.Error())

	app, _, stdout, _ = newTestApp()
	require.Nil(t, runTierCommand(app, conf, "list"))
	require.Contains(t, stdout.String(), "tier pro (id: ti_")
	require.Contains(t, stdout.String(), "- Name: Pro")
	require.Contains(t, stdout.String(), "- Message limit: 1234")

	app, _, stdout, _ = newTestApp()
	require.Nil(t, runTierCommand(app, conf, "change",
		"--message-limit=999",
		"--message-expiry-duration=2d",
		"--email-limit=91",
		"--reservation-limit=98",
		"--attachment-file-size-limit=100m",
		"--attachment-expiry-duration=1d",
		"--attachment-total-size-limit=10G",
		"--attachment-bandwidth-limit=100G",
		"--stripe-monthly-price-id=price_991",
		"--stripe-yearly-price-id=price_992",
		"pro",
	))
	require.Contains(t, stdout.String(), "- Message limit: 999")
	require.Contains(t, stdout.String(), "- Message expiry duration: 48h")
	require.Contains(t, stdout.String(), "- Email limit: 91")
	require.Contains(t, stdout.String(), "- Reservation limit: 98")
	require.Contains(t, stdout.String(), "- Attachment file size limit: 100.0 MB")
	require.Contains(t, stdout.String(), "- Attachment expiry duration: 24h")
	require.Contains(t, stdout.String(), "- Attachment total size limit: 10.0 GB")
	require.Contains(t, stdout.String(), "- Stripe prices (monthly/yearly): price_991 / price_992")

	app, _, stdout, _ = newTestApp()
	require.Nil(t, runTierCommand(app, conf, "remove", "pro"))
	require.Contains(t, stdout.String(), "tier pro removed")
}

func runTierCommand(app *cli.App, conf *server.Config, args ...string) error {
	userArgs := []string{
		"ntfy",
		"--log-level=ERROR",
		"tier",
		"--config=" + conf.File, // Dummy config file to avoid lookups of real file
		"--auth-file=" + conf.AuthFile,
		"--auth-default-access=" + conf.AuthDefault.String(),
	}
	return app.Run(append(userArgs, args...))
}
