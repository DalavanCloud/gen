package {{.Package}}

// this file generated, do not edit

{{.Header}}

import (
	"fmt"
	{{if .Trace}}"log"{{end}}
)

// $Rule is a rule of the grammar.
type $Rule struct {
	symbol  string
	pattern []string
	reduce  func(data []interface{}) interface{}
}

// Action is an entry in the action table.
// Encoding:
//   0: accept
//   n: shift n
//   -n: reduce n
//   (errors are not in the map)
type $Action int

// $ActionTable holds the parser's precomputed state.
// table[state][token] => action to take on token from state.
type $ActionTable []map[string]$Action

// $Parser manages the parsing process.
type $Parser struct {
	actions $ActionTable
	stack   []int
	data    []interface{}
}

// $NewParser constructs a new $Parser, ready for input.
func $NewParser() *$Parser {
	return &$Parser{
		actions: $Actions,
		stack:   []int{0},
		data:    []interface{}{},
	}
}

// Parse processes one token, returning true on a complete parse and
// false when more input is expected.
func (p *$Parser) Parse(tok *{{.TokenType}}) (bool, error) {
	for {
		{{if .Trace}}
		log.Println("")
		log.Printf("stack:%v, data:%v\n", p.stack, p.data)
		log.Printf("tok:%v\n", tok.ParseId())
		{{end}}
		action, ok := p.actions[p.stack[len(p.stack)-1]][tok.ParseId()]
		if !ok {
			log.Println(p.actions[p.stack[len(p.stack)-1]])
			return false, fmt.Errorf("unexpected token: %v", tok)
		}

		if action > 0 {
			// To shift, we consume the current token and put the next
			// state on the stack.
			nextState := int(action)
			{{if .Trace}}
			log.Printf("input %v => shift %#v\n", tok, nextState)
			{{end}}
			p.data = append(p.data, *tok)
			p.stack = append(p.stack, nextState)

			// Ready for another token.
			return false, nil

		} else if action <= 0 {
			// To reduce, we pop off the matching pattern from the stacks.
			rule := $Rules[-action]
			{{if .Trace}}
			log.Printf("input %v => reduce %s -> %s\n", tok, rule.pattern, rule.symbol)
			{{end}}
			popCount := len(rule.pattern)

			// Update the data stack via the reduce function if available.
			oldData := p.data[len(p.data)-popCount:]
			var newData interface{}
			if rule.reduce != nil {
				newData = rule.reduce(oldData)
			} else if popCount == 1 {
				newData = oldData[0]
			} else {
				s := make([]interface{}, popCount)
				copy(s, oldData)
				newData = s
			}
			p.data = p.data[0 : len(p.data)-popCount]
			p.data = append(p.data, newData)

			p.stack = p.stack[0 : len(p.stack)-popCount]

			if action == 0 {
				// Accept.
				return true, nil
			}

			// Advance to the next state.
			state := p.stack[len(p.stack)-1]
			action, ok = p.actions[state][rule.symbol]
			if !ok || action <= 0 {
				// TODO: better error here; can it actually happen?
				panic(fmt.Errorf("parse error near %s: bad next state", tok.Pos))
			}

			p.stack = append(p.stack, int(action))
		}
	}
}

