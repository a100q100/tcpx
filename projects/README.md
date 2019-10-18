## projects
Projects provide many projects based on tcpx.

## Declaration
Most projects now are used to promote tcpx using.This makes that if you want to download projects here, you might git clone the whole tcpx project.^_^

## 1. Jelly
Jelly is a server broker to handle configurations. It saves json marshaled configs and provide api to access them. Jelly can work easily in balance, this makes it well performs in micro service or distribution system. Jelly also can distinguish configs through config id, project environment.

Jelly provides three ways to notify servers to refresh config.

- Clients request in intervals.
- Using tcpx to push refresh command to client(requiring client keep connection alive).
- Web-hook.

More usage detail refer to package jelly
