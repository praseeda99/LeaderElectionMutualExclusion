package main

func (n *Node) StartElection() {
	LogElection("Node %d starting election", n.ID)
	var reply bool
	n.CallNextNeighbor("RPCHandler.Election", n.ID, &reply)
}

func (h *RPCHandler) Election(candidateID int, reply *bool) error {
	n := h.node
	*reply = true

	if candidateID > n.ID {
		LogElection("Node %d forwarding Election message (ID: %d)", n.ID, candidateID)
		go n.CallNextNeighbor("RPCHandler.Election", candidateID, reply)
	} else if candidateID < n.ID {
		LogElection("Node %d replacing Election ID %d with My ID %d", n.ID, candidateID, n.ID)
		go n.CallNextNeighbor("RPCHandler.Election", n.ID, reply)
	} else if candidateID == n.ID {
		LogSuccess("Node %d Election message returned! Becoming LEADER.", n.ID)
		n.BecomeLeader()
	}

	return nil
}

func (n *Node) BecomeLeader() {
	n.mu.Lock()
	alreadyLeader := n.LeaderID == n.ID
	n.LeaderID = n.ID
	n.mu.Unlock()

	if !alreadyLeader {
		LogSuccess("Node %d is now the COORDINATOR", n.ID)
	} else {
		LogElection("Node %d re-announcing leadership to ring", n.ID)
	}

	// Inform everyone in the ring
	var ack bool
	n.CallNextNeighbor("RPCHandler.Coordinator", n.ID, &ack)
}

func (h *RPCHandler) Coordinator(leaderID int, ack *bool) error {
	n := h.node
	*ack = true

	// If the message has returned to the leader who started it, stop propagation
	if n.ID == leaderID {
		LogSuccess("Leader announcement for Node %d complete", leaderID)
		return nil
	}

	n.mu.Lock()
	n.LeaderID = leaderID
	n.mu.Unlock()

	LogSuccess("Node %d recognized Node %d as the new LEADER", n.ID, leaderID)

	// Forward to next neighbor
	go n.CallNextNeighbor("RPCHandler.Coordinator", leaderID, ack)

	return nil
}
