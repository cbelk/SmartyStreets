package smarty

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "errors"
)

const (
    baseURL string = "https://api.smartystreets.com/street-address"
)

// SmartRequest holds the required request data, although some of this data is optional depending on
// other fields being set or not *(see 'Input Fields' section of the  SmartyStreets US street api)*
type SmartRequest struct {
    AuthID       string
    AuthToken    string
    Street       string
    City         string
    State        string
    Zipcode      string
    FreeForm     string  //Entire address stored in this field (NO country info)
    Candidates   int     //Max number of results (MAX 10)
}

// SmartRequestOptional holds the optional request data.
type SmartRequestOptional struct {
    Addressee    string  //Name of recipient, firm or company
    InputID      string  //Unique id that gets copied into output
    Lastline     string  //City, State and ZipCode combined
    Secondary    string  //Apartment, suite, or office number
    Street2      string  //Extra info (eg leave on porch)
    Urbanization string  //Only used with Puerto Rico
}

// JSONData will hold the data to be passed in as the json object in the POST request.
// Currently using this insted of merging the passed in SmartRequest and SmartRequestOptional
// to avoid confilct since street and freeform will populate the same json object field (street:)
type JSONData struct {
    Street       string `json:"street"`
    City         string `json:"city"`
    State        string `json:"state"`
    Zipcode      string `json:"zipcode"`
    Addressee    string `json:"addressee"`
    InputID      string `json:"input_id"`
    Lastline     string `json:"lastline"`
    Secondary    string `json:"secondary"`
    Street2      string `json:"street2"`
    Urbanization string `json:"urbanization"`
    Candidates   int    `json:"candidates"`
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

// PostAddress is used to construct the POST request from the slices of struct references. If there are
// no request options, then a nil value can be passed for reqOps, but if even if just one reqOp is needed,
// then reqOps must contain an entry in its slice to go with the corresponding entry in the reqs slice (even
// if it's just a nil value). Therefore, reqOps can be nil OR has to be the same length as reqs. Authentication
// needs to be provided in at least one req object in the reqs slice. After constructing the request data, the
// appropriate headers are added and the request is POSTed, returning the result.
func PostAddress(reqs []*SmartRequest, reqOps []*SmartRequestOptional) (res *http.Response, err error) {
    if reqs != nil {
        equal := len(reqs) == len(reqOps)
        omitted := reqOps == nil
        if equal || omitted {
            var nurl string
            var data []JSONData
            authed := false
            for k := range reqs {
                if hasAuth(reqs[k]) {
                    authed = true
                    nurl = baseURL + fmt.Sprintf("?auth-id=%s&auth-token=%s", reqs[k].AuthID, reqs[k].AuthToken)
                }
                tmp, e := preparePostData(reqs[k])
                if e == nil {
                    if !omitted {
                        addReqOpData(reqOps[k], tmp)
                    }
                    data = append(data, *tmp)
                } else {
                    err = e
                }
            }
            if authed {
                b, e := json.Marshal(data)
                if e != nil {
                }
                hreq, e := http.NewRequest("POST", nurl, bytes.NewBuffer(b))
                if e == nil {
                    hreq.Header.Set("Content-Type", "application/json")
                    hreq.Header.Set("Host", "api.smartystreets.com")
                    client := &http.Client{}
                    res, e = client.Do(hreq)
                    if e != nil {
                        err = e
                    }
                } else {
                    err = e
                }
            } else {
                err = errors.New("Authentication Required")
            }
        } else {
            err = errors.New("Lengths of []SmartRequest and []SmartRequestOptional should be equal OR []SmartRequestOptional should be nil")
        }
    }
    return
}

// addReqData is used to pack the fields from the SmartRequest object into
// the JSONData object.
func addReqData(req *SmartRequest, data *JSONData) {
    if req.FreeForm != "" {
        data.Street = req.FreeForm
    }
    if req.Street != "" {
        data.Street = req.Street
    }
    if req.City != "" {
        data.City = req.City
    }
    if req.State != "" {
        data.State = req.State
    }
    if req.Zipcode != "" {
        data.Zipcode = req.Zipcode
    }
    if req.Candidates != 0 {
        data.Candidates = req.Candidates
    }
}

// addReqOpData is used to pack the fields from the SmartRequestOptional object into
// the JSONData object.
func addReqOpData(reqOp *SmartRequestOptional, data *JSONData) {
    if reqOp.Addressee != "" {
        data.Addressee = reqOp.Addressee
    }
    if reqOp.InputID != "" {
        data.InputID = reqOp.InputID
    }
    if reqOp.Lastline != "" {
        data.Lastline = reqOp.Lastline
    }
    if reqOp.Secondary != "" {
        data.Secondary = reqOp.Secondary
    }
    if reqOp.Street2 != "" {
        data.Street2 = reqOp.Street2
    }
    if reqOp.Urbanization != "" {
        data.Urbanization = reqOp.Urbanization
    }
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

// preparePostData drives the creation of the json object to be submitted with the
// POST request.
func preparePostData(req *SmartRequest) (*JSONData, error) {
    var data JSONData
    var err error
    if req != nil {
        if validReq(req) {
            addReqData(req, &data)
        } else {
            err = errors.New("Invalid Request")
        }
    } else {
        err = errors.New("SmartRequest cannot be nil")
    }
    return &data, err
}

// prepareReqQuery constructs and returns the query string from the SmartRequest based on the rules defined in
// the 'Input Fields' section of the US street api documentation.
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
// the 'Input Fields' section of the US street api documentation.
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

// validReq is used to check the validity of the submitted request. A valid request must contain at least
// either street + city + state OR street + zipcode OR freeform input per the SmartyStreets US address api.
func validReq(req *SmartRequest) bool {
    if req.Street != "" {
        if (req.City != "" && req.State != "") || req.Zipcode != "" {
            return true
        } else if req.FreeForm != "" {
            return true
        }
    } else if req.FreeForm != "" {
        return true
    }
    return false
}
