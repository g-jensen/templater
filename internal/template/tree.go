package template

import (
	"strings"
)

type treeNode struct {
	name      string
	isFeature bool
	children  []*treeNode
}

func (n *treeNode) findChild(name string) *treeNode {
	for _, c := range n.children {
		if c.name == name {
			return c
		}
	}
	return nil
}

func RenderTree(features []string) string {
	if len(features) == 0 {
		return ""
	}

	featureSet := toSet(features)
	root := buildTree(features, featureSet)
	collapseNonFeatures(root, "")

	var sb strings.Builder
	renderChildren(&sb, root.children, "")
	return sb.String()
}

func toSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		set[item] = true
	}
	return set
}

func buildTree(features []string, featureSet map[string]bool) *treeNode {
	root := &treeNode{}
	for _, path := range features {
		parts := strings.Split(path, "/")
		node := root
		for i, part := range parts {
			fullPath := strings.Join(parts[:i+1], "/")
			child := node.findChild(part)
			if child == nil {
				child = &treeNode{name: part, isFeature: featureSet[fullPath]}
				node.children = append(node.children, child)
			}
			node = child
		}
	}
	return root
}

func collapseNonFeatures(node *treeNode, prefix string) {
	var collapsed []*treeNode
	for _, child := range node.children {
		childPrefix := joinPath(prefix, child.name)

		if !child.isFeature && len(child.children) > 0 {
			collapseNonFeatures(child, childPrefix)
			collapsed = append(collapsed, child.children...)
		} else {
			if prefix != "" {
				child.name = childPrefix
			}
			collapseNonFeatures(child, "")
			collapsed = append(collapsed, child)
		}
	}
	node.children = collapsed
}

func joinPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "/" + name
}

func renderChildren(sb *strings.Builder, nodes []*treeNode, prefix string) {
	for i, node := range nodes {
		isLast := i == len(nodes)-1
		connector, childPrefix := linePrefixes(prefix, isLast)

		sb.WriteString(connector + node.name + "\n")
		renderChildren(sb, node.children, childPrefix)
	}
}

func linePrefixes(prefix string, isLast bool) (connector, childPrefix string) {
	if isLast {
		return prefix + "└── ", prefix + "    "
	}
	return prefix + "├── ", prefix + "│   "
}
