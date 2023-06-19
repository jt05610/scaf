package views

import "github.com/jt05610/scaf/core"

type ShowListMsg struct{}
type ShowDetailMsg struct {
	System *core.System
}
type ShowConfirmMsg struct{}
type ShowCreateMsg struct{}
type ShowEditMsg struct{}
type ShowDeleteMsg struct{}
