package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shokHorizon/jsonRunner/structs"
)

func main() {
	bytes, err := os.ReadFile("example.json")
	nodes := make(map[string]structs.Action)
	conditions := make(map[string]struct{})
	tasksQueue := make(chan structs.Action, 10)
	if err != nil {
		panic("File doesnt exist")
	}
	cfg := structs.Config{}
	json.Unmarshal(bytes, &cfg)
	for _, node := range cfg.Actions {
		nodes[node.Name] = node
	}
	for _, node := range cfg.Conditions {
		nodes[node.Name] = node
		conditions[node.Name] = struct{}{}
	}
	for _, node := range nodes {
		for _, nextNodeName := range node.Next {
			if nextNode, ok := nodes[nextNodeName]; ok {
				nextNode.PrevNodes = append(nextNode.PrevNodes, &node)
				node.NextNodes = append(node.NextNodes, &nextNode)
				nodes[nextNodeName] = nextNode
				nodes[node.Name] = node
			}
		}
	}
	for _, node := range nodes {
		if len(node.PrevNodes) == 0 {
			tasksQueue <- node
		}
	}

	var done bool

	for !done {
		select {
		case node := <-tasksQueue:
			fmt.Println("<: ", node.Name)
			err := node.Exec(nodes)
			fmt.Println(">: ", node.Name)
			if err != nil {
				panic(err)
			}
			if _, ok := conditions[node.Name]; !ok {
				for _, nextNodeName := range node.Next {
					if nextNode, ok := nodes[nextNodeName]; ok {
						fmt.Println("+: ", nextNode.Name)
						tasksQueue <- nextNode
					} else {
						panic("block initialization problem found!")
					}
				}
			} else if node.Result != "" {
				if nextNode, ok := nodes[node.Result]; ok {
					fmt.Println("+", nextNode)
					tasksQueue <- nextNode
				} else {
					panic("block initialization problem found!")
				}
			}
		default:
			done = !done
		}
	}
	fmt.Println("Queue is empty")
}
