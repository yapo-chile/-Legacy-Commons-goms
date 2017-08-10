# goms

## How to run the service

* Create the dir: `~/go/src/github.schibsted.io/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```
  $ cd ~/go/src/github.schibsted.io/Yapo
  $ git clone git@github.schibsted.io:Yapo/goms.git
  ```
* You will need to modify this files and you should change every reference to goms on the imported packages with the name of your service/api
	- api.go
	- Makefile: same as above
	- scripts/api.spec
* You should need to change, add or replace the example endpoint and any of its routes, these can be found in these files
	- sources/handlers.go
	- sources/router.go

* On the top dir execute the make instruction to clean and start:

  ```
  $ cd goms
  $ make clean_start
  ```

* How to check the service?

  ```
  $ make status
  SERVICE RUNNING (PID: 7568)
  make db-status
  make[1]: Entering directory   `/home/user/go/src/github.schibsted.io/Yapo/goms'
  DATABASE RUNNING
  make[1]: Leaving directory `/home/user/go/src/github.schibsted.io/Yapo/goms'
  ```

* How to get more info?

  ```
  $ make info
  YO           : user
  ServerRoot   : /home/user/go/src/github.schibsted.io/Yapo/goms
  DocumentRoot : /home/user/go/src/github.schibsted.io/Yapo/goms/src/public
  API Base URL : http://SERVER:PORT
  DB connect   : psql -h localhost -p PORT goms-db
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
