/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-16 19:03
 * description :
 * history :
 */

package app

import (
	"com/infrastructure"
	"com/service"
	"com/share/glob"
	"fmt"
	"ops/cf/app"
	"os"
	_ "runtime/debug"
	"strconv"
)

func RunSocket(ctx app.Context, port int, debug bool) {

	if gcx, ok := ctx.(*glob.AppContext); ok {
		if !gcx.Loaded {
			gcx.Init(debug)
		}
	} else {
		fmt.Println("app context err")
		os.Exit(1)
		return
	}

	if debug {
		fmt.Println("[Started]:Socket server (with debug) running on port [" +
			strconv.Itoa(port) + "]:")
		infrastructure.DebugMode = true
	} else {
		fmt.Println("[Started]:Socket server running on port [" +
			strconv.Itoa(port) + "]:")
	}

	service.ServerListen("tcp", ":"+strconv.Itoa(port), ctx)
}
