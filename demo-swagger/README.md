# Lab 2: Generate REST API Server with Swagger

This example walks you through a hypothetical project, building a todo list. Originally stolen from https://goswagger.io with various modifications to make it easier to follow.

It uses a todo list because this is well-understood application, so you can focus on the go-swagger pieces. In this example we build a server and a client.



## Pre-requisites

  - This tutorial assumes you run Linux or OSx machine.
  - Make sure to have golang installed in your machine: https://golang.org/dl
  - Setup `$GOPATH` properly, make sure to have `src`, `pkg`, and `bin` folders inside.
  - Create new folder for our lab today: `$GOPATH/src/gojek.com/go-academy`.



## Preparing Swagger Spec (OpenAPI)

Navigate to our newly created directory `$GOPATH/src/gojek.com/go-academy`. From here, we will create our application. Start with `swagger init` command as below:

```
$ swagger init spec \
    --title "A Todo list application" \
    --description "From the todo list tutorial gojek" \
    --version 1.0.0 \
    --scheme http \
    --consumes application/com.gojek.todo-list.v1+json \
    --produces application/com.gojek.todo-list.v1+json
```

This gives you a skeleton `swagger.yml` file:

```yaml
consumes:
- application/com.gojek.todo-list.v1+json
definitions: {}
info:
  description: From the todo list tutorial on goswagger.io
  title: A Todo list application
  version: 1.0.0
paths: {}
produces:
- application/com.gojek.todo-list.v1+json
schemes:
- http
swagger: "2.0"
```

This doesn't do much but it does pass validation:

```
$ pwd
~/go/src/gojek.com/go-academy
$ swagger validate ./swagger.yml
The swagger spec at "./swagger.yml" is valid against swagger specification 2.0
```

Now that you have an empty but valid specification document, it's time to declare some models and endpoints for the API. You'll need a model to represent a todo item, so you can define that in the definitions section. Add below spec to the bottom of your existing `swagger.yml`

```yaml
...
definitions:
  item:
    type: object
    required:
      - description
    properties:
      id:
        type: integer
        format: int64
        readOnly: true
      description:
        type: string
        minLength: 1
      completed:
        type: boolean
```

In this model definition we say that the model `item` is an _object_ with a required property `description`. This item model has 3 properties: `id`, `description`, and `completed`. The `id` property is an int64 value and is marked as _readOnly_, meaning that it will be provided by the API server and it will be ignored when the item is created.

This document also says that the description must be at least 1 char long, which results in a string property that's not a pointer.

At this moment you have enough so that actual code could be generated, but let's continue defining the rest of the API so that the code generation will be more useful. Now that you have a model so you can add some endpoints to list the todo's. Modify the `paths` property in our `swagger.yml` as below:

```yaml
...
paths:
  /:
    get:
      tags:
        - todos
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
...
```

This snippet of yaml defines a `GET /` operation and tags it with _todos_. Tagging things is useful for many tools, for example helping UI tools group endpoints appropriately. Code generators might turn them into 'controllers'. There is also a response defined with a generic description about the response content. Note that some generators will put the description into the http status message. The response also defines endpoint's return type. In this case the endpoint returns a list of todo items, so the schema is an _array_ and the array will contain `item` objects, which you defined previously.

But wait a minute, what if there are 100's of todo items, will we just return all of them for everybody?  It would be good to add a `since` and `limit` param here. The ids will have to be ordered for a `since` param to work but you're in control of that so that's fine. Add `parameters` to our `GET /` path as below:

```yaml
...
paths:
  /:
    get:
      tags:
        - todos
      parameters:
        - name: since
          in: query
          type: integer
          format: int64
        - name: limit
          in: query
          type: integer
          format: int32
          default: 20
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
```

With this new version of the operation you now have query params. These parameters have defaults so users can leave them off and the API will still function as intended.

However, this definition is extremely optimistic and only defines a response for the "happy path". It's very likely that the API will need to return errors too. That means you have to define a model errors, as well as at least one more response definition to cover the error response.

The error definition looks like below. Append this to the bottom of our `swagger.yml`.

