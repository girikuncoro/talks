# Intorduction to Linux Control Groups (Cgroups)

In this session, we are going to talk about Linux Control Groups (cgroups), which provide critical mechanism to easily manage and monitoring system resources. In our context, we are managing container resources. It does this by partitioning things like CPU time, system memory, disk and network bandwidth, into groups, then assigning tasks to those groups.

## Prerequisites

1. Linux Environment, we use Fedora 27 for the session, because Ubuntu is too maintsream at Gojek

2. Install required tools

   ```
   $ sudo yum install iotop libcgroup-tools -y
   ```



## Cgroups Overview

We can find all cgroup info related stuffs under the `/sys/fs/cgroup` directory. As we can see below, there are tons of resource group we can control through Cgroup, i.e. `blkio`, `cpu`, `memory`, etc.

```
$ ls /sys/fs/cgroup/
blkio  cpuacct      cpuset   freezer  memory   net_cls,net_prio  perf_event  systemd
cpu    cpu,cpuacct  devices  hugetlb  net_cls  net_prio          pids        unified
```

Let's look at one example of group called `blkio`. Examine what are the files inside this directory.

```
$ ls -l /sys/fs/cgroup/blkio/
total 0
-r--r--r--  1 root root 0 Aug  9 08:28 blkio.io_service_time_recursive
-r--r--r--  1 root root 0 Aug  9 08:28 blkio.io_serviced
-r--r--r--  1 root root 0 Aug  9 08:28 blkio.io_serviced_recursive
...
```

You will notice that most of the files are prefixed with `blkio` as the controller name. This differentiation is useful when we combine multiple control groups, i.e. `blkio` and `memory` controllers, such that there will be no naming conflicts and very clear distinction on the subsystem that we want to manage.

If we examine the other groups, i.e. `ls /sys/fs/cgroup/memory` we will find similar pattern. However, notice common files that can be found in each subsystem. Each root control group always has below files.

```
$ ls /sys/fs/cgroup/memory
...
tasks
cgroup.procs
notify_on_release
release_agent
```

First, tasks file contain list of PIDs attached to this control group. If we want to assign a process to particular group, we can append the process ID to tasks file.

Second, the `cgroup.procs` is also similar to tasks file, but contains the thread group ids that can be useful if we have multithreaded applications.

Finally, `notify_on_release` and `release_agent` should be used together to take next action when all processes in a control group terminate. We can then write a script into `release_agent` which will be called when the group terminate. This is very handy if we want to do logging or send notifications.

We can see list of groups that we have right now with `lscgroup`.

```
$ lscgroup
```

We see a bunch of groups are inside `system.slice` directory. This is actually created by `systemd` during init process as part of services organization, such that all systemd services groups are covered under `system.slice` directory.

## Demo on Blkio Controller Cgroups

#### Creating New Blkio Group

For demo purpose, we are going to make a child group of `blkio` called `demo1`. We do this by creating new directory called `demo1`, then Linux system will automagically populate with required files, such that it looks exactly like the parent directory, except `release_agent` file, which only presents in parent directory.

```
$ sudo su
# mkdir /sys/fs/cgroup/blkio/demo1
# ls /sys/fs/cgroup/blkio/demo1
```

See that there's new group in the system for `blkio` controller through `lscgroup` called `demo1`.

```
# lscgroup | grep -v slice
```

You can also see which groups your current process runs in (shell process) from `/proc/self/cgroup`.

```
# cat /proc/self/cgroup
```

In this demo, we want to throttle the read or write rate of a particular container (process or group of processes). Let's open 2 terminals. On first terminal, let's execute `iotop` to give us status of anything read or write into our disks with rate.

```
# iotop
```

Second terminal, we are going to create a GB file called `file-abc`, and compare disk read rates for processes that are part of our `demo1` control group, with the ones without control group. We are going to run `dd` and use `/dev/zero` as the input file, output to `file-abc`, byte count of 1MB, then we do this for 3000 times to get us 3 GB file of zeroes.

```
# dd if=/dev/zero of=file-abc bs=1M count=3000
```

As we are doing this, you can see from `iotop` our writing is around 350 MB/s.

Now we want to read from this file. Before going forward, we need to make sure to dump the cache to avoid improvement on disk reads that we are about to make.

