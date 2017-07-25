package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/check"
	"github.com/docker/machine/libmachine/log"
	"github.com/lambda-linux/lambda-machine-local/commands/mcndirs"
)

func cmdConfig(c CommandLine, api libmachine.API) error {
	// Ensure that log messages always go to stderr when this command is
	// being run (it is intended to be run in a subshell)
	log.SetOutWriter(os.Stderr)

	target, err := targetHost(c, api)
	if err != nil {
		return err
	}

	host, err := api.Load(target)
	if err != nil {
		return err
	}

	// The second argument for `Check` method is for swarm. We set it to
	// false
	dockerHost, _, err := check.DefaultConnChecker.Check(host, false)
	if err != nil {
		return fmt.Errorf("Error running connection boilerplate: %s", err)
	}

	log.Debug(dockerHost)

	tlsCACert := filepath.Join(mcndirs.GetMachineDir(), host.Name, "ca.pem")
	tlsCert := filepath.Join(mcndirs.GetMachineDir(), host.Name, "cert.pem")
	tlsKey := filepath.Join(mcndirs.GetMachineDir(), host.Name, "key.pem")

	// TODO(nathanleclaire): These magic strings for the certificate file
	// names should be cross-package constants.
	fmt.Printf("--tlsverify\n--tlscacert=%q\n--tlscert=%q\n--tlskey=%q\n-H=%s\n",
		tlsCACert, tlsCert, tlsKey, dockerHost)

	return nil
}