```yaml
...
definitions:
...
  error:
    type: object
    required:
      - message
    properties:
      code:
        type: integer
        format: int64
      message:
        type: string
```

For the error response you can use the default response, on the assumption that every successful response from your API is defying the odds. Add `default` response to our `GET /` path pointing to the newly added `error` definition.

```yaml
...
paths:
  /:
    get:
      tags:
        - todos
      parameters:
        - name: since
          in: query
          type: integer
          format: int64
        - name: limit
          in: query
          type: integer
          format: int32
          default: 20
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
        default:
          description: generic error response
          schema:
            $ref: "#/definitions/error"
```

At this point you've defined your first endpoint completely. To improve the strength of this contract you could define responses for each of the status codes and perhaps return different error messages for different statuses. For now, the status code will be provided in the error message.

Try validating the specification again with `swagger validate ./swagger.yml` to ensure that code generation will work as expected. Generating code from an invalid specification leads to unpredictable results.

```
$ swagger validate ./swagger.yml
The swagger spec at "./swagger.yml" is valid against swagger specification 2.0
```

Your completed spec should now look like below. Move the `produces`, `schemes`, and `swagger` properties to the top so it looks cleaner.

```yaml
swagger: "2.0"
info:
  description: From the todo list tutorial gojek
  title: A Todo list application
  version: 1.0.0
consumes:
- application/com.gojek.todo-list.v1+json
produces:
- application/com.gojek.todo-list.v1+json
schemes:
- http
paths:
  /:
    get:
      tags:
        - todos
      parameters:
        - name: since
          in: query
          type: integer
          format: int64
        - name: limit
          in: query
          type: integer
          format: int32
          default: 20
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
        default:
          description: generic error response
          schema:
            $ref: "#/definitions/error"
definitions:
  item:
    type: object
    required:
      - description
    properties:
      id:
        type: integer
        format: int64
        readOnly: true
      description:
        type: string
        minLength: 1
      completed:
        type: boolean
  error:
    type: object
    required:
      - message
    properties:
      code:
        type: integer
        format: int64
      message:
        type: string
```



## Generate API Server Code using Swagger Spec

Now, let's do the most interesting thing we can do with swagger: generate code to run our server. Execute the `generate` command as below and see the output.

