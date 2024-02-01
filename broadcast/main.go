package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

func main() {
	n := maelstrom.NewNode()
	topology := make(map[string]interface{})
	broadcastMsgs := make([]float64, 0)

	n.Handle("topology", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		body["type"] = "topology_ok"
		topology = body["topology"].(map[string]interface{})
		delete(body, "topology")

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		receivedMsg := body["message"].(float64)
		broadcastMsgs = append(broadcastMsgs, receivedMsg)

		nodeId := n.ID()
		if neighbors, ok := topology[nodeId].([]string); ok {
			for _, neighbor := range neighbors {
				n.Send(neighbor, receivedMsg)
			}
		}
		body["type"] = "broadcast_ok"
		delete(body, "message")
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		body["type"] = "read_ok"
		body["messages"] = broadcastMsgs

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
