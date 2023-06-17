package views

import "github.com/charmbracelet/bubbles/list"

type Item interface {
	list.Item
	Title() string
	Description() string
}