```
$ pwd
~/go/src/gojek.com/go-academy
$ swagger generate server -A todo-list -f ./swagger.yml
2018/07/27 16:38:52 validating spec /Users/girikuncoro/Hacks/golang/src/go-academy/swagger.yml
2018/07/27 16:38:53 preprocessing spec with option:  minimal flattening
2018/07/27 16:38:53 building a plan for generation
2018/07/27 16:38:53 planning definitions
2018/07/27 16:38:53 planning operations
2018/07/27 16:38:53 grouping operations into packages
2018/07/27 16:38:53 planning meta data and facades
2018/07/27 16:38:53 rendering 2 models
2018/07/27 16:38:53 rendering 1 templates for model error
2018/07/27 16:38:53 name field error
2018/07/27 16:38:53 package field models
2018/07/27 16:38:53 creating generated file "error.go" in "models" as definition
2018/07/27 16:38:53 executed template asset:model
2018/07/27 16:38:53 rendering 1 templates for model item
2018/07/27 16:38:53 name field item
2018/07/27 16:38:53 package field models
2018/07/27 16:38:53 creating generated file "item.go" in "models" as definition
2018/07/27 16:38:53 executed template asset:model
2018/07/27 16:38:53 rendering 1 operation groups (tags)
2018/07/27 16:38:53 rendering 1 operations for todos
2018/07/27 16:38:53 rendering 4 templates for operation todo-list
2018/07/27 16:38:53 name field Get
2018/07/27 16:38:53 package field todos
2018/07/27 16:38:53 creating generated file "get_parameters.go" in "restapi/operations/todos" as parameters
2018/07/27 16:38:53 executed template asset:serverParameter
2018/07/27 16:38:53 name field Get
2018/07/27 16:38:53 package field todos
2018/07/27 16:38:53 creating generated file "get_urlbuilder.go" in "restapi/operations/todos" as urlbuilder
2018/07/27 16:38:53 executed template asset:serverUrlbuilder
2018/07/27 16:38:53 name field Get
2018/07/27 16:38:53 package field todos
2018/07/27 16:38:53 creating generated file "get_responses.go" in "restapi/operations/todos" as responses
2018/07/27 16:38:53 executed template asset:serverResponses
2018/07/27 16:38:53 name field Get
2018/07/27 16:38:53 package field todos
2018/07/27 16:38:53 creating generated file "get.go" in "restapi/operations/todos" as handler
2018/07/27 16:38:53 executed template asset:serverOperation
2018/07/27 16:38:53 rendering 0 templates for operation group todo-list
2018/07/27 16:38:53 rendering support
2018/07/27 16:38:53 rendering 6 templates for application TodoList
2018/07/27 16:38:53 name field TodoList
2018/07/27 16:38:53 package field operations
2018/07/27 16:38:53 creating generated file "configure_todo_list.go" in "restapi" as configure
2018/07/27 16:38:53 executed template asset:serverConfigureapi
2018/07/27 16:38:53 name field TodoList
2018/07/27 16:38:53 package field operations
2018/07/27 16:38:53 creating generated file "main.go" in "cmd/todo-list-server" as main
2018/07/27 16:38:53 executed template asset:serverMain
2018/07/27 16:38:53 name field TodoList
2018/07/27 16:38:53 package field operations
2018/07/27 16:38:53 creating generated file "embedded_spec.go" in "restapi" as embedded_spec
2018/07/27 16:38:53 executed template asset:swaggerJsonEmbed
2018/07/27 16:38:53 name field TodoList
2018/07/27 16:38:53 package field operations
2018/07/27 16:38:53 creating generated file "server.go" in "restapi" as server
2018/07/27 16:38:53 executed template asset:serverServer
2018/07/27 16:38:53 name field TodoList
2018/07/27 16:38:53 package field operations
2018/07/27 16:38:53 creating generated file "todo_list_api.go" in "restapi/operations" as builder
2018/07/27 16:38:53 executed template asset:serverBuilder
2018/07/27 16:38:53 name field TodoList
2018/07/27 16:38:53 package field operations
2018/07/27 16:38:53 creating generated file "doc.go" in "restapi" as doc
2018/07/27 16:38:53 executed template asset:serverDoc
2018/07/27 16:38:53 Generation completed!

For this generation to compile you need to have some packages in your GOPATH:

	* github.com/go-openapi/runtime
	* github.com/jessevdk/go-flags
```

Get the server dependencies of `runtime` and `go-flags` with below command:

```
$ go get -u github.com/go-openapi/runtime
$ go get -u github.com/jessevdk/go-flags
```

Now let's see what are the files have been generated. Execute `tree` to see the list.

```
$ tree
.
├── cmd
│   └── todo-list-server
│       └── main.go
├── models
│   ├── error.go
│   └── item.go
├── restapi
│   ├── configure_todo_list.go
│   ├── doc.go
│   ├── embedded_spec.go
│   ├── operations
│   │   ├── todo_list_api.go
│   │   └── todos
│   │       ├── get.go
│   │       ├── get_parameters.go
│   │       ├── get_responses.go
│   │       └── get_urlbuilder.go
│   └── server.go
└── swagger.yml
```

In this file tree you see that there is a `cmd/todo-list-server` directory. The swagger generator adds "-server" to the application name that you gave via the `-A` argument.

The next section in this tree is the `models` package. This package contains go representations for all item definitions in the swagger spec document.

The last section is `restapi`. The `restapi` package is generated based on the `paths` property in the swagger specification. The go swagger generator uses tags to group the operations into packages.

You can also name operation by specifying an `operationId` in the specification for a path. Let's add `findTodos` as our `operationId` for `GET /`.


```yaml
...
paths:
  /:
    get:
      tags:
        - todos
      operationId: findTodos
      ...
```

Cleanup the `restapi` folder and generate the code again, you will see `operationId` values are used to name the generated files.

