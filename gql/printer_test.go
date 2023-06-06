package gql

import (
	"bytes"
	"github.com/jt05610/scaf/core"
	"testing"
)

func TestPrinter(t *testing.T) {
	tests := []struct {
		name     string
		module   *core.Module
		expected string
	}{
		{
			name: "Simple module",
			module: &core.Module{
				Name: "Simple",
				Port: 8081,
				Date: "05 Jun 2023",
				Types: []*core.Type{
					{
						Name: "Example",
						Fields: []*core.Field{
							{Name: "exampleField", Type: "String"},
						},
					},
				},
				Funcs: []*core.Func{
					{
						Name: "exampleFunc",
						Params: []*core.Field{
							{Name: "exampleParam", Type: "String"},
						},
						Return: []*core.Field{
							{Name: "exampleReturn", Type: "String"},
						},
					},
				},
			},
			expected: `# Auto-generated GraphQL schema. Do not edit. Generated on 05 Jun 2023

schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}
type Example {
    exampleField: String
}

type Query {
    getExample(id: ID!): Example
    listExamples: [Example]
}

type Mutation {
    createExample(exampleField: String): Example
    updateExample(id: ID!, exampleField: String): Example
    deleteExample(id: ID!): Example
    exampleFunc(exampleParam: String): String
}

type Subscription {
    onExampleCreated: Example
    onExampleUpdated: Example
    onExampleDeleted: Example
}
`,
		},
		// Other test cases...
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			printer := NewPrinter(&buf)
			err := printer.Visit(test.module)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if buf.String() != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, buf.String())
			}
		})
	}
}
