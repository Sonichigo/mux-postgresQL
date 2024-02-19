# Product Catelog
A sample url shortener app to test Keploy integration capabilities

## Installation Setup

```bash
git clone https://github.com/keploy/samples-go.git && cd samples-go/mux-sql
go mod download
```

## Installation Keploy

Keploy can be installed on Linux directly and on Windows with the help of WSL. Based on your system archieture, install the keploy latest binary release

1. One Click


### Start Postgres Instance 

Using the docker-compose file we will start our postgres instance:-

```bash
# Start Postgres
docker-compose up -d
```

### Update the Host

> **Since we have setup our sample-app natively set the host to `localhost` on line 10.**

### Capture the Testcases

Now, we will create the binary of our application:-

```zsh
go build
```

Once we have our binary file ready,this command will start the recording of API calls using ebpf:-

```shell
sudo -E keploy record -c "./test-app-product-catelog"
```

Make API Calls using Hoppscotch, Postman or cURL command. Keploy with capture those calls to generate the test-suites containing testcases and data mocks.

#### Generate testcases

To genereate testcases we just need to make some API calls. You can use [Postman](https://www.postman.com/), [Hoppscotch](https://hoppscotch.io/), or simply `curl`

### 1. Generate shortned url

```bash
curl --request POST \
  --url http://localhost:8010/product \
  --header 'content-type: application/json' \
  --data '{
    "name":"Bubbles", 
    "price": 123
}'
```
this will return the response. 
```
{
    "id": 1,
    "name": "Bubbles",
    "price": 123
}
```

#### 2. Redirect to original url from shortened url
1. By using Curl Command
```bash
curl --request GET \
  --url http://localhost:8010/products
```

2. By querying through the browser `http://localhost:8010/products`

Now both these API calls were captured as editable testcases and written to ``keploy/tests folder``. The keploy directory would also have `mocks` files that contains all the outputs of postgres operations. 

![Testcase](./img/testcase.png?raw=true)

Now, let's see the magic! 🪄💫

## Generate Test Runs

Now let's run the test mode (in the mux-sql directory, not the Keploy directory).

```shell
sudo -E keploy test -c "./test-app-product-catelog" --delay 10
```

Once done, you can see the Test Runs on the Keploy server, like this:

![Testrun](./img/testrun.png?raw=true)

So no need to setup fake database/apis like Postgres or write mocks for them. Keploy automatically mocks them and, **The application thinks it's talking to Postgres 😄**

# Using Docker

Keploy can be used on Linux, Windows and MacOS through Docker.

Note: To run Keploy on MacOS through [Docker](https://docs.docker.com/desktop/release-notes/#4252) the version must be ```4.25.2``` or above.

## Create Keploy Alias
To establish a network for your application using Keploy on Docker, follow these steps.

If you're using a docker-compose network, replace keploy-network with your app's `docker_compose_network_name` below.

```shell
alias keploy='sudo docker run --pull always --name keploy-v2 -p 16789:16789 --privileged --pid=host -it -v $(pwd):$(pwd) -w $(pwd) -v /sys/fs/cgroup:/sys/fs/cgroup -v /sys/kernel/debug:/sys/kernel/debug -v /sys/fs/bpf:/sys/fs/bpf -v /var/run/docker.sock:/var/run/docker.sock --rm ghcr.io/keploy/keploy'
```
## Let's start the MongoDB Instance
Using the docker-compose file we will start our mongodb instance:-

```shell
docker-compose up -d
```
> Since we are using docker to run the application, we need to update the `postgres` host on line 10 in `main.go`, update the host to `mux-sql-postgres-1`.
Now, we will create the docker image of our application:-

Now, we will create the docker image of our application:-

```shell
docker build -t mux-app:1.0 .
```

## Capture the Testcases

```zsh
keploy record -c "docker run -p 8010:8010 --name muxSqlApp --network keploy-network mux-app:1.0" --buildDelay 50s
```

![Testcase](./img/testcase.png?raw=true)

### Generate testcases
To genereate testcases we just need to make some API calls. You can use Postman, Hoppscotch, or simply curl

```bash
curl --request POST \
  --url http://localhost:8010/product \
  --header 'content-type: application/json' \
  --data '{
    "name":"coke", 
    "price": 124
}'
```
this will return the response. 

```json
{
    "id": 1,
    "name": "Bubbles",
    "price": 123
}
```

#### 2. Redirect to original url from shortened url
1. By using Curl Command
```bash
curl --request GET \
  --url http://localhost:8010/products
```

2. By querying through the browser `http://localhost:8010/products`

Now both these API calls were captured as editable testcases and written to ``keploy/tests folder``. The keploy directory would also have `mocks` files that contains all the outputs of postgres operations.

## Run the captured testcases
Now that we have our testcase captured, run the test file.

```shell
keploy test -c "sudo docker run -p 8010:8010 --net keploy-network --name muxSqlApp mux-app:1.0" --buildDelay 50s
```
So no need to setup dependencies like mongoDB, web-go locally or write mocks for your testing.

The application thinks it's talking to mongoDB 😄

We will get output something like this:
![Testrun](./img/testrun.png?raw=true)

# Using Keploy with `go-test`

## Installation
Dowload the Keploy go-client SDK, we can use this to even generate realistic mock/stub files for your applications.

```bash
go get -u github.com/keploy/go-sdk/v2
```
## Generating Go-Mocks with unit-test
Keploy can createe **readable/editable** `mocks/stubs` yaml files which can be referenced in any of your unit-tests. Let's add keploy-go-sdk to our existing `main_test.go` file

```go
import(
    "github.com/keploy/go-sdk/v2/keploy"
)

// Inside Our SetMain func
...
err := keploy.New(keploy.Config{
    Mode: keploy.MODE_RECORD,
    Name: "TestMuxSQLApp",
    MuteKeployLogs: false,
    delay: 20,
})
...
```
At the start of each testcases we need to call our `SetMain()`, and at the end of the each test case we need to add the following function which will terminate keploy if not keploy will be running even after unit test is run

```go
keploy.KillProcessOnPort()
```

After making all the necessary changes, we will run `go test` to generate the `Mocks.yml` which will look something like this: - 

```yml
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-0
spec:
    metadata:
        type: config
    postgresrequests:
        - identifier: StartupRequest
          length: 102
          payload: AAAAZgADAAB1c2VyAHBvc3RncmVzAGV4dHJhX2Zsb2F0X2RpZ2l0cwAyAGRhdGFiYXNlAHBvc3RncmVzAGNsaWVudF9lbmNvZGluZwBVVEY4AGRhdGVzdHlsZQBJU08sIE1EWQAA
          startup_message:
            protocolversion: 196608
            parameters:
                client_encoding: UTF8
                database: postgres
                datestyle: ISO, MDY
                extra_float_digits: "2"
                user: postgres
          auth_type: 0
    postgresresponses:
        - header: [R]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 240
                - 226
                - 157
                - 79
          msg_type: 82
          auth_type: 5
    reqtimestampmock: 2024-02-15T07:46:56.70115911Z
    restimestampmock: 2024-02-15T07:46:56.702623523Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-1
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [p]
          identifier: ClientRequest
          length: 102
          password_message:
            password: md5c881c9e7c4396527674db4c8429a4cb1
          msg_type: 112
          auth_type: 0
    postgresresponses:
        - header: [R, S, S, S, S, S, S, S, S, S, S, S, K, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          backend_key_data:
            process_id: 31
            secret_key: 2938221901
          parameter_status:
            - name: application_name
              value: ""
            - name: client_encoding
              value: UTF8
            - name: DateStyle
              value: ISO, MDY
            - name: integer_datetimes
              value: "on"
            - name: IntervalStyle
              value: postgres
            - name: is_superuser
              value: "on"
            - name: server_encoding
              value: UTF8
            - name: server_version
              value: 10.5 (Debian 10.5-2.pgdg90+1)
            - name: session_authorization
              value: postgres
            - name: standard_conforming_strings
              value: "on"
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.713816268Z
    restimestampmock: 2024-02-15T07:46:56.714026464Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-2
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: CREATE TABLE IF NOT EXISTS products ( id SERIAL, name TEXT NOT NULL, price NUMERIC(10,2) NOT NULL DEFAULT 0.00, CONSTRAINT products_pkey PRIMARY KEY (id) )
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: ["N", C, Z]
          identifier: ServerResponse
          length: 102
          payload: TgAAAHVTTk9USUNFAFZOT1RJQ0UAQzQyUDA3AE1yZWxhdGlvbiAicHJvZHVjdHMiIGFscmVhZHkgZXhpc3RzLCBza2lwcGluZwBGcGFyc2VfdXRpbGNtZC5jAEwyMDkAUnRyYW5zZm9ybUNyZWF0ZVN0bXQAAEMAAAARQ1JFQVRFIFRBQkxFAFoAAAAFSQ==
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 67
                - 82
                - 69
                - 65
                - 84
                - 69
                - 32
                - 84
                - 65
                - 66
                - 76
                - 69
          notice_response:
            severity: NOTICE
            severity_unlocalized: NOTICE
            code: 42P07
            message: relation "products" already exists, skipping
            detail: ""
            hint: ""
            position: 0
            internal_position: 0
            internal_query: ""
            where: ""
            schema_name: ""
            table_name: ""
            column_name: ""
            data_type_name: ""
            constraint_name: ""
            file: parse_utilcmd.c
            line: 209
            routine: transformCreateStmt
            unknown_fields: {}
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.716011803Z
    restimestampmock: 2024-02-15T07:46:56.716125213Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-3
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 51
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.722639198Z
    restimestampmock: 2024-02-15T07:46:56.722734651Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-4
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.724505503Z
    restimestampmock: 2024-02-15T07:46:56.724581165Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-5
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.725552815Z
    restimestampmock: 2024-02-15T07:46:56.72563181Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-6
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.727013894Z
    restimestampmock: 2024-02-15T07:46:56.727090098Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-7
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [P, D]
          identifier: ClientRequest
          length: 102
          describe:
            object_type: 83
            name: ""
          parse:
            - name: ""
              query: SELECT id, name, price FROM products LIMIT $1 OFFSET $2
              parameter_oids: []
          msg_type: 68
          auth_type: 0
    postgresresponses:
        - header: ["1", t, T, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          parameter_description:
            parameteroids:
                - 20
                - 20
          ready_for_query:
            txstatus: 73
          row_description: {fields: [{name: [105, 100], table_oid: 16386, table_attribute_number: 1, data_type_oid: 23, data_type_size: 4, type_modifier: -1, format: 0}, {name: [110, 97, 109, 101], table_oid: 16386, table_attribute_number: 2, data_type_oid: 25, data_type_size: -1, type_modifier: -1, format: 0}, {name: [112, 114, 105, 99, 101], table_oid: 16386, table_attribute_number: 3, data_type_oid: 1700, data_type_size: -1, type_modifier: 655366, format: 0}]}
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.728163159Z
    restimestampmock: 2024-02-15T07:46:56.728249487Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-8
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [B, E]
          identifier: ClientRequest
          length: 102
          bind:
            - parameters: [[49, 48], [48]]
              result_format_codes: [1, 0, 0]
          execute:
            - {}
          msg_type: 69
          auth_type: 0
    postgresresponses:
        - header: ["2", C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 83
                - 69
                - 76
                - 69
                - 67
                - 84
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:46:56.729264093Z
    restimestampmock: 2024-02-15T07:46:56.729335088Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-0
spec:
    metadata:
        type: config
    postgresrequests:
        - identifier: StartupRequest
          length: 102
          payload: AAAAZgADAABkYXRhYmFzZQBwb3N0Z3JlcwBjbGllbnRfZW5jb2RpbmcAVVRGOABkYXRlc3R5bGUASVNPLCBNRFkAZXh0cmFfZmxvYXRfZGlnaXRzADIAdXNlcgBwb3N0Z3JlcwAA
          startup_message:
            protocolversion: 196608
            parameters:
                client_encoding: UTF8
                database: postgres
                datestyle: ISO, MDY
                extra_float_digits: "2"
                user: postgres
          auth_type: 0
    postgresresponses:
        - header: [R]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 101
                - 50
                - 209
                - 254
          msg_type: 82
          auth_type: 5
    reqtimestampmock: 2024-02-15T07:47:18.854350212Z
    restimestampmock: 2024-02-15T07:47:18.855718423Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-1
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [p]
          identifier: ClientRequest
          length: 102
          password_message:
            password: md5d16b092aa1ba903b65ba340f77e5a251
          msg_type: 112
          auth_type: 0
    postgresresponses:
        - header: [R, S, S, S, S, S, S, S, S, S, S, S, K, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          backend_key_data:
            process_id: 33
            secret_key: 1982480785
          parameter_status:
            - name: application_name
              value: ""
            - name: client_encoding
              value: UTF8
            - name: DateStyle
              value: ISO, MDY
            - name: integer_datetimes
              value: "on"
            - name: IntervalStyle
              value: postgres
            - name: is_superuser
              value: "on"
            - name: server_encoding
              value: UTF8
            - name: server_version
              value: 10.5 (Debian 10.5-2.pgdg90+1)
            - name: session_authorization
              value: postgres
            - name: standard_conforming_strings
              value: "on"
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.856818816Z
    restimestampmock: 2024-02-15T07:47:18.857135339Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-2
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: CREATE TABLE IF NOT EXISTS products ( id SERIAL, name TEXT NOT NULL, price NUMERIC(10,2) NOT NULL DEFAULT 0.00, CONSTRAINT products_pkey PRIMARY KEY (id) )
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: ["N", C, Z]
          identifier: ServerResponse
          length: 102
          payload: TgAAAHVTTk9USUNFAFZOT1RJQ0UAQzQyUDA3AE1yZWxhdGlvbiAicHJvZHVjdHMiIGFscmVhZHkgZXhpc3RzLCBza2lwcGluZwBGcGFyc2VfdXRpbGNtZC5jAEwyMDkAUnRyYW5zZm9ybUNyZWF0ZVN0bXQAAEMAAAARQ1JFQVRFIFRBQkxFAFoAAAAFSQ==
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 67
                - 82
                - 69
                - 65
                - 84
                - 69
                - 32
                - 84
                - 65
                - 66
                - 76
                - 69
          notice_response:
            severity: NOTICE
            severity_unlocalized: NOTICE
            code: 42P07
            message: relation "products" already exists, skipping
            detail: ""
            hint: ""
            position: 0
            internal_position: 0
            internal_query: ""
            where: ""
            schema_name: ""
            table_name: ""
            column_name: ""
            data_type_name: ""
            constraint_name: ""
            file: parse_utilcmd.c
            line: 209
            routine: transformCreateStmt
            unknown_fields: {}
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.857886127Z
    restimestampmock: 2024-02-15T07:47:18.857968206Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-3
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.859073182Z
    restimestampmock: 2024-02-15T07:47:18.859250713Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-4
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.868834227Z
    restimestampmock: 2024-02-15T07:47:18.868960262Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-5
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.86971555Z
    restimestampmock: 2024-02-15T07:47:18.869898789Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-6
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.872213069Z
    restimestampmock: 2024-02-15T07:47:18.872273982Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-7
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [P, D]
          identifier: ClientRequest
          length: 102
          describe:
            object_type: 83
            name: ""
          parse:
            - name: ""
              query: SELECT name, price FROM products WHERE id=$1
              parameter_oids: []
          msg_type: 68
          auth_type: 0
    postgresresponses:
        - header: ["1", t, T, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          parameter_description:
            parameteroids:
                - 23
          ready_for_query:
            txstatus: 73
          row_description: {fields: [{name: [110, 97, 109, 101], table_oid: 16386, table_attribute_number: 2, data_type_oid: 25, data_type_size: -1, type_modifier: -1, format: 0}, {name: [112, 114, 105, 99, 101], table_oid: 16386, table_attribute_number: 3, data_type_oid: 1700, data_type_size: -1, type_modifier: 655366, format: 0}]}
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.872978648Z
    restimestampmock: 2024-02-15T07:47:18.873076684Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-8
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [B, E]
          identifier: ClientRequest
          length: 102
          bind:
            - parameters: [[49, 49]]
          execute:
            - {}
          msg_type: 69
          auth_type: 0
    postgresresponses:
        - header: ["2", C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 83
                - 69
                - 76
                - 69
                - 67
                - 84
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:18.873905593Z
    restimestampmock: 2024-02-15T07:47:18.873958673Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-0
spec:
    metadata:
        type: config
    postgresrequests:
        - identifier: StartupRequest
          length: 102
          payload: AAAAZgADAABleHRyYV9mbG9hdF9kaWdpdHMAMgBkYXRlc3R5bGUASVNPLCBNRFkAdXNlcgBwb3N0Z3JlcwBkYXRhYmFzZQBwb3N0Z3JlcwBjbGllbnRfZW5jb2RpbmcAVVRGOAAA
          startup_message:
            protocolversion: 196608
            parameters:
                client_encoding: UTF8
                database: postgres
                datestyle: ISO, MDY
                extra_float_digits: "2"
                user: postgres
          auth_type: 0
    postgresresponses:
        - header: [R]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 79
                - 166
                - 144
                - 102
          msg_type: 82
          auth_type: 5
    reqtimestampmock: 2024-02-15T07:47:41.002401312Z
    restimestampmock: 2024-02-15T07:47:41.003831269Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-1
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [p]
          identifier: ClientRequest
          length: 102
          password_message:
            password: md595aa96b9299735c879b04697709f199a
          msg_type: 112
          auth_type: 0
    postgresresponses:
        - header: [R, S, S, S, S, S, S, S, S, S, S, S, K, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          backend_key_data:
            process_id: 34
            secret_key: 286020656
          parameter_status:
            - name: application_name
              value: ""
            - name: client_encoding
              value: UTF8
            - name: DateStyle
              value: ISO, MDY
            - name: integer_datetimes
              value: "on"
            - name: IntervalStyle
              value: postgres
            - name: is_superuser
              value: "on"
            - name: server_encoding
              value: UTF8
            - name: server_version
              value: 10.5 (Debian 10.5-2.pgdg90+1)
            - name: session_authorization
              value: postgres
            - name: standard_conforming_strings
              value: "on"
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.004914122Z
    restimestampmock: 2024-02-15T07:47:41.005138192Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-2
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: CREATE TABLE IF NOT EXISTS products ( id SERIAL, name TEXT NOT NULL, price NUMERIC(10,2) NOT NULL DEFAULT 0.00, CONSTRAINT products_pkey PRIMARY KEY (id) )
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: ["N", C, Z]
          identifier: ServerResponse
          length: 102
          payload: TgAAAHVTTk9USUNFAFZOT1RJQ0UAQzQyUDA3AE1yZWxhdGlvbiAicHJvZHVjdHMiIGFscmVhZHkgZXhpc3RzLCBza2lwcGluZwBGcGFyc2VfdXRpbGNtZC5jAEwyMDkAUnRyYW5zZm9ybUNyZWF0ZVN0bXQAAEMAAAARQ1JFQVRFIFRBQkxFAFoAAAAFSQ==
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 67
                - 82
                - 69
                - 65
                - 84
                - 69
                - 32
                - 84
                - 65
                - 66
                - 76
                - 69
          notice_response:
            severity: NOTICE
            severity_unlocalized: NOTICE
            code: 42P07
            message: relation "products" already exists, skipping
            detail: ""
            hint: ""
            position: 0
            internal_position: 0
            internal_query: ""
            where: ""
            schema_name: ""
            table_name: ""
            column_name: ""
            data_type_name: ""
            constraint_name: ""
            file: parse_utilcmd.c
            line: 209
            routine: transformCreateStmt
            unknown_fields: {}
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.006061763Z
    restimestampmock: 2024-02-15T07:47:41.006143217Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-3
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.006763222Z
    restimestampmock: 2024-02-15T07:47:41.006815927Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-4
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.014207493Z
    restimestampmock: 2024-02-15T07:47:41.014300112Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-5
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.014978781Z
    restimestampmock: 2024-02-15T07:47:41.015027653Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-6
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.016222166Z
    restimestampmock: 2024-02-15T07:47:41.01628837Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-7
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [P, D]
          identifier: ClientRequest
          length: 102
          describe:
            object_type: 83
            name: ""
          parse:
            - name: ""
              query: INSERT INTO products(name, price) VALUES($1, $2) RETURNING id
              parameter_oids: []
          msg_type: 68
          auth_type: 0
    postgresresponses:
        - header: ["1", t, T, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          parameter_description:
            parameteroids:
                - 25
                - 1700
          ready_for_query:
            txstatus: 73
          row_description: {fields: [{name: [105, 100], table_oid: 16386, table_attribute_number: 1, data_type_oid: 23, data_type_size: 4, type_modifier: -1, format: 0}]}
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.016984913Z
    restimestampmock: 2024-02-15T07:47:41.01708124Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-8
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [B, E]
          identifier: ClientRequest
          length: 102
          bind:
            - parameters: [[116, 101, 115, 116, 32, 112, 114, 111, 100, 117, 99, 116], [49, 49, 46, 50, 50]]
              result_format_codes: [1]
          execute:
            - {}
          msg_type: 69
          auth_type: 0
    postgresresponses:
        - header: ["2", D, C, Z]
          identifier: ServerResponse
          length: 102
          payload: MgAAAAREAAAADgABAAAABAAAAAFDAAAAD0lOU0VSVCAwIDEAWgAAAAVJ
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 73
                - 78
                - 83
                - 69
                - 82
                - 84
                - 32
                - 48
                - 32
                - 49
          data_row: [{row_values: ['base64:AAAAAQ==']}, {row_values: ['base64:AAAAAQ==']}, {row_values: ['base64:AAAAAQ==']}]
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:47:41.018185884Z
    restimestampmock: 2024-02-15T07:47:41.018649648Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-0
spec:
    metadata:
        type: config
    postgresrequests:
        - identifier: StartupRequest
          length: 102
          payload: AAAAZgADAABkYXRhYmFzZQBwb3N0Z3JlcwBjbGllbnRfZW5jb2RpbmcAVVRGOABleHRyYV9mbG9hdF9kaWdpdHMAMgB1c2VyAHBvc3RncmVzAGRhdGVzdHlsZQBJU08sIE1EWQAA
          startup_message:
            protocolversion: 196608
            parameters:
                client_encoding: UTF8
                database: postgres
                datestyle: ISO, MDY
                extra_float_digits: "2"
                user: postgres
          auth_type: 0
    postgresresponses:
        - header: [R]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 60
                - 57
                - 10
                - 217
          msg_type: 82
          auth_type: 5
    reqtimestampmock: 2024-02-15T07:48:03.132631653Z
    restimestampmock: 2024-02-15T07:48:03.133999823Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-1
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [p]
          identifier: ClientRequest
          length: 102
          password_message:
            password: md55c98cbf2933503fe5e4df30a02236d03
          msg_type: 112
          auth_type: 0
    postgresresponses:
        - header: [R, S, S, S, S, S, S, S, S, S, S, S, K, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          backend_key_data:
            process_id: 35
            secret_key: 2701389984
          parameter_status:
            - name: application_name
              value: ""
            - name: client_encoding
              value: UTF8
            - name: DateStyle
              value: ISO, MDY
            - name: integer_datetimes
              value: "on"
            - name: IntervalStyle
              value: postgres
            - name: is_superuser
              value: "on"
            - name: server_encoding
              value: UTF8
            - name: server_version
              value: 10.5 (Debian 10.5-2.pgdg90+1)
            - name: session_authorization
              value: postgres
            - name: standard_conforming_strings
              value: "on"
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
            - name: TimeZone
              value: UTC
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.13499989Z
    restimestampmock: 2024-02-15T07:48:03.135317038Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-2
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: CREATE TABLE IF NOT EXISTS products ( id SERIAL, name TEXT NOT NULL, price NUMERIC(10,2) NOT NULL DEFAULT 0.00, CONSTRAINT products_pkey PRIMARY KEY (id) )
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: ["N", C, Z]
          identifier: ServerResponse
          length: 102
          payload: TgAAAHVTTk9USUNFAFZOT1RJQ0UAQzQyUDA3AE1yZWxhdGlvbiAicHJvZHVjdHMiIGFscmVhZHkgZXhpc3RzLCBza2lwcGluZwBGcGFyc2VfdXRpbGNtZC5jAEwyMDkAUnRyYW5zZm9ybUNyZWF0ZVN0bXQAAEMAAAARQ1JFQVRFIFRBQkxFAFoAAAAFSQ==
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 67
                - 82
                - 69
                - 65
                - 84
                - 69
                - 32
                - 84
                - 65
                - 66
                - 76
                - 69
          notice_response:
            severity: NOTICE
            severity_unlocalized: NOTICE
            code: 42P07
            message: relation "products" already exists, skipping
            detail: ""
            hint: ""
            position: 0
            internal_position: 0
            internal_query: ""
            where: ""
            schema_name: ""
            table_name: ""
            column_name: ""
            data_type_name: ""
            constraint_name: ""
            file: parse_utilcmd.c
            line: 209
            routine: transformCreateStmt
            unknown_fields: {}
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.136151906Z
    restimestampmock: 2024-02-15T07:48:03.136267066Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-3
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 49
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.137687609Z
    restimestampmock: 2024-02-15T07:48:03.137762729Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-4
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.139305931Z
    restimestampmock: 2024-02-15T07:48:03.13938076Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-5
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: DELETE FROM products
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 68
                - 69
                - 76
                - 69
                - 84
                - 69
                - 32
                - 48
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.139869648Z
    restimestampmock: 2024-02-15T07:48:03.140005765Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-6
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [Q]
          identifier: ClientRequest
          length: 102
          query:
            string: ALTER SEQUENCE products_id_seq RESTART WITH 1
          msg_type: 81
          auth_type: 0
    postgresresponses:
        - header: [C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 65
                - 76
                - 84
                - 69
                - 82
                - 32
                - 83
                - 69
                - 81
                - 85
                - 69
                - 78
                - 67
                - 69
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.141357103Z
    restimestampmock: 2024-02-15T07:48:03.141466805Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-7
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [P, D]
          identifier: ClientRequest
          length: 102
          describe:
            object_type: 83
            name: ""
          parse:
            - name: ""
              query: INSERT INTO products(name, price) VALUES($1, $2)
              parameter_oids: []
          msg_type: 68
          auth_type: 0
    postgresresponses:
        - header: ["1", t, "n", Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          parameter_description:
            parameteroids:
                - 25
                - 1700
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.141999066Z
    restimestampmock: 2024-02-15T07:48:03.142197304Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-8
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [B, E]
          identifier: ClientRequest
          length: 102
          bind:
            - parameters: [[80, 114, 111, 100, 117, 99, 116, 32, 48], [49, 48]]
          execute:
            - {}
          msg_type: 69
          auth_type: 0
    postgresresponses:
        - header: ["2", C, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 73
                - 78
                - 83
                - 69
                - 82
                - 84
                - 32
                - 48
                - 32
                - 49
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.143742797Z
    restimestampmock: 2024-02-15T07:48:03.143895497Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-9
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [P, D]
          identifier: ClientRequest
          length: 102
          describe:
            object_type: 83
            name: ""
          parse:
            - name: ""
              query: SELECT name, price FROM products WHERE id=$1
              parameter_oids: []
          msg_type: 68
          auth_type: 0
    postgresresponses:
        - header: ["1", t, T, Z]
          identifier: ServerResponse
          length: 102
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          parameter_description:
            parameteroids:
                - 23
          ready_for_query:
            txstatus: 73
          row_description: {fields: [{name: [110, 97, 109, 101], table_oid: 16386, table_attribute_number: 2, data_type_oid: 25, data_type_size: -1, type_modifier: -1, format: 0}, {name: [112, 114, 105, 99, 101], table_oid: 16386, table_attribute_number: 3, data_type_oid: 1700, data_type_size: -1, type_modifier: 655366, format: 0}]}
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.144801111Z
    restimestampmock: 2024-02-15T07:48:03.144945602Z
---
version: api.keploy.io/v1beta1
kind: Postgres
name: mock-10
spec:
    metadata:
        type: config
    postgresrequests:
        - header: [B, E]
          identifier: ClientRequest
          length: 102
          bind:
            - parameters: [[49]]
          execute:
            - {}
          msg_type: 69
          auth_type: 0
    postgresresponses:
        - header: ["2", D, C, Z]
          identifier: ServerResponse
          length: 102
          payload: MgAAAAREAAAAHAACAAAACVByb2R1Y3QgMAAAAAUxMC4wMEMAAAANU0VMRUNUIDEAWgAAAAVJ
          authentication_md5_password:
            salt:
                - 0
                - 0
                - 0
                - 0
          command_complete:
            - command_tag:
                - 83
                - 69
                - 76
                - 69
                - 67
                - 84
                - 32
                - 49
          data_row: [{row_values: [Product 0, "10.00"]}, {row_values: [Product 0, "10.00"]}, {row_values: [Product 0, "10.00"]}]
          ready_for_query:
            txstatus: 73
          msg_type: 90
          auth_type: 0
    reqtimestampmock: 2024-02-15T07:48:03.145444782Z
    restimestampmock: 2024-02-15T07:48:03.14581301Z
```

Voila!! We have our mock ready to be used in any of our other unit-tests files.

<!-- ## Getting Code Coverage with Keploy --!>