```
$ rm -rf restapi
$ swagger generate server -A todo-list -f ./swagger.yml
$ tree
.
├── cmd
│   └── todo-list-server
│       └── main.go
├── models
│   ├── error.go
│   └── item.go
├── restapi
│   ├── configure_todo_list.go
│   ├── doc.go
│   ├── embedded_spec.go
│   ├── operations
│   │   ├── todo_list_api.go
│   │   └── todos
│   │       ├── find_todos.go
│   │       ├── find_todos_parameters.go
│   │       ├── find_todos_responses.go
│   │       └── find_todos_urlbuilder.go
│   └── server.go
└── swagger.yml
```

You can see that the files under `restapi/operations/todos` now use the `operationId` as part of the generated file names.

At this point can start the server, but first let's see what `--help` gives you. First install the server binary and then run it:

```
$ go install ./cmd/todo-list-server
$ todo-list-server --help
Usage:
  todo-list-server [OPTIONS]

From the todo list tutorial on goswagger.io

Application Options:
      --scheme=            the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=   grace period for which to wait before shutting down the server (default: 10s)
      --max-header-size=   controls the maximum number of bytes the server will read parsing the request header's keys and values, including the
                           request line. It does not limit the size of the request body. (default: 1MiB)
      --socket-path=       the unix socket to listen on (default: /var/run/todo-list.sock)
      --host=              the IP to listen on (default: localhost) [$HOST]
      --port=              the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=      limit the number of outstanding requests
      --keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
                           (default: 3m)
      --read-timeout=      maximum duration before timing out read of the request (default: 30s)
      --write-timeout=     maximum duration before timing out write of the response (default: 60s)
      --tls-host=          the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=          the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=   the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=           the private key to use for secure conections [$TLS_PRIVATE_KEY]
      --tls-ca=            the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=  limit the number of outstanding requests
      --tls-keep-alive=    sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=  maximum duration before timing out read of the request
      --tls-write-timeout= maximum duration before timing out write of the response

Help Options:
  -h, --help               Show this help message
```

If you run your application now it will start on a random port by default. This might not be what you want, so you can configure a port through `--port` command line argument or a `PORT` env var.

```
$ todo-list-server --port 9000
serving todo list at http://127.0.0.1:9000
```

You can use `curl` to check your API:

```
$ curl -i http://127.0.0.1:9000
```
```http
HTTP/1.1 501 Not Implemented
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 11:33:41 GMT
Content-Length: 57

"operation todos.FindTodos has not yet been implemented"
```

As you can see, the generated API isn't very usable yet, but we know it runs and does something. To make it useful you'll need to implement the actual logic behind those endpoints. And you'll also want to add some more endpoints, like adding a new todo item and updating an existing item to change its description or mark it completed.

To supporting  adding a todo item you should define a `POST` operation. Append below `post` snippet below your `get` property for `/` path.

```yaml
...
paths:
  /:
    get:
...
    post:
      tags:
        - todos
      operationId: addOne
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/item"
      responses:
        201:
          description: Created
          schema:
            $ref: "#/definitions/item"
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
...
```

This snippet has something new. You see that the parameters are defined using a `schema` that references our exist item model. Remember that we defined this object's `id` key as _readOnly_, so it will not be accepted as part of the `POST` body.

Next you can define a `DELETE` to remove a todo item from the list. Add new path `/{id}` and fill it up with `delete` snippet as below.

```yaml
...
paths:
  /:
...
  /{id}:
    parameters:
      - type: integer
        format: int64
        name: id
        in: path
        required: true
    delete:
      tags:
        - todos
      operationId: destroyOne
      parameters:
        - type: integer
          format: int64
          name: id
          in: path
          required: true
      responses:
        204:
          description: Deleted
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
...
```

This time you're you're defining a parameter that is part of the `path`. This operation will look in the URI templated path for an `id`.  Since there's nothing to return after a delete, the success response  is `204 No Content`.

Finally, you need to define a way to update an existing item.

