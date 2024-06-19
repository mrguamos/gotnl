GO SSH Tunnel Library

# How to use

```GO
client, listener, err := gotnl.Tunnel(gotnl.Config{
    BastionHost: viper.GetString("bastion.host"),
    BastionPort: viper.GetString("bastion.port"),
    BastionUser: viper.GetString("bastion.username"),
    TargetHost:  viper.GetString("target.host"),
    TargetPort:  viper.GetString("ssh.remote.port"),
    LocalPort:   viper.GetString("ssh.local.port"),
    SSHKey:      viper.GetString("ssh.key"),
    Passphrase:  viper.GetString("ssh.passphrase"),
})
if err != nil {
    log.Fatal(err)
}
```