package extractor

import (
	"go/ast"
	"go/token"
)

// ConcurrencyExtractor extracts concurrency-related constructs
type ConcurrencyExtractor struct{}

// NewConcurrencyExtractor creates a new concurrency extractor
func NewConcurrencyExtractor() *ConcurrencyExtractor {
	return &ConcurrencyExtractor{}
}

// Extract extracts concurrency info from an AST node
func (e *ConcurrencyExtractor) Extract(node ast.Node, ctx *Context) ([]ExtractionResult, error) {
	var results []ExtractionResult

	switch n := node.(type) {
	case *ast.GoStmt:
		// Goroutine invocation
		results = append(results, &ConcurrencyInfo{
			Type:     "goroutine",
			Details:  map[string]interface{}{"call": exprToString(n.Call)},
			Location: getLocation(n, ctx),
		})

	case *ast.ChanType:
		// Channel type declaration
		direction := "bidirectional"
		if n.Dir == ast.SEND {
			direction = "send-only"
		} else if n.Dir == ast.RECV {
			direction = "receive-only"
		}

		results = append(results, &ConcurrencyInfo{
			Type:    "channel",
			Details: map[string]interface{}{
				"direction": direction,
				"elementType": exprToString(n.Value),
			},
			Location: getLocation(n, ctx),
		})

	case *ast.SelectStmt:
		// Select statement
		results = append(results, &ConcurrencyInfo{
			Type:     "select",
			Details:  map[string]interface{}{"branchCount": len(n.Body.List)},
			Location: getLocation(n, ctx),
		})

	case *ast.SendStmt:
		// Channel send operation
		results = append(results, &ConcurrencyInfo{
			Type:     "send",
			Details:  map[string]interface{}{
				"channel": exprToString(n.Chan),
				"value": exprToString(n.Value),
			},
			Location: getLocation(n, ctx),
		})

	case *ast.UnaryExpr:
		// Receive operation (<-ch)
		if n.Op == token.ARROW {
			results = append(results, &ConcurrencyInfo{
				Type:     "receive",
				Details:  map[string]interface{}{"channel": exprToString(n.X)},
				Location: getLocation(n, ctx),
			})
		}

	case *ast.DeferStmt:
		// Defer statement
		results = append(results, &ConcurrencyInfo{
			Type:     "defer",
			Details:  map[string]interface{}{"call": exprToString(n.Call)},
			Location: getLocation(n, ctx),
		})

	case *ast.CallExpr:
		// Check for panic, recover, or sync package calls
		if ident, ok := n.Fun.(*ast.Ident); ok {
			if ident.Name == "panic" {
				results = append(results, &ConcurrencyInfo{
					Type:     "panic",
					Details:  map[string]interface{}{"argCount": len(n.Args)},
					Location: getLocation(n, ctx),
				})
			} else if ident.Name == "recover" {
				results = append(results, &ConcurrencyInfo{
					Type:     "recover",
					Details:  map[string]interface{}{},
					Location: getLocation(n, ctx),
				})
			}
		}
	}

	return results, nil
}
