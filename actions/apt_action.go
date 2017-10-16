/*
Apt Action

Install packages and their dependencies to the target rootfs with 'apt'.

Yaml syntax:
 - action: apt
   recommends: bool
   packages:
     - package1
     - package2

Mandatory properties:

- packages -- list of packages to install

Optional properties:

- recommends -- boolean indicating if suggested packages will be installed

- allow-services -- allow to (re-)start services from packages pre/post-install scripts.
Boolean property defaults to 'false'.
*/
package actions

import (
	"github.com/go-debos/debos"
)

type AptAction struct {
	debos.BaseAction `yaml:",inline"`
	Recommends       bool
	Packages         []string
	AllowServices    bool `yaml:"allow-services"`
}

func (apt *AptAction) Run(context *debos.DebosContext) error {
	apt.LogStart()
	aptOptions := []string{"apt-get", "-y"}

	if apt.AllowServices != true {
		debos.ServicePolicyHelper(debos.DenyServices)
		defer debos.ServicePolicyHelper(debos.AllowServices)
	}

	if !apt.Recommends {
		aptOptions = append(aptOptions, "--no-install-recommends")
	}

	aptOptions = append(aptOptions, "install")
	aptOptions = append(aptOptions, apt.Packages...)

	c := debos.NewChrootCommand(context.Rootdir, context.Architecture)
	c.AddEnv("DEBIAN_FRONTEND=noninteractive")

	err := c.Run("apt", "apt-get", "update")
	if err != nil {
		return err
	}
	err = c.Run("apt", aptOptions...)
	if err != nil {
		return err
	}
	err = c.Run("apt", "apt-get", "clean")
	if err != nil {
		return err
	}

	return nil
}
