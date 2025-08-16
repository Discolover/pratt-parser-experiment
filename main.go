package main

import (
	"fmt"
	"strings"
)

type Token int

const (
	EOI Token = iota
	Number
	Plus
	Minus
	Times
	Less
	Greater
)

type Node struct {
	v Token
	l *Node
	r *Node
}

func (t Token) String() string {
	return []string{"EOI", "Number", "Plus", "Minus", "Times", "Less", "Greater"}[t]
}

func (t Token) Character() string {
	return []string{"", "N", "+", "-", "*", "<", ">"}[t]
}

type Parser struct {
	stream []Token
}

func (p *Parser) getToken() (t Token) {
	if len(p.stream) <= 0 {
		t = EOI
	} else {
		t = p.stream[len(p.stream)-1]
		p.stream = p.stream[:len(p.stream)-1]
	}

	return
}

func isBinaryOperator(t Token) bool {
	return t == Plus || t == Minus || t == Times || t == Less || t == Greater
}

func (p *Parser) parseExpression() *Node {
	left := &Node{v: p.getToken()}

	if left.v != Number {
		panic("Expected `Number`")
	}

	next := p.getToken()
	if isBinaryOperator(next) {
		right := p.parseExpression()
		return &Node{v: next, l: left, r: right}
	}

	return left
}

func traverse(node *Node, depth int) {
	if node == nil {
		return
	}
	traverse(node.l, depth+1)
	traverse(node.r, depth+1)

	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(node.v)
}

func drawTree(root *Node) {
	type Element struct {
		node  *Node
		depth int
		index int
	}

	var queue []Element

	queue = append(queue, Element{root, 0, 0})

	maxDepth := -1

	index2node := map[int]*Node{}

	for len(queue) > 0 {
		element := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		nd := element.node
		depth := element.depth

		if depth > maxDepth {
			maxDepth = depth
		}

		index := element.index

		index2node[index] = nd

		nPerDepth := 1 << depth // 2^depth

		firstIndex := nPerDepth - 1
		lastIndex := (nPerDepth - 1) * 2

		nBefore := index - firstIndex
		nAfter := lastIndex - index

		if nd.r != nil {
			queue = append(queue, Element{nd.r, depth + 1, index + nBefore*2 + nAfter + 2})
		}

		if nd.l != nil {
			queue = append(queue, Element{nd.l, depth + 1, index + nBefore*2 + nAfter + 1})
		}
	}

	maxDepth++

	gaps := make([]int, maxDepth)

	gaps[len(gaps)-1] = 1
	for i := len(gaps) - 2; i >= 0; i-- {
		gaps[i] = gaps[i+1]*2 + 1
	}

	index := 0
	for i := 0; i < maxDepth; i++ {
		fmt.Print(strings.Repeat(" ", gaps[i]/2))
		n := 1 << i
		for j := 0; j < n; j++ {
			if node, ok := index2node[index]; ok {
				fmt.Print(node.v.Character())
			} else {
				fmt.Print(" ")
			}

			index++

			if j+1 < n {
				fmt.Print(strings.Repeat(" ", gaps[i]))
			}
		}
		fmt.Println()
	}

	// h = maxDepth + 1; w = (1 << maxDepth) * 2 - 1

}

func main() {
	p := Parser{stream: []Token{Number, Minus, Number, Plus, Number, Plus, Number, Minus, Number}}

	// for t := p.getToken(); t != EOI; t = p.getToken() {
	// 	fmt.Println(t)
	// }

	root := p.parseExpression()

	traverse(root, 0)
	fmt.Println()
	drawTree(root)
}
