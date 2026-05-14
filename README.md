# About

[![build](https://github.com/rgl/use-ssh-over-hyper-v-socket/actions/workflows/build.yml/badge.svg)](https://github.com/rgl/use-ssh-over-hyper-v-socket/actions/workflows/build.yml)

An example Go application that executes a command using SSH over a Hyper-V VMBus Socket (aka vsock, aka hvsock).

## Usage

Start a Linux based Hyper-V Virtual Machine named `vm1` (e.g. like in https://github.com/rgl/ansible-hyperv-ubuntu-vm).

**NB** The `sshd` daemon running inside the Virtual Machine must be listening in the `vsock::22` address (like configured in https://github.com/rgl/ubuntu-vagrant).

**NB** A Hyper-V Socket can only be accessed by a user in the `Administrators` or `Hyper-V Administrators` group.

In a PowerShell session, execute:

```powershell
.\use-ssh-over-hyper-v-socket.exe `
    -username vagrant `
    -password vagrant `
    -vmid ((Get-VM vm1).ID) `
    -command 'ps -efww --forest'
```

This can also be executed inside a Windows Guest Hyper-V VM, to target a local service:

```powershell
# NB the available service-id GUIDs are listed in the Hyper-V Host registry key at:
#       HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\GuestCommunicationServices
.\use-ssh-over-hyper-v-socket.exe `
    -username vagrant `
    -password vagrant `
    -vmid localhost `
    -service-id '0000138A-FACB-11E6-BD58-64006A7986D3' `
    -command 'whoami.exe /all'
```

This can also be executed inside a Linux Guest Hyper-V VM, to target a local service:

```bash
./use-ssh-over-hyper-v-socket `
    -username vagrant `
    -password vagrant `
    -vmid localhost `
    -port 5002 `
    -command 'id'
```
