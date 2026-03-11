package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

type RPCHandler struct {
	node *Node
}

func (n *Node) StartRPC() {

	handler := &RPCHandler{node: n}

	rpc.Register(handler)

	addr := ":" + IntToString(n.Port)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("RPC server running on", addr)

	rpc.Accept(listener)
}

func (n *Node) CallPeer(peerID int, method string, args interface{}, reply interface{}) error {

	addr := n.Peers[peerID]
	LogComm("Node %d calling %s on Node %d (%s)", n.ID, method, peerID, addr)

	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return err
	}

	client := rpc.NewClient(conn)
	// Make sure to close client, but if we time out, we should still let it be garbage collected or closed.
	// We'll close the connection to abort the in-flight request if it times out
	defer client.Close()

	call := client.Go(method, args, reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			LogWarn("RPC call %s to Node %d failed: %v", method, peerID, call.Error)
		} else {
			LogDebug("RPC call %s to Node %d returned successfully", method, peerID)
		}
		return call.Error
	case <-time.After(1 * time.Second):
		LogWarn("RPC call %s to Node %d timed out", method, peerID)
		return fmt.Errorf("rpc call timeout to %s", addr)
	}
}
