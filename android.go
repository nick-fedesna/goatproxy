package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"errors"
	"github.com/nick-fedesna/android/adb"
)

func launchAndroid() error {

	name, err := os.Hostname()
	if err != nil {
		return err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return err
	}

	devices := adb.ListDevices(nil)
	if len(devices) == 0 {
		return errors.New("Oops: No devices found.\n")
	}
	// TODO: chooser if more than one
	d := devices[0]

	dumpsys := []string { "shell", "dumpsys", "package", *pkg }
	dump, err := d.ExecSync(dumpsys...)
	if err != nil {
		return err
	}

	component := string(dump)
	actionMain := "android.intent.action.MAIN:"
	indexMain := strings.Index(component, actionMain)
	if indexMain == -1 {
		return errors.New("Oops: '" + *pkg + "' not installed.")
	}

	component = component[indexMain:]
	indexMain = strings.Index(component, *pkg)
	component = component[indexMain:]
	indexMain = strings.Index(component, "filter")
	component = component[:indexMain - 1]

	cmd := []string { "shell", "am", "start" }
	cmd = append(cmd, "-a", "android.intent.action.MAIN")
	cmd = append(cmd, "-n", component)
	cmd = append(cmd, "-e", "goatProxyHosts", strings.Join(addrs, string(0xf09f9090)))
	cmd = append(cmd, "-e", "goatProxyPort", (*port)[1:])
	cmd = append(cmd, "--activity-clear-top")

	result, err := d.ExecSync(cmd...)
	if err != nil {
		return err
	}

	fmt.Println(string(result))

	return nil
}
