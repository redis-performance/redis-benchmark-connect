# redis-benchmark-connect

## How to use

Clone the repo and build the application by running:
```bash
go build -o rbc main.go
```

use the flags to connect and measure the performance of a redis instance:

```text
Please provide at least one flag
  -caCertFile string
        Path a file containing CA certificates
  -certFile string
        Path to client certificate
  -certKey string
        Path to client private key
  -hello
        Send Hello command after connection
  -help
        Display usage
  -ip string
        Redis server IP address (default "localhost")
  -numConnections int
        Number of connections to establish (default 100)
  -outputFile
        Send per connection timing to a file, provide a file name
  -parallel
        Run connections in parallel
  -password string
        Redis server password
  -port string
        Redis server port (default "6379")
  -setCommand
        Send additional SET command for every connection
  -tls
        Use TLS for connection
  -tlsSkipVerify
        Skip verification of server certificate
  -tlsVersion string
        TLS version (1.2 or 1.3) (default "1.2")
  -version
        Display version
```

## OS Configuration

For large amount of connections need to configure the target OS to release open ports quicker:

```bash
sudo sysctl -w net.ipv4.tcp_fin_timeout=10
sudo sysctl -w net.ipv4.tcp_tw_reuse=1
```

in some cases the dollowing command is also bneeded on the client side to enable more file descriptors

```bash
ulimit -n 40960
```