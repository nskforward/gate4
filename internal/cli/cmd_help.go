package cli

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) cmdHelp(context.Context) error {
	maxCommandLen := 32
	commandInfo := func(command, description, example string) {
		fmt.Printf(" %s %s %s\n %s  $ %s\n\n", command, strings.Repeat(" ", maxCommandLen-len(command)), description, strings.Repeat(" ", maxCommandLen), example)
	}
	printHeader("Help")
	commandInfo("help", "show help menu", "gate4 help")
	commandInfo("cert init-ca", "generate ssl CA key and cert", "gate4 cert init-ca")
	commandInfo("cert issue", "generate and sign the ssl certicate", "gate4 cert issue")
	commandInfo("cert list-active", "show connected certs", "gate4 cert list-active")
	commandInfo("account list", "show accounts list", "gate4 account list")
	commandInfo("account add", "add a new account", "gate4 account add")
	commandInfo("account del <key>", "delete a account", "gate4 account del <key>")
	commandInfo("quotes <symbol> <key>", "subscribe for quote stream", "gate4 quotes <symbol> <key>")
	commandInfo("positions <key>", "show list of account positions", "gate4 positions <key>")
	commandInfo("schedule <symbol> <key>", "show the symbol trading schedule", "gate4 schedule <symbol> <key>")
	commandInfo("session <symbol> <key>", "show the current symbol schedule session", "gate4 session <symbol> <key>")
	return nil
}
