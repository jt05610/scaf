package gql

import "text/template"

var gqlTemplate = template.Must(template.New("gql").Parse(`# Auto-generated GraphQL schema. Do not edit. Generated on {{.Date}}

schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}

{{- range .Types}}
type {{.Name}} {
    {{- range .Fields}}
    {{.Name}}: {{.Type}}
    {{- end}}
}
{{- end}}

type Query {
    {{- range .Types}}
    get{{.Name}}(id: ID!): {{.Name}}
    list{{.Name}}s: [{{.Name}}]
    {{- end}}
}

type Mutation {
    {{- range .Types}}
    create{{.Name}}({{range .Fields}}{{.Name}}: {{.Type}}{{if not .Last}},{{end}}{{end}}): {{.Name}}
    update{{.Name}}(id: ID!, {{range .Fields}}{{.Name}}: {{.Type}}{{if not .Last}},{{end}}{{end}}): {{.Name}}
    delete{{.Name}}(id: ID!): {{.Name}}
    {{- end}}
    {{- range .Funcs}}
    {{.Name}}({{range .Params}}{{.Name}}: {{.Type}}{{if not .Last}},{{end}}{{end}}): {{range .Return}}{{.Type}}{{if not .Last}},{{end}}{{end}}
    {{- end}}
}

type Subscription {
    {{- range .Types}}
    on{{.Name}}Created: {{.Name}}
    on{{.Name}}Updated: {{.Name}}
    on{{.Name}}Deleted: {{.Name}}
    {{- end}}
}
`))
