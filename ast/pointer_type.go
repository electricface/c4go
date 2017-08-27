package ast

type PointerType struct {
	Addr     Address
	Type     string
	Children []Node
}

func parsePointerType(line string) *PointerType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &PointerType{
		Addr:     ParseAddress(groups["address"]),
		Type:     groups["type"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *PointerType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
