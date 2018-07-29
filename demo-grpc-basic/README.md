# Lab 2.2: Build GO-POINTS Service with gRPC

This lab gets you started with gRPC in Go with a simple GO-POINTS project. The project is about GO-POINTS because it is a very simple application which returns only random point to user, so you can focus on the gRPC pieces. This lab borrows the base content of official gRPC tutorial with various additions. By walking through this lab you will learn how to:

* Define a service in a `.proto` file.
* Generate server and client code using the protocol buffer compiler.
* Use the Go gRPC API to write a simple client and server for your service.



## Prerequisites

* This lab assumes you run Linux or OSx machine.

* Make sure to have golang installed in your machine: https://golang.org/dl

* Setup `$GOPATH` properly, make sure to have `src`, `pkg`, and `bin` folders inside.

* Create new folder for this lab: `$GOPATH/src/gojek.com/go-points`.

* **Install gRPC** 

  Use the following command to install gRPC:

  ```go
  $ go get -u google.golang.org/grpc
  ```

* **Install Protocol Buffers v3**

  Install the protoc compiler that is used to generate gRPC service code. The simplest way to do this is to download pre-compiled binaries for your platform (`protoc-<version>-<platform>.zip`) from [here](https://github.com/google/protobuf/releases):

  * Unzip this file. 

  * Move `bin/protoc` to your `/usr/local/bin`. Also move `include/*` to your `/usr/local/include`.

    ```
    $ mv bin/protoc /usr/local/bin
    $ mv include/* /usr/local/include
    ```

  * Update environment variable `PATH` to include the path to protoco binary file.

  * Verify you have `protoc` in your environment.

    ```
    $ which protoc
    /usr/local/bin/protoc
    ```

  Next, install the protoc plugin for Go.

  ```go
  $ go get -u github.com/golang/protobuf/protoc-gen-go
  ```

  The compiler plugin, `protoc-gen-go`, will be installed in `$GOBIN`, defaulting to `$GOPATH/bin`. It must be in your `$PATH` for the protocol compiler, protoc, to find it.

  ```sh
  $ export PATH=$PATH:$GOPATH/bin
  ```



## Defining the service

Navigate to our newly created directory `$GOPATH/src/gojek.com/go-points`. From here, make another directory called `gopoints` without "-" (dash) and create `point.proto` file.

```.
$ cd $GOPATH/src/gojek.com/go-points
$ mkdir gopoints
$ touch gopoints/gopoints.proto
```

Our first step is to define the gRPC service and the method request and response types using protocol buffers. Now, define `GoPoints` service inside the `gopoints.proto` file. We should also specify that we are going to use protocol buffer `proto3`.

```protobuf
syntax = "proto3";

// The GoPoints service definition
service GoPoints {}
```

Then let's define `rpc` methods inside our service definition, speicfying their request and response types. gRPC lets us define four kinds of service method: *simple RPC*, *server-side streaming RPC*, *client-side streaming RPC*, and *bidirectional streaming RPC*. Let's just use the *simple RPC* method for this lab, so our `GoPoints` service definition looks like this:

```protobuf
...
// The GoPoints service definition
service GoPoints {
    // Get random point
    rpc GetPoint(PointRequest) returns(PointReply) {}
}
```

Our `.proto` file also contains protocol buffer message type definitions for all the request and response types used in our service methods. Define our `PointRequest` and `PointReply` messages. Append below snippet at the bottom of our `.proto` file.

```protobuf
...
// The request message containing user's name
message PointRequest {
  string name = 1;
}

// The response message containing random point
message PointReply {
  string point = 1;
}
```

If you get confused or not sure your code is proper as you follow, you can always refer to our [Github](https://github.com/girikuncoro/talks/blob/master/demo-grpc-basic) for this lab.



## Generating client and server code

Next we need to generate the gRPC client and server interfaces from our `.proto` service definition. We can do this using the protocol buffer compiler `protoc` with a special gRPC Go plugin.

From the root of `go-points` directory, run:

```
$ protoc -I gopoints/ gopoints/gopoints.proto --go_out=plugins=grpc:gopoints
```

Running this command generates the following file under the `go-points/gopoints` directory:

```
$ tree
.
└── gopoints
    ├── gopoints.pb.go
    └── gopoints.proto
```

The `gopoints.pb.go` file contains:

* All the protocol buffer code to populate, serialize and retrieve our request and response message types
* An interface type (or *stub*) for clients to call with the methods defined in the `GoPoints` service.
* An interface type for servers to implement, also with the methods defined in the `GoPoints` service.



## Creating the server

First let's take a look on how we create a `GoPoints` server. There are two parts to make our `GoPoints` service do its job:

* Implementing the service interface generated from our service definition: doing the actual "work" of our service.
* Running a gRPC server to listen for requests from clients and dispatch them to the right service implementation.

#### Implementing GoPoints

From the root of `go-points` directory, make a new directory called `server` and `server.go` inside it.

```
$ cd $GOPATH/src/gojek.com/go-points
$ mkdir server
$ touch server/server.go
```

Now, fill up `server.go` with the complete snippet below, the explanation comes afterwards.

```go
package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "gojek.com/go-points/gopoints"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9001"
)

type goPointsServer struct{}

// GetPoint implements gopoints.GetPoint
func (s *goPointsServer) GetPoint(ctx context.Context, in *pb.PointRequest) (*pb.PointReply, error) {
	point := strconv.Itoa(rand.Intn(100))
	return &pb.PointReply{Point: "Hi " + in.Name + ", you received " + point + " points from GO-JEK"}, nil
}

func main() {
    log.Printf("GO-POINTS server is listening on port %s ...", port)
    
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGoPointsServer(grpcServer, &goPointsServer{})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```

As we can see above, `goPointsServer` implements our one and only method: `GetPoint`, which simply gets `Name` from the client and returns random points.

```go
func (s *goPointsServer) GetPoint(ctx context.Context, in *pb.PointRequest) (*pb.PointReply, error) {
	point := strconv.Itoa(rand.Intn(100))
	return &pb.PointReply{Point: "Hi " + in.Name + ", you received " + point + " points from GO-JEK"}, nil
}
```

The method is passed a context object for the RPC and the client's `PointRequest` protocol buffer request. It returns `PointReply` protocol buffer object with the response information and an error. In the method, we populate the `PointReply` with appropriate information, and `return` it along with a `nil` error to tell gRPC that we have finished dealing with the RPC and that the `PointReply` can be returned to the client.

Once we have implemented our method, we also need to start up the gRPC server so that clients can actually use our service. As we can see in `server.go`, we have done:

* Specifying the port we want to listen for client requests with hardcoded `9001` port.

  ```
  ...
  const (
  	port = ":9001"
  )
  
  func main() {
  ...
      lis, err := net.Listen("tcp", port)
  ...
  ```

* Create an instance of gRPC server.

  ```
  ...
      grpcServer := grpc.NewServer()
  ...
  ```

* Register our service implementation with gRPC server.

  ```
  ...
  	pb.RegisterGoPointsServer(grpcServer, &goPointsServer{})
  	// Register reflection service on gRPC server.
  	reflection.Register(grpcServer)
  ...
  ```

* Call `Serve()` on the server with our port details to do blocking wait until the process is killed. Note that calling function in `if` block like below is Golang idiom, in case you are new to Go.

  ```
  ...
      if err := grpcServer.Serve(lis); err != nil {
          ...
      }
  ...
  ```



#### Starting the server

Finally, let's just run our server with `go run` since this is a simple code snippet.

```
$ go run server/server.go
2018/07/28 22:22:26 GO-POINTS server is listening on port :9001 ...
```



## Creating the Go client

In this section, we will look at creating a Go client for our `GoPoints` service. Let's make a new directory called `client`.  

```
$ cd $GOPATH/src/gojek.com/go-points
$ mkdir client
```

Inside `client` directory, create a file called `client.go`  with below complete snippet. Explanation comes afterwards.

```go
package main

import (
	"log"
	"os"
	"time"

	pb "gojek.com/go-points/gopoints"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:9001"
	defaultName = "Nadiem"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGoPointsClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetPoint(ctx, &pb.PointRequest{Name: name})
	if err != nil {
		log.Fatalf("could not receive point: %v", err)
	}
	log.Printf("GO-POINTS: %s", resp.Point)
}
```

#### Creating a stub

As we can see in `client.go`, we have done:

* Create a gRPC *channel* to communicate with the server, in order to call service method. We create this by passing the server address and port number to `grpc.Dial()`.  We can also use `DialOptions` to set the auth credentials (e.g. TLS, JWT credentials) in `grpc.Dial` if the service we request requires that. But for our `GoPoints` service, we don't need this.

  ```
  ...
  	conn, err := grpc.Dial(address, grpc.WithInsecure())
  	if err != nil {
          ...
  	}
  	defer conn.Close()
  ...
  ```

* Once gRPC *channel* is setup, we need a client *stub* to perform RPCs. We get this using the `NewGoPointsClient` method provided in the `pb` package we generated from our `.proto`.

  ```
  ...
  	client := pb.NewGoPointsClient(conn)
  ...
  ```

#### Calling service method

Now, let's take a look at how we call our service method. Note that in gRPC-Go, RPCs operate in a blocking/syncrhonous mode, which means that the RPC call waits for the server to respond, and will either return a response or an error.

* Calling our RPC `GetPoint` is as straightforward as calling a local method. As we can see below, we call the method on the stub we got earlier. In our method parameter, we create and populate a request protocol buffer object (`PointRequest`). We also pass a `context.Context` object which lets us change our RPC's behaviour if necessary, such as time-out/cancel an RCP in flight.

  ```
  ...
      resp, err := client.GetPoint(ctx, &pb.PointRequest{Name: name})
  	if err != nil {
  	    ...
  	}
  ...
  ```

* If the call doesn't return an error, then we can read the response informatino from the server from the first return value.

  ```
  ...
  	log.Printf("GO-POINTS: %s", resp.Point)
  ```

  

#### Running our client

Finally, we have completed all client and server code. 

```
$ pwd
$GOPATH/src/gojek.com/go-points

$ tree
.
├── client
│   └── client.go
├── gopoints
│   ├── gopoints.pb.go
│   └── gopoints.proto
└── server
    └── server.go

3 directories, 4 files
```



With the `GoPoints` server running in background, let's run the client simply by doing:

```
$ go run client/client.go
2018/07/28 23:03:04 GO-POINTS: Hi nadiem, you received 81 points from GO-JEK
```



## Optional Challenge

* Add one more method to `gopoints.proto` and regenerate the code.
* Create one more client with Ruby based gRPC. The client should be able to get similar output to our Golang client and the response of the other method you just created. Follow the official [Ruby tutorial](https://grpc.io/docs/quickstart/ruby.html) for reference.

