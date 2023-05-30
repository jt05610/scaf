package codegen

import "time"

type Options struct {
	Package      string        `yaml:"package"`
	UIPortStart  int           `yaml:"ui_port_start"`
	APIPortStart int           `yaml:"api_port_start"`
	PortTimeout  time.Duration `yaml:"port_timeout"`
}
