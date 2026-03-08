package main

type RequestMessage struct {
	NodeID    int
	Timestamp int
}

// RPCHandler.RequestCS is called by a remote node that wants the critical section.
// We reply immediately unless we have higher priority, in which case we defer.
func (h *RPCHandler) RequestCS(req RequestMessage, reply *bool) error {

	n := h.node

	n.mu.Lock()

	sendNow := false

	if n.State == "RELEASED" {
		sendNow = true
	} else if n.State == "WANTED" {
		// Lower timestamp wins; tie-break by lower node ID
		if req.Timestamp < n.RequestTimestamp ||
			(req.Timestamp == n.RequestTimestamp && req.NodeID < n.ID) {
			sendNow = true
		} else {
			n.DeferredRequests = append(n.DeferredRequests, req.NodeID)
			LogCS("Node %d deferred reply to Node %d", n.ID, req.NodeID)
		}
	} else {
		// HELD — always defer
		n.DeferredRequests = append(n.DeferredRequests, req.NodeID)
		LogCS("Node %d deferred reply to Node %d", n.ID, req.NodeID)
	}

	n.mu.Unlock()

	if sendNow {
		LogCS("Node %d granted reply to Node %d", n.ID, req.NodeID)
	}

	*reply = sendNow
	return nil
}

// RPCHandler.GrantCS is called by a peer that previously deferred our CS request
// and is now releasing the critical section.
func (h *RPCHandler) GrantCS(senderID int, reply *bool) error {

	n := h.node

	n.mu.Lock()
	n.ReplyCount++
	needed := len(n.Peers) - 1
	ready := n.ReplyCount >= needed && n.State == "WANTED"
	ch := n.csReadyCh
	n.mu.Unlock()

	LogCS("Node %d received deferred grant from Node %d", n.ID, senderID)

	if ready && ch != nil {
		select {
		case ch <- struct{}{}:
		default:
		}
	}

	*reply = true
	return nil
}

func (n *Node) RequestCS() {

	n.mu.Lock()
	n.LogicalClock++
	n.State = "WANTED"
	n.RequestTimestamp = n.LogicalClock
	n.ReplyCount = 0
	n.csReadyCh = make(chan struct{}, 1)
	n.mu.Unlock()

	LogCS("Node %d requesting critical section", n.ID)

	needed := len(n.Peers) - 1

	for id := range n.Peers {

		if id == n.ID {
			continue
		}

		LogCS("Sending REQUEST to Node %d", id)

		var reply bool
		req := RequestMessage{NodeID: n.ID, Timestamp: n.RequestTimestamp}

		err := n.CallPeer(id, "RPCHandler.RequestCS", req, &reply)

		if err != nil {
			// Peer is down — treat as implicit grant
			LogWarn("Node %d unreachable, counting as grant", id)
			n.mu.Lock()
			n.ReplyCount++
			n.mu.Unlock()
		} else if reply {
			LogSuccess("Reply granted by Node %d", id)
			n.mu.Lock()
			n.ReplyCount++
			n.mu.Unlock()
		} else {
			LogWarn("Reply deferred by Node %d", id)
		}
	}

	// Block until all deferred replies arrive via RPCHandler.GrantCS
	n.mu.Lock()
	alreadyReady := n.ReplyCount >= needed
	n.mu.Unlock()

	if !alreadyReady {
		LogInfo("Node %d waiting for deferred replies...", n.ID)
		<-n.csReadyCh
	}

	n.mu.Lock()
	n.State = "HELD"
	n.mu.Unlock()

	LogSuccess("Node %d entered critical section", n.ID)
}

func (n *Node) EnterCS() {

	n.mu.RLock()
	state := n.State
	n.mu.RUnlock()

	if state == "HELD" {
		LogInfo("Node %d is already in the critical section", n.ID)
	} else {
		LogWarn("Node %d is not in the critical section (state: %s)", n.ID, state)
	}
}

func (n *Node) ExitCS() {

	n.mu.Lock()

	if n.State != "HELD" {
		n.mu.Unlock()
		LogWarn("Cannot exit critical section: node is not inside it")
		return
	}

	n.State = "RELEASED"
	n.RequestTimestamp = 0
	n.ReplyCount = 0
	deferred := n.DeferredRequests
	n.DeferredRequests = nil

	n.mu.Unlock()

	LogInfo("Node %d exited critical section", n.ID)

	// Send queued-up grants now that we are released
	for _, peerID := range deferred {
		var reply bool
		err := n.CallPeer(peerID, "RPCHandler.GrantCS", n.ID, &reply)
		if err != nil {
			LogWarn("Failed to send grant to Node %d: %v", peerID, err)
		} else {
			LogInfo("Node %d sent grant to Node %d", n.ID, peerID)
		}
	}
}
