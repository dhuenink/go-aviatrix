// +build debug

package goaviatrix

import "log"

func debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
