package main

import (
	"sort"
	"sync"
)

type Node struct {
	ID   int
	Port int

	Peers map[int]string

	LeaderID int
	Alive    bool

	mu sync.RWMutex

	// Ricart–Agrawala
	LogicalClock     int
	State            string
	RequestTimestamp int
	ReplyCount       int
	DeferredRequests []int
	csReadyCh        chan struct{}

	// Files
	Files map[string]int
}

func (n *Node) GetNextNeighbor() int {
	n.mu.RLock()
	defer n.mu.RUnlock()

	ids := []int{}
	for id := range n.Peers {
		ids = append(ids, id)
	}
	// Sort IDs to define ring order
	sort.Ints(ids)

	for i, id := range ids {
		if id == n.ID {
			nextIdx := (i + 1) % len(ids)
			return ids[nextIdx]
		}
	}
	return n.ID
}

func NewNode(id int, port int, peers map[int]string) *Node {

	if len(peers) == 0 {
		peers = map[int]string{
			1: "10.97.141.218:8001",
			2: "10.97.141.82:8002",
			3: "10.97.141.35:8003",
			4: "10.97.141.158:8004",
		}
	}

	return &Node{
		ID:   id,
		Port: port,

		Peers: peers,

		Alive: true,

		State: "RELEASED",

		Files: make(map[string]int),
	}
}
