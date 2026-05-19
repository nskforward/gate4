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
	// commandInfo("cert list", "show connected certs", "gate4 cert list")
	commandInfo("account list", "show accounts list", "gate4 account list")
	commandInfo("account add", "add a new account", "gate4 account add")
	commandInfo("account del <key>", "delete a account", "gate4 account del <key>")
	commandInfo("quotes <symbol> <key>", "subscribe for the quote stream", "gate4 quotes <symbol> <key>")
	commandInfo("position <symbol> <key>", "show the symbol position", "gate4 position <symbol> <key>")
	commandInfo("positions <key>", "show the list of all positions", "gate4 positions <key>")
	commandInfo("schedule <symbol> <key>", "show the symbol trading schedule", "gate4 schedule <symbol> <key>")
	commandInfo("my-trades <key>", "subscribe for the account trades", "gate4 my-trades <key>")
	commandInfo("asset <symbol> <key>", "show the asset general inforamtion", "gate4 asset <symbol> <key>")
	return nil
}
