# goms

## How to run the service

* Create the dir: `~/go/src/github.schibsted.io/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.schibsted.io/Yapo
  $ git clone git@github.schibsted.io:Yapo/goms.git
  ```

* You will need to modify these files and you should change every reference to goms on the imported packages with the name of your service/api
	- api.go
	- Makefile: same as above
	- scripts/api.spec
* You should need to change, add or replace the example endpoint and any of its routes, these can be found in these files
	- service/handlers.go
	- service/router.go

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd goms
  $ make start
  ```

* How to check the service?

  ```
  $ make status
  SERVICE RUNNING (PID: 7568)
  ```

* How to get more info?

  ```
  $ make info
  YO           : user
  ServerRoot   : /home/user/go/src/github.schibsted.io/Yapo/goms
  DocumentRoot : /home/user/go/src/github.schibsted.io/Yapo/goms/src/public
  API Base URL : http://SERVER:PORT
  ```

* If you change the code:

  ```
  $ make start
  ```

* How to run the tests

  ```
  $ make test
  ```

* How to check format

  ```
  $ make check -s
  ```

* Update to your own service:
  - Create a repo for your new service on: https://github.schibsted.io/Yapo
  - Rename your goms dir to your service name:
  ```
    $ mv goms YourService
  ```
  - Update origin: 
  ```
  # https://help.github.com/articles/changing-a-remote-s-url/
  $ git remote set-url origin git@github.schibsted.io:Yapo/YourService.git
  ```

## Endpoints
### GET  /api/v1/healthcheck
Reports whether the service is up and ready to respond.

> When implementing a new service, you MUST keep this endpoint
and update it so it replies according to your service status!

#### Request
No request parameters

#### Response
* Status: Ok message, representing service health

```javascript
200 OK
{
	"Status": "OK"
}
```

### GET  /api/v1/inject
Exhibits dependency injection feature

#### Request
No request parameters

#### Response
* Resource: The resource being injected
* Resource.Name: Name of the resource
* Resource.Usage: Intended usage

```javascript
200 OK
{
	"Resource": {
		"Name": "A Resource",
		"Usage": "Being injected"
	}
}
```
## Error Codes
Please update this table with the error codes you use.
Keep it as http standard compliant as possible.

|      **ERROR**         |     **HTTP Code**    	|      **MESSAGE**
|----------------------- | --------------------- | ---------------------
|Unprocessable Entity    |         422           |  JSON motherflower, do you speak it?

### Contact
dev@schibsted.cl
