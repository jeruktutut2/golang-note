package setup

import "golang-note/globals"

func Globals() {
	globals.Session = make(map[string]string)
	globals.SSEClients = []globals.SSEClient{}
}
