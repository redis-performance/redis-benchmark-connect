# redis-benchmark-connect

## OS Configuration

For large amount of connections need to configure the target OS to release open ports quicker:

```bash
sudo sysctl -w net.ipv4.tcp_fin_timeout=10
sudo sysctl -w net.ipv4.tcp_tw_reuse=1
```
