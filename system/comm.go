package system

type Protocol string

const (
	Serial Protocol = "serial"
	GQL    Protocol = "graphql"
	RPC    Protocol = "rpc"
)

type Comm struct {
	Protocol Protocol
}
