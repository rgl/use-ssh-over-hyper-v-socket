package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Microsoft/go-winio"
	"github.com/Microsoft/go-winio/pkg/guid"
)

func openConnection(vmidStr string) (net.Conn, string, error) {
	log.Printf("Connecting to %s on port %d...", vmidStr, *sshPortFlag)
	var err error
	var vmid guid.GUID
	if vmidStr == "localhost" {
		vmid = winio.HvsockGUIDLoopback()
	} else {
		vmid, err = guid.FromString(vmidStr)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse VM ID: %w", err)
		}
	}
	var serviceID guid.GUID
	if *sshServiceIDFlag == "" {
		serviceID = winio.VsockServiceID(uint32(*sshPortFlag))
	} else {
		serviceID, err = guid.FromString(*sshServiceIDFlag)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse service ID: %w", err)
		}
	}
	addr := &winio.HvsockAddr{
		VMID:      vmid,
		ServiceID: serviceID,
	}
	conn, err := winio.Dial(context.Background(), addr)
	return conn, addr.String(), err
}
