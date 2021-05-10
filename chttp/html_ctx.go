package chttp

import (
	"context"
	"net/http"
)

type ctxKey string

const ctxKeyComponentTree = ctxKey("chttp/component-tree")

// ComponentTree represents a single rendering of HTML components. It exposes methods that allows
// managing and communication between parent and child components
type ComponentTree struct {
	events []string
}

// GetComponentTree returns the current component tree being rendered for the given http.Request
func GetComponentTree(r *http.Request) *ComponentTree {
	tree := r.Context().Value(ctxKeyComponentTree)
	if tree == nil {
		panic("component tree does not exist")
	}

	return tree.(*ComponentTree)
}

func newComponentTree() *ComponentTree {
	return &ComponentTree{
		events: make([]string, 0),
	}
}

// Broadcast broadcasts the event that can be listened to by components in the entire tree
func (t *ComponentTree) Broadcast(event string) {
	t.events = append(t.events, event)
}

func requestWithComponentTree(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), ctxKeyComponentTree, newComponentTree())

	return r.WithContext(ctx)
}
