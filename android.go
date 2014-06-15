package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"github.com/wmbest2/android/adb"
)

func LaunchAndroid(component string, port string) {

	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	devices := adb.ListDevices(nil)

	if len(devices) == 0 {
		fmt.Println("Oops: No devices found.\n")
		return;
	} else {
		// TODO: chooser if more than one

		cmd := []string { "shell", "am", "start", "-a", "android.intent.action.MAIN" }
		cmd = append(cmd, "-n", component)
		cmd = append(cmd, "-e", "goatProxyHosts", strings.Join(addrs, string(0xf09f9090)))
		cmd = append(cmd, "-e", "goatProxyPort", port[1:])
		cmd = append(cmd, "--activity-clear-top")

		command := strings.Join(cmd, " ")
		fmt.Println(command)

		d := devices[0]
		result, err := d.ExecSync(cmd...)

		if err != nil {
			fmt.Println("ERROR: ", err)
		} else {
			fmt.Println(string(result))
		}
	}

}