```yaml
...
paths:
...
  /{id}:
...
    put:
      tags: ["todos"]
      operationId: updateOne
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/item"
      responses:
        '200':
          description: OK
          schema:
            $ref: "#/definitions/item"
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
    delete:
...
```

There are two approaches typically taken for updates. A `PUT` indicates that the entire entity is being replaced, and a `PATCH` indicates that only the fields provided in the request should be updated. In the above example, you can see that the `PUT` "brute force" approach is being used.

Another thing to note is that because the `/{id}` path is shared for both `DELETE` and `PUT`, they can share a `parameters` definition.

At this point you should have a complete specification for the todo list API:

```yaml
swagger: "2.0"
info:
  description: From the todo list tutorial gojek
  title: A Todo list application
  version: 1.0.0
consumes:
- application/com.gojek.todo-list.v1+json
produces:
- application/com.gojek.todo-list.v1+json
schemes:
- http
paths:
  /:
    get:
      tags:
        - todos
      operationId: find_todos
      parameters:
        - name: since
          in: query
          type: integer
          format: int64
        - name: limit
          in: query
          type: integer
          format: int32
          default: 20
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
        default:
          description: generic error response
          schema:
            $ref: "#/definitions/error"
    post:
      tags:
        - todos
      operationId: addOne
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/item"
      responses:
        201:
          description: Created
          schema:
            $ref: "#/definitions/item"
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
  /{id}:
    parameters:
      - type: integer
        format: int64
        name: id
        in: path
        required: true
    put:
      tags: ["todos"]
      operationId: updateOne
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/item"
      responses:
        '200':
          description: OK
          schema:
            $ref: "#/definitions/item"
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
    delete:
      tags:
        - todos
      operationId: destroyOne
      parameters:
        - type: integer
          format: int64
          name: id
          in: path
          required: true
      responses:
        204:
          description: Deleted
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
definitions:
  item:
    type: object
    required:
      - description
    properties:
      id:
        type: integer
        format: int64
        readOnly: true
      description:
        type: string
        minLength: 1
      completed:
        type: boolean
  error:
    type: object
    required:
      - message
    properties:
      code:
        type: integer
        format: int64
      message:
        type: string
```

This is a good time to sanity check and by validating the schema:

```
$ swagger validate ./swagger.yml
The swagger spec at "./swagger.yml" is valid against swagger specification 2.0
```

Now you're ready to generate the API and start filling in the actual operations:

```
$ swagger generate server -A TodoList -f ./swagger.yml
...
$ tree
.
├── cmd
│   └── todo-list-server
│       └── main.go
├── models
│   ├── error.go
│   └── item.go
├── restapi
│   ├── configure_todo_list.go
│   ├── doc.go
│   ├── embedded_spec.go
│   ├── operations
│   │   ├── todo_list_api.go
│   │   └── todos
│   │       ├── add_one.go
│   │       ├── add_one_parameters.go
│   │       ├── add_one_responses.go
│   │       ├── add_one_urlbuilder.go
│   │       ├── destroy_one.go
│   │       ├── destroy_one_parameters.go
│   │       ├── destroy_one_responses.go
│   │       ├── destroy_one_urlbuilder.go
│   │       ├── find_todos.go
│   │       ├── find_todos_parameters.go
│   │       ├── find_todos_responses.go
│   │       ├── find_todos_urlbuilder.go
│   │       ├── update_one.go
│   │       ├── update_one_parameters.go
│   │       ├── update_one_responses.go
│   │       └── update_one_urlbuilder.go
│   └── server.go
└── swagger.yml

6 directories, 25 files
```



## Implement POST, GET, PUT, DELETE Handlers

To implement the core of your application you start by editing `restapi/configure_todo_list.go`. This file is safe to edit. Its content will not be overwritten if you run `swagger generate` again the future.

The simplest way to implement this application is to simply store all the todo items in a golang `map`. This provides a simple way to move forward without bringing in complications like a database or files.

To do this you'll need a map and a counter to track the last assigned id. Add this snippet to `configure_todo_list.go` on top of the file after `import` section.

