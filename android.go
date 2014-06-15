package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"errors"
	"github.com/wmbest2/android/adb"
)

func launchAndroid() error {

	name, err := os.Hostname()
	if err != nil {
		return err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return err;
	}

	devices := adb.ListDevices(nil)

	if len(devices) == 0 {
		return errors.New("Oops: No devices found.\n");
	}

	// TODO: chooser if more than one
	d := devices[0]

	find := []string { "shell", "dumpsys", "package", *pkg }
	find = append(find, "|", "grep", "-A1", "android.intent.action.MAIN:")
	find = append(find, "|", "grep", *pkg)
//	find = append(find, "|", "sed", "-e", "\"s:.*\\(" + *pkg + "/.*\\) filter.*:\\1:g\"")

	fmt.Println(strings.Join(find, " "))

	main, err := d.ExecSync(find...)

	if err != nil {
		return err;
	}

	trimmed := strings.Trim(string(main), " ")
	component := strings.Split(trimmed, " ")
	fmt.Println("Android Component: ", component[1]);

	cmd := []string { "shell", "am", "start" }
	cmd = append(cmd, "-a", "android.intent.action.MAIN")
	cmd = append(cmd, "-n", component[1])
	cmd = append(cmd, "-e", "goatProxyHosts", strings.Join(addrs, string(0xf09f9090)))
	cmd = append(cmd, "-e", "goatProxyPort", (*port)[1:])
	cmd = append(cmd, "--activity-clear-top")

	fmt.Println(strings.Join(cmd, " "))

	result, err := d.ExecSync(cmd...)

	if err != nil {
		return err;
	}

	fmt.Println(string(result))

	return nil;
}
