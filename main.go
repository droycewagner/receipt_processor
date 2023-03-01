/*
From 
https://github.com/droycewagner/receipt_processor

This furnishes a RESTful API which processes receipts POSTed as JSON blocks, 
assigns (by hashing) an ID to each receipt, and stores a point score for each 
ID, which can be accessed via a GET request. 
*/

package main

import (
    "fmt"
    "log"
    "net/http"
    "math"
    "strings"
    "strconv"
    "regexp"
    "encoding/json"
    "io/ioutil"
    "crypto/sha256"
    "encoding/hex"
    "os"
)

// stores pairs (string, int) corresponding to receipt id's and their point score.
var points = make(map[string]int)

// given a json receipt as a byte array, return a SHA256 hash. 
func makeID (byteArray []byte) string {
    hasher := sha256.New()
    hasher.Write(byteArray)
    sha := hex.EncodeToString(hasher.Sum(nil))
    return sha
}

// struct to parse json receipts
type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price string `json:"price"`
}

//struct to parse json receipts
type Receipt struct {
    Retailer  string `json:"retailer"`
    PurchaseDate  string `json:"purchaseDate"`
    PurchaseTime   string `json:"purchaseTime"`
    Total  string `json:"total"`
    Items  []Item `json:"items"`
}

// Computes the point value of a receipt, given a Receipt struct
func computePoints(receipt Receipt) int {
    points := 0.
    
    // One point for every alphanumeric character in the retailer name.
    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    retailer_alphanum := reg.ReplaceAllString(receipt.Retailer, "")
    points += float64(len(retailer_alphanum))
    
    total, _ := strconv.ParseFloat(receipt.Total,64)
    
    // 50 points if the receipt total is a round dollar amount with no cents.
    if total == math.Trunc(total) {
        points += 50
    }
    
    // 25 points if the receipt total is a multiple of 0.25.
    if 4*total == math.Trunc(4*total) {
        points += 25
    }
    
    // 5 points for every two items on the receipt.
    points += float64(5*(len(receipt.Items[:])/2))
    
    // For each item, if trimmed length . . .
    for _, item := range receipt.Items {
        if len(strings.TrimSpace(item.ShortDescription)) % 3 == 0 {
            price, _ := strconv.ParseFloat(item.Price,64)
            points += math.Ceil(price*0.2)
        }
    }
    
    // 6 points if the day in the purchase date is odd.
    var dateParts = strings.Split(receipt.PurchaseDate,"-")
    dayDigit, _ := strconv.Atoi(dateParts[len(dateParts)-1])
    if dayDigit % 2 == 1 {
        points += 6
    }
    
    // 10 points if the time of purchase is between 2pm and 4pm.
    var timeParts = strings.Split(receipt.PurchaseTime,":")
    hour, _ := strconv.Atoi(timeParts[0])
    if hour >= 14 && hour < 16 {
        points += 10
    }
    
    return int(points)
}

// returns the points score from the JSON block in fileName; used for unit tests. 
func pointsFromFile (fileName string) int {
    // get contents of receipt from file
    exFile, err := os.Open(fileName)
    if err != nil {
      fmt.Println(err)
    }
    defer exFile.Close()
    
    // unpack file contents to Receipt object
    byteValue, _ := ioutil.ReadAll(exFile)
    var receipt Receipt
    json.Unmarshal(byteValue, &receipt)
    
    // run test
    return computePoints(receipt)
}

// handles GET and POST requests
func ReceiptREST(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        // first make sure that the URL is nominally correct
        splitURL:=strings.Split(r.URL.Path,"/")
        if splitURL[1] != "receipts" {
            http.Error(w, "404 NOT FOUND", http.StatusNotFound)
        }
        if splitURL[3] != "points" {
            http.Error(w, "404 NOT FOUND", http.StatusNotFound)
        }
        // check that the specified ID is valid
        _, ok := points[splitURL[2]]
        if ! ok {
            http.Error(w, "404 NOT FOUND", http.StatusNotFound)
        } else {
            // return the point score
            fmt.Fprintf(w, "{ \"points\" : %d }", points[splitURL[2]])
        }
    case "POST":
        if r.URL.Path != "/receipts/process" {
            http.Error(w, "404 NOT FOUND", http.StatusNotFound)
        } else {
            // unpack JSON object into a Receipt struct
            var receipt Receipt
            reqBody, _ := ioutil.ReadAll(r.Body)
            err := json.Unmarshal(reqBody,&receipt)
            // throw error if json does not unpack properly
            if err !=nil {
                http.Error(w, "404 NOT FOUND\n invalid JSON block", http.StatusNotFound)
            }
            // create a new ID & store the receipt points
            newid := makeID(reqBody)
            points[newid] = computePoints(receipt)
            // return the id
            fmt.Fprintf(w, "{ \"id\" : \"%s\" }", newid)
        }
    default:
        fmt.Fprintf(w, "Only GET and POST methods are supported.")
    }
}

func main() {
    http.HandleFunc("/", ReceiptREST)
    
    fmt.Printf("Listening on port 8080...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
