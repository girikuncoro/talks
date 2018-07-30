# Lab 4: Container Packaging

This lab is originally taken from Docker community training for beginner. In this lab, we will look at some basic Docker commands and a simple build-ship-run workflow. We’ll start by running some simple containers, then we’ll use a Dockerfile to build a custom app. Finally, we’ll look at how to use bind mounts to modify a running container as you might if you were actively developing using Docker.

* Task 1: Run some simple Docker containers
* Task 2: Package and run a custom app using Docker
* Task 3: Modify a Running Website



## Prerequisites

* This lab assumes you are running a Linux or OSx machine

* Create DockerID from https://hub.docker.com

* Use your DockerID to download Docker community edition. Download the right one for your environment from https://www.docker.com/community-edition

* Clone https://github.com/dockersamples/linux_tweet_app. We will be playing around with this for our lab.

  ```
  $ git clone https://github.com/dockersamples/linux_tweet_app
  ```

  

## Task1: Run some simple Docker containers

There are different ways to use containers. These include:

1. **To run a single task:** This could be a shell script or a custom app.
2. **Interactively:** This connects you to the container similar to the way you SSH into a remote server.
3. **In the background:** For long-running services like websites and databases.

In this section you’ll try each of those options and see how Docker manages the workload.



### Run a single task in an Alpine Linux container

In this step we’re going to start a new container and tell it to run the `hostname` command. The container will start, execute the `hostname` command, then exit.

1. Run the following command in your machine.

   ```
   $ docker container run alpine hostname
   ```

   The output below shows that the `alpine:latest` image could not be found locally. When this happens, Docker automatically *pulls* it from Docker Hub.

   After the image is pulled, the container’s hostname is displayed (`888e89a3b36b` in the example below).

   ```sh
    Unable to find image 'alpine:latest' locally
    latest: Pulling from library/alpine
    88286f41530e: Pull complete
    Digest: sha256:f006ecbb824d87947d0b51ab8488634bf69fe4094959d935c0c103f4820a417d
    Status: Downloaded newer image for alpine:latest
    888e89a3b36b
   ```

2. Docker keeps a container running as long as the process it started inside the container is still running. In this case the `hostname` process exits as soon as the output is written. This means the container stops. However, Docker doesn’t delete resources by default, so the container still exists in the `Exited` state.

   List all containers.

   ```
   $ docker container ls --all
   ```

   Notice that your Alpine Linux container is in the `Exited` state.

   ```
   CONTAINER ID        IMAGE                 COMMAND                  CREATED             STATUS                     PORTS               NAMES
   1ada29881185        alpine                "hostname"               2 minutes ago       Exited (0) 2 minutes ago                       keen_jang
   ```

   > **Note:** The container ID *is* the hostname that the container displayed. In the example above it’s `888e89a3b36b`.

Containers which do one task and then exit can be very useful. You could build a Docker image that executes a script to configure something. Anyone can execute that task just by running the container - they don’t need the actual scripts or configuration information.



### Run an interactive CentOS container

You can run a container based on a different version of Linux than is running on your Docker host.

In the next example, we are going to run an CentOS container on top of your Linux Docker host (assuming you are not running CentOS machine, otherwise you can try different image).

1. Run a Docker container and access its shell.

   ```
   $ docker container run --interactive --tty --rm centos bash
   ```

   In this example, we’re giving Docker three parameters:

   - `--interactive` says you want an interactive session.
   - `--tty` allocates a pseudo-tty.
   - `--rm` tells Docker to go ahead and remove the container when it’s done executing.

   The first two parameters allow you to interact with the Docker container.

   We’re also telling the container to run `bash` as its main process (PID 1).

   When the container starts you’ll drop into the bash shell with the default prompt `root@<container id>:/#`. Docker has attached to the shell in the container, relaying input and output between your local session and the shell session in the container.

2. Run the following commands in the container.

   `ls /` will list the contents of the root director in the container, `ps aux` will show running processes in the container, `cat /etc/os-release` will show which Linux distro the container is running, in this case CentOS 7 LTS.

   ```
   $ ls /
   bin  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
   ```

   ```
   $ ps aux
   USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
   root         1  0.7  0.1  11828  2804 pts/0    Ss   12:14   0:00 bash
   root        15  0.0  0.1  51716  3424 pts/0    R+   12:14   0:00 ps aux
   ```

   ```
   $ cat /etc/os-release
   NAME="CentOS Linux"
   VERSION="7 (Core)"
   ...
   ```

3. Type `exit` to leave the shell session. This will terminate the `bash` process, causing the container to exit.

   ```
    exit
   ```

   > **Note:** As we used the `--rm` flag when we started the container, Docker removed the container when it stopped. This means if you run another `docker container ls --all`you won’t see the CentOS container.

Notice that our host VM is not running CentOS Linux, yet we were able to run an CentOS container. As previously mentioned, the distribution of Linux inside the container does not need to match the distribution of Linux running on the Docker host.

However, Linux containers require the Docker host to be running a Linux kernel. For example, Linux containers cannot run directly on Windows Docker hosts. The same is true of Windows containers - they need to run on a Docker host with a Windows kernel.

Interactive containers are useful when you are putting together your own image. You can run a container and verify all the steps you need to deploy your app, and capture them in a Dockerfile.

> You *can* [commit](https://docs.docker.com/engine/reference/commandline/commit/) a container to make an image from it - but you should avoid that wherever possible. It’s much better to use a repeatable [Dockerfile](https://docs.docker.com/engine/reference/builder/) to build your image. You’ll see that shortly.

