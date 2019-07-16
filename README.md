# GDSM (MATTD)

_aka MATT Daemon (Map Active TCP Tunnels) or simply mattd_

Similar to [EPMD (Erlang Port Mapper Daemon)](http://erlang.org/doc/man/epmd.html) but for Go!

### What does this do?

This is a manager based (single nde manager) solution to enable distributed Golang apps can exist.

Any worker node you have can connect to the manager. If any other worker nodes are connected to the manager, they will be updated with the new worker in a pipe delimeted list of IPs.

Example:

```ocaml
192.168.16.3:8081|192.168.16.4:8081
```

The managers IP/DNS name needs to be known.

If the manager goes down, the workers will keep the same list of workers until the manager comes back up. All workers will attempt to reconnect every second.

Semaphores are heavily utilized. No race conditions should occur.

### How to use?

```go
package main

import (
	"github.com/selfup/gdsm/gdsm"
)

func main() {
	gdsm.BootMattDaemon()
}
```

For the MANAGER node, just expose an ENV: `MANAGER=true go run main.go`

For the worker nodes: `UPLINK=manager_dns_or_ip_and:port go run main.go`

If running on the same IP you will need to assign separate PORT ENVs for each process:

Example (different shells/tabs/panes/terminals):

```
MANAGER=true go run main.go
UPLINK=localhost:8081 PORT=8082 go run main.go
UPLINK=localhost:8081 PORT=8083 go run main.go
UPLINK=localhost:8081 PORT=8084 go run main.go
```

Please reference the quite simple `docker-compose.yml` to understand the order and ENV variables needed.

Example logs of workers and a manager booting and attaching:

```ocaml
manager_1  | 2019/07/16 21:38:10 GDSM IS UP ON: 0.0.0.0:8081
workers_1  | 2019/07/16 21:38:12 GDSM IS UP ON: 0.0.0.0:8081
workers_2  | 2019/07/16 21:38:11 GDSM IS UP ON: 0.0.0.0:8081
workers_2  | 2019/07/16 21:38:11 192.168.16.3:8081
workers_1  | 2019/07/16 21:38:12 192.168.16.3:8081|192.168.16.4:8081
workers_1  | 2019/07/16 21:38:12 dial tcp 192.168.16.2:8081: connect: ..connected
workers_2  | 2019/07/16 21:38:11 dial tcp 192.168.16.2:8081: connect: ..connected
workers_2  | 2019/07/16 21:38:12 192.168.16.3:8081|192.168.16.4:8081
```

### Using netcat to ask the nodes for info

`nc ip_or_host port`

Then ask for questions in the shell:

`servers`

OR

`workers`

ALSO

`clients`

### Registry

registry.gitlab.com

### Docker Image (1.5MB)

registry.gitlab.com/selfup/gdsm:latest

### Release Repo

https://gitlab.com/selfup/gdsm

### Watch

You will need `entr`

`./scripts/watch.sh`

### Watch, build container, and run manager/workers with docker-compose

You will need `entr`

`./scripts/docker.watch.sh`