```go
import (
...
// the variables we need throughout our implementation
var items = make(map[int64]*models.Item)
var lastID int64

func configureFlags(
...
```



#### DELETE Handler (destroyOne)

Let's start with implementing delete handler. Add `deleteItem` helper function first and mutex lock. In Golang, we have to manage locking shared resources to prevent race condition. In this case, our shared resource is the `items` map.

```
...
var itemsLock = &sync.Mutex{}

func deleteItem(id int64) error {
	itemsLock.Lock()
	defer itemsLock.Unlock()

	_, exists := items[id]
	if !exists {
		return errors.NotFound("not found: item %d", id)
	}

	delete(items, id)
	return nil
}
...
```

Now, find the `configureAPI` function and implement destroy handler as below using `deleteItem` helper:

```go
func configureAPI(...) {
...
	api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
		if err := deleteItem(params.ID); err != nil {
			return todos.NewDestroyOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return todos.NewDestroyOneNoContent()
	})
...
}
```

After deleting the item from the store, you need to provide a response. The code generator created responders for each response you defined in the the swagger specification, and you can see how one of those is being used in the sample code above.



#### GET Handler (findTodos)

Now, let's implement our GET handler which is called `FindTodos`. Add helper function `allItems` that will act as our main logic of `FindTodos` handler.

```
...
func allItems(since int64, limit int32) (result []*models.Item) {
	result = make([]*models.Item, 0)
	for id, item := range items {
		if len(result) >= int(limit) {
			return
		}
		if since == 0 || id > since {
			result = append(result, item)
		}
	}
	return
}
...
```

Let's use `allItems` helper and refactor `api.TodosFindTodosHandler` in `configureAPI` function. We will be using our generated `models.Error` here, so don't forget to import `models`.

```
import (
...
  "gojek.com/go-academy/models"
)
...
function configureAPI(...) {
...
	api.TodosFindTodosHandler = todos.FindTodosHandlerFunc(func(params todos.FindTodosParams) middleware.Responder {
		mergedParams := todos.NewFindTodosParams()
		mergedParams.Since = swag.Int64(0)
		if params.Since != nil {
			mergedParams.Since = params.Since
		}
		if params.Limit != nil {
			mergedParams.Limit = params.Limit
		}
		return todos.NewFindTodosOK().WithPayload(allItems(*mergedParams.Since, *mergedParams.Limit))
	})
...
}
```



#### POST Handler (addOne)

Now, we will need `addItem` helper function to implement our `addOne` handler.

```
...
func addItem(item *models.Item) error {
	if item == nil {
		return errors.New(500, "item must be present")
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	newID := newItemID()
	item.ID = newID
	items[newID] = item

	return nil
}
...
```

And then, add handler to `configureAPI` function for `addOne`  that uses `addItem` helper.

```
...
function configureAPI(...) {
...
	api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
		if err := addItem(params.Body); err != nil {
			return todos.NewAddOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return todos.NewAddOneCreated().WithPayload(params.Body)
	})
...
}
...
```



### UPDATE Handler (updateOne)

Finally, let's implement `updateOne` handler to complete our API server. Add `updateItem` helper function.

```
...
func updateItem(id int64, item *models.Item) error {
	if item == nil {
		return errors.New(500, "item must be present")
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	_, exists := items[id]
	if !exists {
		return errors.NotFound("not found: item %d", id)
	}

	item.ID = id
	items[id] = item
	return nil
}
...
```

Use this and add `updateOne` handler to `configureAPI` function.

```
...
function configureAPI(...) {
...
	api.TodosUpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams) middleware.Responder {
		if err := updateItem(params.ID, params.Body); err != nil {
			return todos.NewUpdateOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return todos.NewUpdateOneOK().WithPayload(params.Body)
	})
...
}
...
```

Now, we have completed all handlers of our API server, let's generate, install, and run our server again.

```
$ swagger generate server -A TodoList -f ./swagger.yml
...
$ go install ./cmd/todo-list-server
$ todo-list-server --port 9000
2018/07/27 19:12:14 Serving todo list at http://127.0.0.1:9000
```

