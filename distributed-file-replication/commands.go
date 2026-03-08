package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartCommandListener(n *Node) {

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print(ColorCyan + ">> " + ColorReset)

		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "" {
			continue
		}

		switch cmd {

		case "election":
			n.StartElection()

		case "cs":
			go n.RequestCS()

		case "enter":
			n.EnterCS()

		case "exitcs":
			n.ExitCS()

		case "replicate":
			n.ReplicateFile("report.txt", "distributed file")

		case "snapshot":
			n.StartSnapshot()

		case "files":
			n.ListFiles()

		case "status":
			PrintStatus(n)

		case "help":
			PrintHelp()

		default:
			LogWarn("Unknown command: %s. Type 'help' for a list of commands.", cmd)
		}
	}
}

func PrintHelp() {
	fmt.Println(ColorBlue + "---- AVAILABLE COMMANDS ----" + ColorReset)
	fmt.Println("  election   - Start leader election (Chang-Roberts)")
	fmt.Println("  cs         - Request Critical Section (Ricart-Agrawala)")
	fmt.Println("  enter      - Check if inside critical section")
	fmt.Println("  exitcs     - Exit Critical Section")
	fmt.Println("  replicate  - Start file replication")
	fmt.Println("  snapshot   - Start snapshot algorithm")
	fmt.Println("  files      - List replicated files")
	fmt.Println("  status     - Show current node status")
	fmt.Println("  help       - Show this menu")
	fmt.Println(ColorBlue + "----------------------------" + ColorReset)
}
