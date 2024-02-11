A program for file transfer over a network using TCP protocol. You can independently select host or client, and sender or receiver. You need to specify the mode, address, files, and destination directory. When connecting, the first messages sent are send or receive headers, if they match "receive - send" or "send - receive" then sending starts. This protects against erroneous connections and incompatible program versions

# Install

```bash
make
sudo make install
```
## Uninstall

```bash
sudo make uninstall
```

# Usage

```bash
pft hs <port> [files]
pft hr <port> <destdir>
pft cs <addr> <port> [files]
pft cr <addr> <port> <destdir>
```

* __hs__ - host sender
* __hr__ - host receiver
* __cs__ - client sender
* __cr__ - client receiver
* __addr__ - host ip or domain, specified for the client
* __port__ - transfer port
* __destdir__ - The folder where the received files will be uploaded, specify only for the receiver
* __files__ - files to be sent, specify only for the sender, separated by a space

# Examples

Send two archives to the device with IP=192.168.123.123, through 23232 port, the received files will be placed in the current directory.
```bash
pft hr 23232 . #host
pft cs 192.168.123.123 23232 archive1.tar archive2.tar #client
```

Send to the first connected client all files from the `dir/` folder on port 8841. 
The second command connects to the server and retrieves files from it into the
`downloads/` folder
```bash
pft hs 8841 dir/* 
pft cr some_domain.com 8841 downloads/
```