```
# free -m
# echo 3 > /proc/sys/vm/drop_caches
# free -m
```

Let's now configure `blkio.throttle.read_bps_device` to include our configuration that limit disk reads for 5MB per second to our `sda` disk. As per the [documentation](https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt), the configuration format is as below:

```
$ echo "<major>:<minor>  <rate_bytes_per_second>" > blkio.throttle.read_bps_device
```

This requires major and minor id of device. From the information below, if we want to get `/dev/sda` , the major ID is `8` and minor ID is `0`.

```
$ ls -l /dev/sda*
brw-rw----. 1 root disk 8, 0 Aug 10 10:45 /dev/sda
brw-rw----. 1 root disk 8, 1 Aug 10 10:45 /dev/sda1
```

Let's echo the numbers to `blkio.throttle.read_bps_device` and limit the rate to 5 MB/s.

```
# cd /sys/fs/cgroup/blkio/demo1
# echo "8:0 5000000" > blkio.throttle.read_bps_device
```

Before we do comparison, let's run `dd` command again without group. We use dd command to read file-abc from disk and write it to `/dev/null`. Monitor the `iotop` on the other terminal when doing this.

```
# dd if=file-abc of=/dev/null
```

We can see that our read rate from disk is at 350+ MB/s. How is this compared with running `dd` with our `demo1` group? Let's do that but first, clean up the disk cache again.

```
# echo 3 > /proc/sys/vm/drop_caches
# free -m
```

Let's now rerun the `dd` command with `cgexec`.

```
# cgexec -g blkio:demo1 dd if=file-abc of=/dev/null
```

When doing this, we can see `iotop` is reporting reading rate at 5 MB/s. Our `demo1` group is working! How is this useful in implementing container runtime is that we can make sure the container processes that do not have any build in rate limiting can throttle it back.



## Demo on Memory Cgroups

Let's now try out different controller: memory subsystem. In this demo, we are going to see how to limit amount of memory a process (or set of processes) can use. We will first write a super simple Python script to simulate hungry memory program.

```
# cd
# cat <<EOF> greedy.py
import time
f = open("/dev/urandom", "rb")
data = bytearray()
i = 0
while True:
    data.extend(f.read(10000000))
    i += 1
    print("%dmb" % (i*10,))
    time.sleep(1)
EOF
```

The heart of the program is the for loop, in which each iteration will try and grab 10 MB of memory every second, print updates along the way, forever. Let's try to run this out as is, but be careful not to forget to cancel at some point, otherwise your machine will get crashed.

```
# python3 greedy.py
10mb
20mb
30mb
40mb
50mb
60mb
70mb
80mb
90mb
^CTraceback (most recent call last):
  File "greedy.py", line 11, in <module>
    time.sleep(1)
KeyboardInterrupt
```

Now, let's go into the `/sys/fs/cgroup/memory` directory and create a subgroup called `demo2` here. The content inside seems familiar.

```
# mkdir /sys/fs/cgroup/memory/demo2
# cd /sys/fs/cgroup/memory/demo2
# ls /sys/fs/cgroup/memory/demo2
```

The plan on this demo is to set the limit for physical memory and swap usage. We can do this by modifying the files `memory.limit_in_bytes` and `memory.swappiness`. Let's limit memory of this `demo2` group to 50 MB. We also want to disable swap to make demo interesting, otherwise our program will start eating up swap, we want the program to hit hard limit.

```
# cd /sys/fs/cgroup/memory/demo2
# echo "50000000" > memory.limit_in_bytes
# echo "0" > memory.swappiness
```

Let's now run the python greedy program again with `demo2` group using `cgexec`.

```
# cgexec -g memory:demo2 python3 greedy.py
10mb
20mb
30mb
Killed
```

It's amazing! When the process hit the limit, the process was killed at the 50MB limit. Let's take a look at last couple lines of `dmesg`, which include kernel messages on our process being killed.

```
# dmesg
...
[ 1466.966113] Memory cgroup out of memory: Kill process 9049 (python3) score 1070 or sacrifice child
[ 1466.967299] Killed process 9049 (python3) total-vm:73056kB, anon-rss:48296kB, file-rss:5384kB, shmem-rss:0kB
```



## References

* https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt
* https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt
* https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
* https://en.wikipedia.org/wiki/Cgroups
* https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/resource_management_guide/index