If you are facing problem or not sure whether your implemention is correct, you can compare `configure_todo_list.go` and `swagger.yml` from our [Github](https://github.com/girikuncoro/talks/blob/master/demo-swagger/swagger.yml).



## Test Our Generated API Server

So assuming you implemented all endpoints correctly, you're all set to test our running server. Let's test out the `GET /` path.

```bash
$ curl -i localhost:9000
```
```http
HTTP/1.1 200 OK
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 13:44:37 GMT
Content-Length: 3

[]
```

We can see above returning `200` code and empty list since we have not added anything yet. 

Swagger is all nothing but contract. Remember that we defined the server will only consume `application/com.gojek.todo-list.v1+json` content type in our `swagger.yml`. Try to add an item to our server.

```bash
$ curl -i localhost:9000 -d "{\"description\":\"message $RANDOM\"}"
```

```http
HTTP/1.1 415 Unsupported Media Type
Content-Type: application/json
Date: Fri, 27 Jul 2018 13:49:12 GMT
Content-Length: 145

{"code":415,"message":"unsupported media type \"application/x-www-form-urlencoded\", only [application/com.gojek.todo-list.v1+json] are allowed"}
```
We can see above returning `Unsupported Media Type`. We need to pass request with proper content type. Let's do it again.

```bash
$ curl -i localhost:9000 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/com.gojek.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 13:52:38 GMT
Content-Length: 39

{"description":"message 26463","id":1}
```
The server returns `201` which means new data has successfully been added. Let's add couple more.

```bash
$ curl -i localhost:9000 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/com.gojek.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 13:53:51 GMT
Content-Length: 39

{"description":"message 18912","id":2}
```
```bash
$ curl -i localhost:9000 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/com.gojek.todo-list.v1+json'
```

```http
HTTP/1.1 201 Created
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 13:56:38 GMT
Content-Length: 39

{"description":"message 14556","id":3}
```

We understand that server assigns unique `id` to response incrementally. This is due to `newItemID` helper that we added in `configure_todo_list.go`. Now we can hit `GET /` again and see all 3 items we have added.

```bash
$ curl -i localhost:9000
```

```http
HTTP/1.1 200 OK
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 13:59:19 GMT
Content-Length: 119

[{"description":"message 18912","id":2},{"description":"message 14556","id":3},{"description":"message 26463","id":1}]
```

Now, let's modify item 3 with PUT request.

```bash
$ curl -i localhost:9000/3 -X PUT -H 'Content-Type: application/com.gojek.todo-list.v1+json' -d '{"description":"topup gopay"}'
```
```http
HTTP/1.1 200 OK
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 14:02:49 GMT
Content-Length: 37

{"description":"topup gopay","id":3}
```
We can see server returns `200` which means item has been successfully updated. We can see reflected change by calling `GET /` again.

```bash
$ curl -i localhost:9000
```
```http
HTTP/1.1 200 OK
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 14:03:59 GMT
Content-Length: 117

[{"description":"topup gopay","id":3},{"description":"message 26463","id":1},{"description":"message 18912","id":2}]
```
Finally, let's try to delete one of the items.

```bash
$ curl -i localhost:9000/1 -X DELETE -H 'Content-Type: application/com.gojek.todo-list.v1+json'
```
```http
HTTP/1.1 204 No Content
Date: Fri, 27 Jul 2018 14:05:03 GMT
```
We learn that server returns `204` which means item has been successfully deleted. We can no longer see item 1 when we get the item list.

```bash
$ curl -i localhost:9000
```
```http
HTTP/1.1 200 OK
Content-Type: application/com.gojek.todo-list.v1+json
Date: Fri, 27 Jul 2018 14:05:17 GMT
Content-Length: 78

[{"description":"message 18912","id":2},{"description":"topup gopay","id":3}]
```



### Challenges

* Open https://editor.swagger.io, and copy paste your `swagger.yml` there. Play around with the UI.
* From inside the UI, generate client with language of your choice (can be ruby, python, or anything else). Implement a simple application by using the generated client library/sdk to interact with our Golang generated Todo List server.

