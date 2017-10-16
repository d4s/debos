package debos

import (
	"fmt"
	"os"
)

type ServicePolicy int

const (
	DenyServices = iota
	AllowServices
)

const debianPolicyHelper = "/usr/sbin/policy-rc.d"

/*
ServicePolicyHelper allow or prohibit start/stop services during
packages installation on OS level.
Currently supports only debian-based family.
*/
func ServicePolicyHelper(action ServicePolicy) error {

	var helper = []byte(`#!/bin/sh
exit 101
`)
	switch action {
	case DenyServices:
		if _, err := os.Stat(debianPolicyHelper); os.IsNotExist(err) {
			return nil
		}
		if err := os.Remove(debianPolicyHelper); err != nil {
			return err
		}
	case AllowServices:
		if _, err := os.Stat(debianPolicyHelper); os.IsExist(err) {
			return fmt.Errorf("Policy helper file '%s' exists already", debianPolicyHelper)
		}
		pf, err := os.Create(debianPolicyHelper)
		if err != nil {
			return err
		}
		defer pf.Close()

		if _, err := pf.Write(helper); err != nil {
			return err
		}

		if err := pf.Chmod(0755); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Services policy type is not supported")
	}

	return nil
}
