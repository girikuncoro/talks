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

  

## Task 1: Run some simple Docker containers

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



### Run a background MySQL container

Background containers are how you’ll run most applications. Here’s a simple example using MySQL.

1. Run a new MySQL container with the following command.

   ```
    $ docker container run \
        --detach \
        --name mydb \
        -e MYSQL_ROOT_PASSWORD=my-secret-pw \
        mysql:latest
   ```

   - `--detach` will run the container in the background.
   - `--name` will name it **mydb**.
   - `-e` will use an environment variable to specify the root password (NOTE: This should never be done in production).

   As the MySQL image was not available locally, Docker automatically pulled it from Docker Hub.

   ```
    Unable to find image 'mysql:latest' locallylatest: Pulling from library/mysql
    aa18ad1a0d33: Pull complete
    fdb8d83dece3: Pull complete
    75b6ce7b50d3: Pull complete
    ed1d0a3a64e4: Pull complete
    8eb36a82c85b: Pull complete
    41be6f1a1c40: Pull complete
    0e1b414eac71: Pull complete
    914c28654a91: Pull complete
    587693eb988c: Pull complete
    b183c3585729: Pull complete
    315e21657aa4: Pull complete
    Digest: sha256:0dc3dacb751ef46a6647234abdec2d47400f0dfbe77ab490b02bffdae57846ed
    Status: Downloaded newer image for mysql:latest
    41d6157c9f7d1529a6c922acb8167ca66f167119df0fe3d86964db6c0d7ba4e0
   ```

   As long as the MySQL process is running, Docker will keep the container running in the background.

2. List the running containers.

   ```
   $ docker container ls
   ```

   Notice your container is running.

   ```
    CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS            NAMES
    3f4e8da0caf7        mysql:latest        "docker-entrypoint..."   52 seconds ago      Up 51 seconds       3306/tcp            mydb
   ```

3. You can check what’s happening in your containers by using a couple of built-in Docker commands: `docker container logs` and `docker container top`.

   ```
   $ docker container logs mydb
   ```

   This shows the logs from the MySQL Docker container.

   ```
   ...
      2017-09-29T16:02:58.605004Z 0 [Note] Executing 'SELECT * FROM INFORMATION_SCHEMA.TABLES;' to get a list of tables using the deprecated partition engine. You may use the startup option '--disable-partition-engine-check' to skip this check.
      2017-09-29T16:02:58.605026Z 0 [Note] Beginning of list of non-natively partitioned tables
      2017-09-29T16:02:58.616575Z 0 [Note] End of list of non-natively partitioned tables
   ```

   Let’s look at the processes running inside the container.

   ```
   $ docker container top mydb
   ```

   You should see the MySQL daemon (`mysqld`) is running in the container.

   ```
    PID                 USER                TIME                COMMAND
    2876                999                 0:00                mysqld
   ```

   Although MySQL is running, it is isolated within the container because no network ports have been published to the host. Network traffic cannot reach containers from the host unless ports are explicitly published.

4. List the MySQL version using `docker container exec`.

   `docker container exec` allows you to run a command inside a container. In this example, we’ll use `docker container exec` to run the command-line equivalent of `mysql --user=root --password=$MYSQL_ROOT_PASSWORD --version`inside our MySQL container.

   ```
   $ docker exec -it mydb \
       mysql --user=root --password=$MYSQL_ROOT_PASSWORD --version
   ```

   You will see the MySQL version number, as well as a handy warning.

   ```
    mysql: [Warning] Using a password on the command line interface can be insecure.
    mysql  Ver 8.0.11 for Linux on x86_64 (MySQL Community Server - GPL)
   ```

5. You can also use `docker container exec` to connect to a new shell process inside an already-running container. Executing the command below will give you an interactive shell (`sh`) inside your MySQL container.

   ```
   $ docker exec -it mydb sh
   ```

   Notice that your shell prompt has changed. This is because your shell is now connected to the `sh` process running inside of your container.

6. Let’s check the version number by running the same command again, only this time from within the new shell session in the container.

   ```
   # mysql --user=root --password=$MYSQL_ROOT_PASSWORD --version
   ```

   Notice the output is the same as before.

7. Type `exit` to leave the interactive shell session.

   ```
    exit
   ```

