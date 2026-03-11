package main

import (
	"fmt"
	"strconv"
)

// ANSI color codes for terminal
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
)

func LogInfo(format string, a ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", a...)
}

func LogWarn(format string, a ...interface{}) {
	fmt.Printf(ColorYellow+"[WARN] "+format+ColorReset+"\n", a...)
}

func LogElection(format string, a ...interface{}) {
	fmt.Printf(ColorCyan+"[ELECTION] "+format+ColorReset+"\n", a...)
}

func LogCS(format string, a ...interface{}) {
	fmt.Printf(ColorGreen+"[MUTEX] "+format+ColorReset+"\n", a...)
}

func LogSuccess(format string, a ...interface{}) {
	fmt.Printf(ColorMagenta+"[SUCCESS] "+format+ColorReset+"\n", a...)
}

func LogDebug(format string, a ...interface{}) {
	fmt.Printf(ColorWhite+"[DEBUG] "+format+ColorReset+"\n", a...)
}

func LogComm(format string, a ...interface{}) {
	fmt.Printf(ColorBlue+"[COMM] "+format+ColorReset+"\n", a...)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func PrintStatus(n *Node) {
	n.mu.RLock()
	id := n.ID
	leader := n.LeaderID
	state := n.State
	clock := n.LogicalClock
	files := n.Files
	n.mu.RUnlock()

	fmt.Println(ColorBlue + "---- NODE STATUS ----" + ColorReset)
	fmt.Println("Node ID:", id)
	fmt.Println("Leader:", leader)
	fmt.Println("State:", state)
	fmt.Println("Clock:", clock)
	fmt.Println("Files:", files)
	fmt.Println(ColorBlue + "---------------------" + ColorReset)
}

func (n *Node) StartHeartbeat() {
	LogInfo("Heartbeat service started for Node %d", n.ID)
	// We could add more periodic status logs here if desired
}
