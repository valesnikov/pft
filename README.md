A program for file transfer over a network using TCP protocol. You can independently select host or client, and sender or receiver. You need to specify the mode, address, files, and destination directory. When connecting, the first messages sent are send or receive headers, if they match "receive - send" or "send - receive" then sending starts. This protects against erroneous connections and incompatible program versions

# Install
It can be installed from the [releases page](https://github.com/faceleft/pft/releases). The program consists of 1 executable file and does not require any dependencies.
Or building from source code.
```bash
go build
```
The project is compiled into 1 binary executable file, which should be moved to the PATH of your choice
```bash
#linux
sudo mv ./pft /usr/local/bin/pft
```

# Usage

```bash
pft hs -p <port> [files]
pft hr -p <port> -d <destdir>
pft cs -a <addr> -p <port> [files]
pft cr -a <addr> -p <port> -d <destdir>
```

* __hs, sh__ - host sender
* __hr, rh__ - host receiver
* __cs, sc__ - client sender
* __cr, rc__ - client receiver
* __-a --address {addr}__ - host ip or domain, specified for the client, default local machine
* __-p --port {port}__ - transfer port, default 29192
* __-d --destdir {dir}__ - The folder where the received files will be uploaded, specify only for the receiver, default "."
* __files__ - files to be sent, specify only for the sender, separated by a space

# Examples

Send two archives to the device with IP=192.168.123.123, through default (29192) port, the received files will be placed in the current directory.
```bash
pft hr #host
pft cs -a 192.168.123.123 archive1.tar archive2.tar #client
```

Send to the first connected client all files from the `dir/` folder on port 10104. 
The second command connects to the server and retrieves files from it into the
`downloads/` folder
```bash
pft hs -p 10104 -d dir/* 
pft cr -a some_domain.com -p 10104 -d downloads/
```