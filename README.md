# receipt_processor

A stateless RESTful API implemented in Go
=========================================

This project implements an API given [here](https://github.com/fetch-rewards/receipt-processor-challenge). The API has two endpoints. In short: 
* When a properly structured JSON object, representing a commercial receipt, is POSTed to `/receipts/process`, an ID is created by hashing the JSON object, a point score is calculated for the receipt, and the ID is returned as a JSON object. 
* When a GET request is sent to `/receipts/{id}/points`, given a valid ID, the point score for that ID is returned as a JSON object. 


Dependencies
------------

A user should have [Go](https://go.dev/doc/install) installed on their system (this project has been tested with version 1.19.6). Only the standard libraries provided with the Go installation are used. 


Running the project
-------------------

That points are being correctly computed for the examples given at the [documentation](https://github.com/fetch-rewards/receipt-processor-challenge) can be verified with a unit test: 

    go test

should return `PASS`. 

To start the application on port 8080, navigate to the root directory of the project and run 

`go run ./main.go`

Various examples of receipt data are provided in the `examples` folder: 

    curl -d @examples/simple-receipt.json http://localhost:8080/receipts/process
    { "id" : "9adae37a4510d62bb7e3679d4c098ca84d5cdf308a866626c914542eab4f46bc" }

    curl localhost:8080/receipts/9adae37a4510d62bb7e3679d4c098ca84d5cdf308a866626c914542eab4f46bc/points
    { "points" : 31 }

    curl -d @examples/morning-receipt.json http://localhost:8080/receipts/process
    { "id" : "a7624a895e6f511fb6376412e3c0f1a42ba57f9acae38f8fd808550c024db5ac" }

    curl http://localhost:8080/receipts/a7624a895e6f511fb6376412e3c0f1a42ba57f9acae38f8fd808550c024db5ac/points
    { "points" : 15 }
