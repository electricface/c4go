package ast

type WhileStmt struct {
	Addr     Address
	Position string
	Children []Node
}

func parseWhileStmt(line string) *WhileStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &WhileStmt{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *WhileStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
