package raft

import "strings"

func ParseNodes(NodeList string) []string {
	return strings.Split(NodeList, ",")
}
