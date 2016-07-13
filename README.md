# smarty

This is a go package for working with the SmartyStreets US Street Address API. 

To install this package you can run:<br>
`go get github.com/cbelk/smarty`

The package will verify that the minimum information required by the API has been provided.
The minimum requirements are:
  * AuthID
  * AuthToken
  * Either:
    * Street, City, State
        <br>OR
    * Street, ZipCode
        <br>OR
    * Street (freeform)

The package supports both the GET and POST requests and returns an array of SmartResponse objects.

An example program using GET:
*Note: If you don't want to set any optional arguments you can pass in nil for the SmartRequestOptional parameter.*
```
import (
    "fmt"
    "github.com/cbelk/smarty"
    "log"
)

func main () {
    req := new(smarty.SmartRequest)
    req.AuthID = "your-auth-id-here"
    req.AuthToken = "your token"
    req.Street = "123 Some St."
    req.City = "Heresville"
    req.State = "CA"

    res, err := smarty.GetAddress(req, nil)
    if err != nil {
        log.Fatal(err)
    }

    smart, err := smarty.ParseResponse(res)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Smart: %v\n", smart)
}
```
With the POST request you need to supply a slice of SmartRequests and a slice of SmartRequestOptions
where element 'i' in SmartRequestOptions goes with element 'i' in SmartRequests. If you have no 
options, you can supply a nil value for the SmartRequestOptions slice, but if you use
SmartRequestOptions then you need to provide a corresponding SmartRequestOptions object for each
element in the SmartRequests slice even if it's just an empty object. Also, the AuthID and AuthToken
need to be provided with at least one SmartRequests.

An example program using POST:
```
import (
    "fmt"
    "github.com/cbelk/smarty"
    "log"
)

func main () {
    var reqs []*smarty.SmartRequest
    var reqOps []*smarty.SmartRequestOptional

    req1 := new(smarty.SmartRequest)
    req1.AuthID = "your-auth-id-here"
    req1.AuthToken = "your token"
    req1.Street = "123 Some St."
    req1.City = "Heresville"
    req1.State = "CA"
    reqOp1 := new(smarty.SmartRequestOptional)
    reqOp1.InputID = "addCA"

    req2 := new(smarty.SmartRequest)
    req2.Street = "456 This Drive"
    req2.Zipcode = "98765"
    reqOp2 := new(smarty.SmartRequestOptional)

    req3 := new(smarty.SmartRequest)
    req3.FreeForm = "789 My Rd Theresville, WY 01234"
    reqOp3 := new(smarty.SmartRequestOptional)
    reqOp3.Street2 = "leave on porch"
    reqOp3.InputID = "shippingAddress"

    reqs = append(reqs, req1)
    reqs = append(reqs, req2)
    reqs = append(reqs, req3)

    reqOps = append(reqOps, reqOp1)
    reqOps = append(reqOps, reqOp2)
    reqOps = append(reqOps, reqOp3)

    res, err := smarty.PostAddress(reqs, reqOps)
    if err != nil {
        log.Fatal(err)
    }

    smart, err := smarty.ParseResponse(res)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Smart: %v\n", smart)
}
```

For a description of the input fields or questions about the API refer to the SmartyStreets documentation:<br>
[Smarty Streets US Address API](https://smartystreets.com/docs/us-street-api)
