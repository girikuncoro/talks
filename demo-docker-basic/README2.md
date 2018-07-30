# Lab 4.2: Docker Networking

This lab is originally taken from Docker community training for beginner-intermediate. In this lab, we will learn about Docker Networking concepts. You will get your hands dirty by doing:

- Task 1: Networking Basics
- Task 2: Bridge Networking



## Prerequisites

- This lab assumes you are running a Linux machine (Ubuntu preferred, use VM if you're not, with Vagrant is recommended)
- Docker engine must be running in your environment
- If you for some reason don't have access to Linux with Docker, let's just use this inteactive [Linux shell](https://training.play-with-docker.com/docker-networking-hol) 



## Task 1: Networking Basics

### Docker network command

The `docker network` command is the main command for configuring and managing container networks. Run the `docker network` command from the first terminal.

```
$ docker network
Usage:	docker network COMMAND

Manage networks

Options:
      --help   Print usage

Commands:
  connect     Connect a container to a network
  create      Create a network
  disconnect  Disconnect a container from a network
  inspect     Display detailed information on one or more networks
  ls          List networks
  prune       Remove all unused networks
  rm          Remove one or more networks

Run 'docker network COMMAND --help' for more information on a command.
```

The command output shows how to use the command as well as all of the `docker network` sub-commands. As you can see from the output, the `docker network` command allows you to create new networks, list existing networks, inspect networks, and remove networks. It also allows you to connect and disconnect containers from networks.



### List networks

Run a `docker network ls` command to view existing container networks on the current Docker host.

```
$ docker network ls
```

```
NETWORK ID          NAME                DRIVER              SCOPE
3430ad6f20bf        bridge              bridge              local
a7449465c379        host                host                local
06c349b9cc77        none                null                local
```

The output above shows the container networks that are created as part of a standard installation of Docker.

New networks that you create will also show up in the output of the `docker network ls` command.

You can see that each network gets a unique `ID` and `NAME`. Each network is also associated with a single driver. Notice that the “bridge” network and the “host” network have the same name as their respective drivers.



### Inspect a network

The `docker network inspect` command is used to view network configuration details. These details include; name, ID, driver, IPAM driver, subnet info, connected containers, and more.

Use `docker network inspect <network>` to view configuration details of the container networks on your Docker host. The command below shows the details of the network called `bridge`.

```
$ docker network inspect bridge
[
    {
        "Name": "bridge",
        "Id": "128ab629632f5f759f18ffc6bd8d63682b267ad98e7e1b11a5b2c763f2a30065",
        "Created": "2018-07-10T02:44:30.424095011Z",
        "Scope": "local",
        "Driver": "bridge",
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": null,
            "Config": [
                {
                    "Subnet": "172.17.0.0/16",
                    "Gateway": "172.17.0.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "ConfigFrom": {
            "Network": ""
        },
        "ConfigOnly": false,
        "Containers": {},
        "Options": {
            "com.docker.network.bridge.default_bridge": "true",
            "com.docker.network.bridge.enable_icc": "true",
            "com.docker.network.bridge.enable_ip_masquerade": "true",
            "com.docker.network.bridge.host_binding_ipv4": "0.0.0.0",
            "com.docker.network.bridge.name": "docker0",
            "com.docker.network.driver.mtu": "1500"
        },
        "Labels": {}
    }
]
```

> **NOTE:** The syntax of the `docker network inspect` command is `docker network inspect <network>`, where `<network>` can be either network name or network ID. In the example above we are showing the configuration details for the network called “bridge”. Do not confuse this with the “bridge” driver.



### List network driver plugins

The `docker info` command shows a lot of interesting information about a Docker installation.

Run the `docker info` command and locate the list of network plugins.

```
docker info
Containers: 0
 Running: 0
 Paused: 0
 Stopped: 0
Images: 0
Server Version: 17.03.1-ee-3
Storage Driver: aufs
<Snip>
Plugins:
 Volume: local
 Network: bridge host macvlan null overlay
Swarm: inactive
Runtimes: runc
<Snip>
```

The output above shows the **bridge**, **host**,**macvlan**, **null**, and **overlay** drivers.



## Task 2: Bridge Networking

### The Basics

Every clean installation of Docker comes with a pre-built network called **bridge**. Verify this with the `docker network ls`.

```
$ docker network ls
```

```
NETWORK ID          NAME                DRIVER              SCOPE
3430ad6f20bf        bridge              bridge              local
a7449465c379        host                host                local
06c349b9cc77        none                null                local
```

The output above shows that the **bridge** network is associated with the *bridge* driver. It’s important to note that the network and the driver are connected, but they are not the same. In this example the network and the driver have the same name - but they are not the same thing!

The output above also shows that the **bridge** network is scoped locally. This means that the network only exists on this Docker host. This is true of all networks using the *bridge* driver - the *bridge* driver provides single-host networking.

All networks created with the *bridge* driver are based on a Linux bridge (a.k.a. a virtual switch).

Install the `brctl` command and use it to list the Linux bridges on your Docker host. You can do this by running `sudo apt-get install bridge-utils`.

```
apk update
apk add bridge
```

Then, list the bridges on your Docker host, by running `brctl show`.

```
$ brctl show
```

```
bridge name	bridge id		STP enabled	interfaces
docker0		8000.024252ed52f7	no
```

The output above shows a single Linux bridge called **docker0**. This is the bridge that was automatically created for the **bridge** network. You can see that it has no interfaces currently connected to it.

You can also use the `ip a` command to view details of the **docker0** bridge.

```
$ ip a
<Snip>
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default
    link/ether 02:42:52:ed:52:f7 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 scope global docker0
       valid_lft forever preferred_lft forever
```



### Connect a container

The **bridge** network is the default network for new containers. This means that unless you specify a different network, all new containers will be connected to the **bridge** network.

Create a new container by running `docker container run -dt ubuntu sleep infinity`.

```
$ docker container run -dt ubuntu sleep infinity
```

```
Unable to find image 'ubuntu:latest' locally
latest: Pulling from library/ubuntu
d54efb8db41d: Pull complete
f8b845f45a87: Pull complete
e8db7bf7c39f: Pull complete
9654c40e9079: Pull complete
6d9ef359eaaa: Pull complete
Digest: sha256:dd7808d8792c9841d0b460122f1acf0a2dd1f56404f8d1e56298048885e45535
Status: Downloaded newer image for ubuntu:latest
846af8479944d406843c90a39cba68373c619d1feaa932719260a5f5afddbf71
```

This command will create a new container based on the `ubuntu:latest` image and will run the `sleep` command to keep the container running in the background. You can verify our example container is up by running `docker container list`.

```
$ docker container list
```

```
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
846af8479944        ubuntu              "sleep infinity"    55 seconds ago      Up 54 seconds                           heuristic_boyd
```

As no network was specified on the `docker run` command, the container will be added to the **bridge** network.

Run the `brctl show` command again.

```
$ brctl show
```

```
bridge name	bridge id		STP enabled	interfaces
docker0		8000.024252ed52f7	no		vethd630437
```

Notice how the **docker0** bridge now has an interface connected. This interface connects the **docker0** bridge to the new container just created.

You can inspect the **bridge** network again, by running `docker network inspect bridge`, to see the new container attached to it.

```
$ docker network inspect bridge
<Snip>
        "Containers": {
            "846af8479944d406843c90a39cba68373c619d1feaa932719260a5f5afddbf71": {
                "Name": "heuristic_boyd",
                "EndpointID": "1265c418f0b812004d80336bafdc4437eda976f166c11dbcc97d365b2bfa91e5",
                "MacAddress": "02:42:ac:11:00:02",
                "IPv4Address": "172.17.0.2/16",
                "IPv6Address": ""
            }
        },
<Snip>
```



### Test network connectivity

The output to the previous `docker network inspect` command shows the IP address of the new container. In the previous example it is “172.17.0.2” but yours might be different.

Ping the IP address of the container from the shell prompt of your Docker host by running `ping -c5 <IPv4 Address>`. Remember to use the IP of the container in **your** environment.

```
$ ping -c5 172.17.0.2
PING 172.17.0.2 (172.17.0.2) 56(84) bytes of data.
64 bytes from 172.17.0.2: icmp_seq=1 ttl=64 time=0.055 ms
64 bytes from 172.17.0.2: icmp_seq=2 ttl=64 time=0.031 ms
64 bytes from 172.17.0.2: icmp_seq=3 ttl=64 time=0.034 ms
64 bytes from 172.17.0.2: icmp_seq=4 ttl=64 time=0.041 ms
64 bytes from 172.17.0.2: icmp_seq=5 ttl=64 time=0.047 ms

--- 172.17.0.2 ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 4075ms
rtt min/avg/max/mdev = 0.031/0.041/0.055/0.011 ms
```

The replies above show that the Docker host can ping the container over the **bridge** network. But, we can also verify the container can connect to the outside world too. Lets log into the container, install the `ping` program, and then ping `www.github.com`.

First, we need to get the ID of the container started in the previous step. You can run `docker ps` to get that.

```
$ docker container list
```

```
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
846af8479944        ubuntu              "sleep infinity"    7 minutes ago       Up 7 minutes                            heuristic_boyd
```

Next, lets run a shell inside that ubuntu container, by running `docker exec -it <CONTAINER ID> /bin/bash`.

```
$ docker container exec -it yourcontainerid /bin/bash
root@846af8479944:/#
```

Next, we need to install the ping program. So, lets run `apt-get update && apt-get install -y iputils-ping`.

```
# apt-get update && apt-get install -y iputils-ping
```

Lets ping www.github.com by running `ping -c5 www.github.com`

```
#  ping -c5 www.github.com
```

```
PING www.docker.com (104.239.220.248) 56(84) bytes of data.
64 bytes from 104.239.220.248: icmp_seq=1 ttl=45 time=38.1 ms
64 bytes from 104.239.220.248: icmp_seq=2 ttl=45 time=37.3 ms
64 bytes from 104.239.220.248: icmp_seq=3 ttl=45 time=37.5 ms
64 bytes from 104.239.220.248: icmp_seq=4 ttl=45 time=37.5 ms
64 bytes from 104.239.220.248: icmp_seq=5 ttl=45 time=37.5 ms

--- www.docker.com ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 4003ms
rtt min/avg/max/mdev = 37.372/37.641/38.143/0.314 ms
```

Finally, lets disconnect our shell from the container, by running `exit`.

```
exit
```

We should also stop this container so we clean things up from this test, by running `docker stop <CONTAINER ID>`.

```
$ docker container stop <yourcontainerid>
```

This shows that the new container can ping the internet and therefore has a valid and working network configuration.



### Configure NAT for external connectivity

In this step we’ll start a new **NGINX** container and map port 8080 on the Docker host to port 80 inside of the container. This means that traffic that hits the Docker host on port 8080 will be passed on to port 80 inside the container.

> **NOTE:** If you start a new container from the official NGINX image without specifying a command to run, the container will run a basic web server on port 80.

Start a new container based off the official NGINX image by running `docker container run --name web1 -d -p 8080:80 nginx`.

```
$ docker container run --name web1 -d -p 8080:80 nginx
```

```
Unable to find image 'nginx:latest' locally
latest: Pulling from library/nginx
6d827a3ef358: Pull complete
b556b18c7952: Pull complete
03558b976e24: Pull complete
9abee7e1ef9d: Pull complete
Digest: sha256:52f84ace6ea43f2f58937e5f9fc562e99ad6876e82b99d171916c1ece587c188
Status: Downloaded newer image for nginx:latest
4e0da45b0f169f18b0e1ee9bf779500cb0f756402c0a0821d55565f162741b3e
```

Review the container status and port mappings by running `docker container list`.

```
$ docker container list
```

```
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                           NAMES
4e0da45b0f16        nginx               "nginx -g 'daemon ..."   2 minutes ago       Up 2 minutes        443/tcp, 0.0.0.0:8080->80/tcp   web1
```

The top line shows the new **web1** container running NGINX. Take note of the command the container is running as well as the port mapping - `0.0.0.0:8080->80/tcp` maps port 8080 on all host interfaces to port 80 inside the **web1** container. This port mapping is what effectively makes the containers web service accessible from external sources (via the Docker hosts IP address on port 8080).

Now that the container is running and mapped to a port on a host interface you can test connectivity to the NGINX web server.

To complete the following task you will need the IP address of your Docker host. This will need to be an IP address that you can reach (e.g. your lab is hosted in Azure so this will be the instance’s Public IP - the one you SSH’d into). Just point your web browser to the IP and port 8080 of your Docker host. Also, if you try connecting to the same IP address on a different port number it will fail.

If for some reason you cannot open a session from a web browser, you can connect from your Docker host using the `curl 127.0.0.1:8080` command.

```
$ curl 127.0.0.1:8080
```

```
<!DOCTYPE html>
<html>
<Snip>
<head>
<title>Welcome to nginx!</title>
    <Snip>
<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

If you try and curl the IP address on a different port number it will fail.

> **NOTE:** The port mapping is actually port address translation (PAT).