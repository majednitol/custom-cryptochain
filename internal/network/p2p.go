package network

import (
	"encoding/json"
	"net"
)

type Node struct {
	Peers []string
}

func (n *Node) Broadcast(msg Message) {
	for _, peer := range n.Peers {
		conn, _ := net.Dial("tcp", peer)
		json.NewEncoder(conn).Encode(msg)
		conn.Close()
	}
}
