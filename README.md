# SSH Tunneling Package

This Go package provides functionality for establishing and managing SSH tunnels.  
It simplifies the process of setting up secure SSH connections,   
making it easier to securely connect to remote servers.

## Features

- Easy setup of SSH tunnels
- Secure handling of SSH authentication

## Usage

This package can be used in any Go application that requires secure SSH connections and tunneling functionality.

## Installation

Use the standard `go get` command to install this package:

```bash
go get github.com/mrguamos/gotnl
```

# How to use

```GO
client, listener, err := gotnl.Tunnel(gotnl.Config{
    BastionHost: "bastion.host.com",
    BastionPort: "22",
    BastionUser: "ec2-user",
    TargetHost:  "rds.com",
    TargetPort:  "5432",
    LocalPort:   "5433",
    SSHKey:      "/Users/user/key",
    Passphrase:  "optional",
})
if err != nil {
    log.Fatal(err)
}
```