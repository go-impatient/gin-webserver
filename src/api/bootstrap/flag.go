package bootstrap

var debug bool

// Args provides all startup parameters for web-api service.
var Args args

type args struct {
	Host       string
	Port       int
	ConfigFile string
}

// SetDebug sets application running mode.
func SetDebug(mode string) {
	if mode == "debug" {
		debug = true
	}
}

// Debug returns application running mode.
func Debug() bool {
	return debug
}

