package tree

func Tree(name string) map[string][]string {
	return map[string][]string{
		".":                       {name + ".yaml", "Caddyfile"},
		"api":                     {name + ".graphql"},
		"modules":                 {},
		"modules/blinky":          {"blinky.yaml"},
		"modules/blinky/firmware": {},
		"modules/blinky/service":  {},
		"modules/blinky/ui":       {},
		"external":                {},
	}
}
