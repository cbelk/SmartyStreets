package smarty

import (
    "fmt"
    "net/http"
    "net/url"
    "errors"
)

const (
    baseURL string = "https://api.smartystreets.com/street-address"
)

// SmartRequest holds the required request data, although some of this data is optional depending on
// other fields being set or not *(see 'Input Fields' section of the  SmartyStreets us street api)*
type SmartRequest struct {
    AuthID      string
    AuthToken   string
    Street      string
    City        string
    State       string
    Zipcode     string
    FreeForm    string  //Entire address stored in this field (NO country info)
    Candidates  int     //Max number of results (MAX 10)
}

// SmartRequestOptional holds the optional request data.
type SmartRequestOptional struct {
    Addressee   string  //Name of recipient, firm or company
    InputID     string  //Unique id that gets copied into output
    Lastline    string  //City, State and ZipCode combined
    Secondary   string  //Apartment, suite, or office number
    Street2     string  //Extra info (eg leave on porch)
}

// GetAddress is used to construct the GET request from the referenced request structs. It then sends
// the request to the SmartyStreets api and returns the response and/or an error. The optional request
// parameters can be omitted by passing in nil for reqOp. Note: The req paramter cannot take a nil value
// since it holds the authentication info.
func GetAddress(req *SmartRequest, reqOp *SmartRequestOptional) (res *http.Response, err error) {
    if req != nil {
        if hasAuth(req) {
            nurl := baseURL + fmt.Sprintf("?auth-id=%s&auth-token=%s", req.AuthID, req.AuthToken)
            query, e := prepareReqQuery(req)
            if e == nil {
                nurl += query
                appendCandidates(req, &nurl)
                if reqOp != nil {
                    nurl += prepareReqOpQuery(reqOp)
                }
                res, e = http.Get(nurl)
                if e != nil {
                    err = e
                }
            } else {
                err = e
            }
        } else {
            err = errors.New("Authentication paramaters required")
        }
    }
    return
}

// appendCandidates is used to determine whether the candidates value has been set and add the appropriate
// value to the query string. The value is confined to the range [1-10] with a default value of 1 per the 
// api specifications.
func appendCandidates(req *SmartRequest, nurl *string) {
    if req.Candidates > 0 && req.Candidates <= 10 {
        *nurl += fmt.Sprintf("&candidates=%d", req.Candidates)
    } else if req.Candidates > 10 {
        *nurl += "&candidates=10"
    } else {
        *nurl += "&candidates=1"
    }
}

// hasAuth is used to determine whether the authentication info has been set.
func hasAuth(req *SmartRequest) bool {
    if req.AuthID != "" && req.AuthToken != "" {
        return true
    }
    return false
}

// prepareReqQuery constructs and returns the query string from the SmartRequest based on the rules defined in
// the 'Input Fields' section of the api documentation.
func prepareReqQuery(req *SmartRequest) (query string, err error) {
    if req.Street != "" {
        query += fmt.Sprintf("&street=%s", url.QueryEscape(req.Street))
        if req.City != "" && req.State != "" {
            query += fmt.Sprintf("&city=%s&state=%s", url.QueryEscape(req.City), url.QueryEscape(req.State))
            if req.Zipcode != "" {
                query += fmt.Sprintf("&zipcode=%s", url.QueryEscape(req.Zipcode))
            }
        } else if req.Zipcode != "" {
            query += fmt.Sprintf("&zipcode=%s", url.QueryEscape(req.Zipcode))
        } else {
            err = errors.New("Either street + city + state OR street + zipcode required if not using freeform addressing")
        }
    } else if req.FreeForm != "" {
        query += fmt.Sprintf("&street=%s", url.QueryEscape(req.FreeForm))
    } else {
        err = errors.New("Street address OR freeform required")
    }
    return
}

// prepareReqOpQuery constructs and returns the query string from the SmartRequestOptions based on the rules defined in
// the 'Input Fields' section of the api documentation.
func prepareReqOpQuery(reqOp *SmartRequestOptional) (query string) {
    if reqOp.Addressee != "" {
        query += "&addressee=" + url.QueryEscape(reqOp.Addressee)
    }
    if reqOp.InputID != "" {
        query += "&input_id=" + url.QueryEscape(reqOp.InputID)
    }
    if reqOp.Lastline != "" {
        query += "&lastline=" + url.QueryEscape(reqOp.Lastline)
    }
    if reqOp.Secondary != "" {
        query += "&secondary=" + url.QueryEscape(reqOp.Secondary)
    }
    if reqOp.Street2 != "" {
        query += "&street2=" + url.QueryEscape(reqOp.Street2)
    }
    return
}
