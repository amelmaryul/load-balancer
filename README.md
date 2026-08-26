# Go Load Balancer
A small load balancer written in Go. It consists of two independent tools you can compose together:

- **L4**: a TCP level load balancer. You run one instance per backend cluster. This layer handles connection relaying, least connection load balancing, health checks and automatic failovers across the servers in that cluster
- **L7**: an HTTP level load balancer. You run one instance in front of your L4 instances. It does content based routing by parsing incoming HTTP request and forwards them to the correct L4 based on path prefix.

Both L4 and L7 need to be configured with a config.json file that describes your own services and routes.

## Features
**L4**
- Raw bidirectional TCP proxying between client and backend services
- Least connection load balancing
- Background health checks with automatic failover
- Thread safe
- api/stats endpoint exposing live per server connection counts and health status

**L7**
- HTTP request parsing
- Content switching based on path prefix
- api/stats endpoint that aggregates stats from every L4 it connects to

## Configuration

### L4 (`config.json`)

```json
{
  "tcp_port": ":60",
  "http_port": ":70",
  "servers": [
    { "port": "<host>:<port>", "connections": 0, "isAlive": true },
    { "port": "<host>:<port>", "connections": 0, "isAlive": true },
    { "port": "<host>:<port>", "connections": 0, "isAlive": true }
  ]
}
```

- `tcp_port` — where this L4 instance listens for traffic to relay
- `http_port` — where this L4 instance serves its own `/api/stats`
- `servers` — the backends in this cluster. Any host:port works here. You could have it locally with localhost:80 or pointed to a real IP/hostname in a deployed environment

Make sure to have the config.json file on the same directory level as L4 main file.

### L7 (`config.json`)

```json
{
  "lbport": ":8080",
  "routes": [
    {
      "pathPrefix": "/players",
      "backend": "http://<host>:<port>/api/players",
      "stats": "http://<host>:<port>/api/stats"
    }
  ]
}
```

- `lbport` — where L7 listens for incoming client requests
- `routes` — one entry per cluster: the path prefix to match, the L4's
  address to forward matching requests to, and that L4's stats endpoint


## Running it

1. Start your backend servers for each cluster.
2. Start an L4 instance per cluster (`go run .` inside the L4 directory,
   with that cluster's `config.json`).
3. Start L7, configured with a route per L4 instance.
4. Send requests to L7's port and it will route them to the right cluster
5. Hit l7's `/api/stats` to see aggregated heath and connection data across every cluster.





