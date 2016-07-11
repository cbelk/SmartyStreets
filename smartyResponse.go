package smarty

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
)

// SmartResponse holds the fields the Root according to the SmartyStreets US street api output field definitions.
type SmartResponse struct {
    InputID                     string          `json:"input_id"`
    InputIndex                  int             `json:"input_index"`
    CandidateIndex              int             `json:"candidate_index"`
    Addressee                   string          `json:"addressee"`
    DeliveryLine1               string          `json:"delivery_line_1"`
    DeliveryLine2               string          `json:"delivery_line_2"`
    LastLine                    string          `json:"last_line"`
    DeliveryPointBarcode        string          `json:"delivery_point_barcode"`
    Components                  SmartComponents `json:"components"`
    Metadata                    SmartMetaData   `json:"metadata"`
    Analysis                    SmartAnalysis   `json:"analysis"`
}

// SmartComponents holds the fields the Components according to the SmartyStreets US street api output field definitions.
type SmartComponents struct {
    Urbanization                string
    PrimaryNumber               string          `json:"primary_number"`
    StreetName                  string          `json:"street_name"`
    StreetPredirection          string          `json:"street_predirection"`
    StreetPostdirection         string          `json:"street_postdirection"`
    StreetSuffix                string          `json:"street_suffix"`
    SecondaryNumber             string          `json:"secondary_number"`
    SecondaryDesignator         string          `json:"secondary_designator"`
    ExtraSecondaryNumber        string          `json:"extra_secondary_number"`
    ExtraSecondaryDesignator    string          `json:"extra_secondary_designator"`
    PMBdesignator               string          `json:"pmb_designator"`
    PMBnumber                   string          `json:"pmb_number"`
    CityName                    string          `json:"city_name"`
    DefaultCityName             string          `json:"default_city_name"`
    StateAbbreviation           string          `json:"state_abbreviation"`
    Zipcode                     string          `json:"zipcode"`
    Plus4Code                   string          `json:"plus4_code"`
    DeliveryPoint               string          `json:"delivery_point"`
    DeliveryPointCheckDigit     string          `json:"delivery_point_check_digit"`
}

// SmartMetaData holds the fields the Metadata according to the SmartyStreets US street api output field definitions.
type SmartMetaData struct {
    RecordType                  string          `json:"record_type"`
    Ziptype                     string          `json:"zip_type"`
    CountyFips                  string          `json:"county_fips"`
    CountyName                  string          `json:"county_name"`
    CarrierRoute                string          `json:"carrier_route"`
    CongressionalDistrict       string          `json:"congressional_district"`
    BuildingDefaultIndicator    string          `json:"building_default_indicator"`
    RDI                         string          `json:"rdi"`
    ElotSequence                string          `json:"elot_sequence"`
    ElotSort                    string          `json:"elot_sort"`
    Latitude                    float64         `json:"latitude"`
    Longitude                   float64         `json:"longitude"`
    Precision                   string          `json:"precision"`
    TimeZone                    string          `json:"time_zone"`
    UTCoffset                   float32         `json:"utc_offset"`
    DST                         bool            `json:"dst"`
}

// SmartAnalysis holds the fields the Analysis according to the SmartyStreets US street api output field definitions.
type SmartAnalysis struct {
    DPVmatchCode                string          `json:"dpv_match_code"`
    DPVfootnotes                string          `json:"dpv_footnotes"`
    DPVcmra                     string          `json:"dpv_cmra"`
    DPVvacant                   string          `json:"dpv_vacant"`
    Active                      string          `json:"active"`
    EWSmatch                    string          `json:"ews_match"`
    Footnotes                   string          `json:"footnotes"`
    LACSLinkCode                string          `json:"lackslink_code"`
    LACSLinkIndicator           string          `json:"lacslink_indicator"`
    SuiteLinkMatch              string          `json:"suitelink_match"`
}

// ParseResponse is used to check the status code and unpack the returned json (assuming 200 statuscode) into
// a SmartResponse object.
func ParseResponse(res *http.Response) ([]SmartResponse, error) {
    var smart []SmartResponse
    var err error
    if res.StatusCode == 200 {
//        json.NewDecoder(res.Body).Decode(&smart)
        b, e := ioutil.ReadAll(res.Body)
        if e != nil {
            err = e
        } else {
            e = json.Unmarshal(b, &smart)
            if e != nil {
                err = e
            }
        }
   } else {
        err = getStatusError(res.StatusCode)
   }
   return smart, err
}

// getStatusError is used to check the status code and produce an error according to the SmartyStreets US
// street api status codes and results section.
func getStatusError(statusCode int) (err error) {
    if statusCode == 401 {
        err = errors.New("Unauthorized")
    }
    if statusCode == 402 {
        err = errors.New("Payment Required")
    }
    if statusCode == 413 {
        err = errors.New("Request Entity Too Large")
    }
    if statusCode == 400 {
        err = errors.New("Bad Request (Malformed Payload)")
    }
    if statusCode == 429 {
        err = errors.New("Too Many Requests")
    }
    return
}
