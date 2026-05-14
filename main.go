package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	versionFlag      = flag.Bool("version", false, "show version")
	vmidFlag         = flag.String("vmid", "", "hyper-v VM ID (GUID or localhost) or Context ID (UINT or one of local, hypervisor, host)")
	sshUsernameFlag  = flag.String("username", "vagrant", "ssh username")
	sshPasswordFlag  = flag.String("password", "", "ssh password")
	sshPortFlag      = flag.Uint("port", 22, "ssh port")
	sshServiceIDFlag = flag.String("service-id", "", "the destination hyper-v socket Service ID (GUID).\nwhen this is used, the -port parameter is ignored.")
	commandStdinFlag = flag.String("stdin", "", "data to pass into the command stdin")
	commandFlag      = flag.String("command", "ps -efww --forest", "command to execute")
	version          = "0.0.0-dev"
	revision         = "0000000000000000000000000000000000000000"
)

func main() {
	log.SetOutput(os.Stdout) // for not disturbing PowerShell...

	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s+%s\n", version, revision)
		return
	}

	if *vmidFlag == "" {
		fmt.Fprintln(os.Stderr, "ERROR: You must set the -vmid parameter.")
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("Executing the %s command...", *commandFlag)

	exitCode, output, err := executeCommand(*commandStdinFlag, *commandFlag)
	if err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}

	log.Printf("Command ended with exit code %d and output:\n%s", exitCode, output)

	os.Exit(exitCode)
}

func executeCommand(stdin string, command string) (int, string, error) {
	config := &ssh.ClientConfig{
		User:            *sshUsernameFlag,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if *sshPasswordFlag != "" {
		config.Auth = append(config.Auth, ssh.Password(*sshPasswordFlag))
	}

	conn, addr, err := openConnection(*vmidFlag)
	if err != nil {
		return -1, "", fmt.Errorf("failed to open ssh connection: %w", err)
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}
	client := ssh.NewClient(c, chans, reqs)
	defer client.Close()

	log.Printf("Connected from %s (%s) to %s (%s)",
		client.LocalAddr(),
		client.ClientVersion(),
		client.RemoteAddr(),
		client.ServerVersion())

	log.Printf("Creating SSH session to %s...", addr)

	session, err := client.NewSession()
	if err != nil {
		return -1, "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	if stdin != "" {
		session.Stdin = bytes.NewBufferString(stdin)
	}

	output, err := session.CombinedOutput(command)
	if err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			return e.ExitStatus(), string(output), nil
		}
		return -1, "", fmt.Errorf("failed to run command: %w", err)
	}

	return 0, string(output), nil
}
