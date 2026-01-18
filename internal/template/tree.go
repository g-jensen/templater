package template

import (
	"strings"
)

type treeNode struct {
	name     string
	children []*treeNode
}

func RenderTree(features []string) string {
	if len(features) == 0 {
		return ""
	}

	root := &treeNode{}
	for _, path := range features {
		parts := strings.Split(path, "/")
		insertPath(root, parts)
	}

	var sb strings.Builder
	renderNode(&sb, root.children, "")
	return sb.String()
}

func insertPath(node *treeNode, parts []string) {
	if len(parts) == 0 {
		return
	}

	name := parts[0]
	var child *treeNode
	for _, c := range node.children {
		if c.name == name {
			child = c
			break
		}
	}

	if child == nil {
		child = &treeNode{name: name}
		node.children = append(node.children, child)
	}

	insertPath(child, parts[1:])
}

func renderNode(sb *strings.Builder, nodes []*treeNode, prefix string) {
	for i, node := range nodes {
		isLast := i == len(nodes)-1

		if isLast {
			sb.WriteString(prefix + "└── " + node.name + "\n")
		} else {
			sb.WriteString(prefix + "├── " + node.name + "\n")
		}

		var childPrefix string
		if isLast {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}
		renderNode(sb, node.children, childPrefix)
	}
}
