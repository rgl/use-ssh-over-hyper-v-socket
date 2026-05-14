package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/mdlayher/vsock"
)

func parseCID(cidStr string) (uint32, error) {
	switch strings.ToLower(cidStr) {
	case "local", "localhost":
		return vsock.Local, nil
	case "hypervisor":
		return vsock.Hypervisor, nil
	case "host":
		return vsock.Host, nil
	}
	cid, err := strconv.ParseUint(cidStr, 0, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid CID: %s (must be number or special name)", cidStr)
	}
	return uint32(cid), nil
}

func openConnection(cidStr string) (net.Conn, string, error) {
	cid, err := parseCID(cidStr)
	if err != nil {
		return nil, "", err
	}
	log.Printf("Connecting to %s (#%d) on port %d...", cidStr, cid, *sshPortFlag)
	conn, err := vsock.Dial(cid, uint32(*sshPortFlag), nil)
	if err != nil {
		if cid == vsock.Local {
			log.Printf("WARNING: You are trying to connect to the loopback address; ensure you have loaded the vsock_loopback linux module (e.g. sudo modprobe vsock_loopback).")
		}
		return nil, "", fmt.Errorf("failed to dial: %w", err)
	}
	return conn, conn.RemoteAddr().String(), err
}
