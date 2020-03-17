package models

import (
	wa "github.com/Rhymen/go-whatsapp"
	"sort"
)

type ResponderFunc func(msg wa.TextMessage, conn wa.Conn) bool

var responders = make(map[int][]ResponderFunc)

func RegisterResponder(responder ResponderFunc, priority int) {
	responders[priority] = append(responders[priority], responder)

}

func CallResponders(msg wa.TextMessage, conn wa.Conn) bool {
	keys := make([]int, 0, len(responders))
	for k := range responders {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		for _, responder := range responders[k] {
			if responder(msg, conn) {
				return true
			}
		}
	}

	return false
}