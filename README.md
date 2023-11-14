# redis-benchmark-connect

## How to use

Clone the repo and build the application by running:
```bash
go build -o rbc main.go
```

use the flags to connect and measure the performance of a redis instance:

```text
Please provide at least one flag
  -certFile string
        Path to client certificate
  -certKey string
        Path to client private key
  -help
        Display usage
  -ip string
        Redis server IP address (default "your-redis-host")
  -numConnections int
        Number of connections to establish (default 100)
  -password string
        Redis server password
  -port string
        Redis server port (default "6379")
  -setCommand
        Send additional SET command for every connection
  -tls
        Use TLS for connection
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
