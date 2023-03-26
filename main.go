package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shokHorizon/jsonRunner/structs"
)

func main() {
	if len(os.Args) != 2 {
		panic("Not enough arguments")
	}
	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	nodes := make(map[string]*structs.Action)
	conditions := make(map[string]struct{})
	tasksQueue := make(chan *structs.Action, 10)
	cfg := structs.Config{}
	json.Unmarshal(bytes, &cfg)
	for i, node := range cfg.Actions {
		nodes[node.Name] = &cfg.Actions[i]
	}
	for _, node := range cfg.Conditions {
		nodes[node.Name] = &node
		conditions[node.Name] = struct{}{}
	}
	for _, node := range nodes {
		fmt.Println(node.Name)
		for _, nextNodeName := range node.Next {
			if nextNode, ok := nodes[nextNodeName]; ok {
				nextNode.PrevNodes = append(nextNode.PrevNodes, node)
				node.NextNodes = append(node.NextNodes, nextNode)
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
			err := node.Exec(nodes)
			if err != nil {
				panic(err)
			}
			if _, ok := conditions[node.Name]; !ok {
				for _, nextNode := range node.NextNodes {
					tasksQueue <- nextNode
				}
			} else if node.Result != "" {
				if nextNode, ok := nodes[node.Result]; ok {
					tasksQueue <- nextNode
				} else {
					panic("block initialization problem found!")
				}
			}
		default:
			done = !done
		}
	}
	cfgOut := structs.Config{}
	for _, node := range nodes {
		if _, ok := conditions[node.Name]; ok {
			cfgOut.Conditions = append(cfg.Conditions, *node)
		} else {
			cfgOut.Actions = append(cfg.Actions, *node)
		}
	}
	logs, err := json.MarshalIndent(cfgOut, "", "    ")
	if err != nil {
		panic(err)
	}
	os.WriteFile(os.Args[1]+".log", logs, 0777)

	fmt.Println("Queue is empty")
}
