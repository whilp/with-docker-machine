// with-docker-machine runs commands in an environment defined by `docker-machine`.

/*
 * Copyright (c) 2015 Will Maier <wcmaier@m.aier.us>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
)

var (
	version  string
	fVersion = flag.Bool("version", false, "print version and exit")
	fMachine = flag.String("machine", "default", "docker machine name")
)

// usage prints a helpful error message.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s COMMAND [... ARGS]\n\n", self())
	fmt.Fprintf(os.Stderr, "Run COMMAND in an environment defined by docker-machine.\n\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nArguments:\n")
	fmt.Fprintf(os.Stderr, "  COMMAND  the command to run (typically 'docker')\n")
	fmt.Fprintf(os.Stderr, "  ARGS     optional arguments to COMMAND\n")
	os.Exit(2)
}

func main() {
	var (
		exit = 0
		m    *machine
		err  error
	)
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if *fVersion {
		fmt.Printf("%s %s\n", self(), version)
		os.Exit(exit)
	} else if len(args) == 0 {
		usage()
	}

	if m, err = dockerMachineInspect(*fMachine); err != nil {
		log.Fatal(err)
	}

	env := machineEnv(m)
	os.Exit(runCommandOrPanic(args, env))
}

// machineEnv constructs a command environment that a child process
// (typically docker) can use to interact with a docker machine.
func machineEnv(m *machine) []string {
	verify := "1"
	if !m.HostOptions.EngineOptions.TlsVerify {
		verify = "0"
	}
	port := "2376"
	host := fmt.Sprintf("%s:%s", m.Driver.IpAddress, port)
	certPath := m.HostOptions.AuthOptions.StorePath
	name := m.Driver.MachineName

	return []string{
		"DOCKER_TLS_VERIFY=" + verify,
		"DOCKER_HOST=" + host,
		"DOCKER_CERT_PATH=" + certPath,
		"DOCKER_MACHINE_NAME=" + name,
	}
}

type machine struct {
	Driver struct {
		IpAddress   string
		MachineName string
	}
	HostOptions struct {
		AuthOptions struct {
			StorePath string
		}
		EngineOptions struct {
			TlsVerify bool
		}
	}
}

// dockerMachineInspect returns information about a docker machine.
func dockerMachineInspect(name string) (*machine, error) {
	var (
		m = new(machine)
		b bytes.Buffer
	)
	cmd := exec.Command("docker-machine", "inspect", name)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = &b

	if err := cmd.Run(); err != nil {
		return m, err
	}
	json.Unmarshal(b.Bytes(), &m)
	return m, nil
}

// runCommand runs command args in env.
func runCommand(args, env []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}

// runCommandOrPanic runs a command in an environment, returning its
// exit code or panicing on error.
func runCommandOrPanic(args, env []string) int {
	exit := 0
	if err := runCommand(args, env); err != nil {
		if exitError, ok := err.(*exec.ExitError); !ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			exit = waitStatus.ExitStatus()
		} else {
			log.Panic(err)
		}
	}
	return exit
}

// self returns the program's name.
func self() string {
	return path.Base(os.Args[0])
}
