package cli

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) cmdHelp(context.Context) error {
	maxCommandLen := 18
	commandInfo := func(command, description, example string) {
		fmt.Printf(" %s %s %s\n %s  $ %s\n\n", command, strings.Repeat(" ", maxCommandLen-len(command)), description, strings.Repeat(" ", maxCommandLen), example)
	}
	printHeader("Help")
	commandInfo("help", "show help menu", "gate4 help")
	commandInfo("cert init-ca", "generate ssl CA key and cert and store to files", "gate4 cert init-ca")
	commandInfo("cert issue", "generate and sign the robot ssl certicate", "gate4 cert issue")
	commandInfo("cert list-active", "show connected certs", "gate4 cert list-active")
	commandInfo("account list", "show accounts list", "gate4 account list")
	commandInfo("account add", "add a new account", "gate4 account add")
	commandInfo("account del <key>", "delete a account", "gate4 account del <key>")
	commandInfo("subscribe quotes <symbol> <account_key>", "subscribe for quote stream", "gate4 subscribe quotes <symbol> <account_key>")
	commandInfo("positions <account_key>", "show list of account positions", "gate4 positions <account_key>")
	commandInfo("schedule <symbol> <account_key>", "show the symbol trading schedule", "gate4 schedule <symbol> <account_key>")
	return nil
}
