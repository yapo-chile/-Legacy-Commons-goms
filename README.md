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
    $ go build
  ```

* How to run the tests

  ```
	$ make test
  ```

* How to check format

  ```
	$ make check -s
  ```


## Endpoints
* GET  /api/v1/theendpoint

## Error Codes
     **ERROR**          |     **HTTP Code**    	|	**MESSAGE**
---------------------   | ---------------------	| ---------------------
InvalidParams		|	410		|	INVALID_PARAMS

### Contact
dev@schibsted.cl
