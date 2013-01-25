package bid

import (
	"fmt"
	"log"
	"time"
	"strings"
	"errors"
	"strconv"
	"encoding/xml"
	"os/exec"
	"math/big"
	"crypto/rand"
)

var PendingBids map[string] *Bid
var CompletedBids map[string] *CompletedBid

var simulateBids bool

type BidForm struct {
	Itemid string
	Amount string
	Time string
	SecondsBefore string
}

type CompletedBid struct {
	Succeded bool
	Timestamp string
	HighBidder string
	CurrentPrice string
	FailureMessage string
}

// main Bid structure used to store information for pending bids
type Bid struct {
	Itemid string
	Amount float64
	SleepTime time.Duration `json:"-"`
	SleepTimeS string
	BidTime time.Time
	AuctionEnd time.Time
	Logger *log.Logger `json:"-"`
}


type PlaceOfferResponseAck struct {
    XMLName xml.Name `xml:"PlaceOfferResponse"`
    Timestamp string   `xml:"Timestamp"`
    Ack string   `xml:"Ack"`
}

type Error struct {
	ShortMessage string
	LongMessage string
	ErrorCode string
	SeverityCode string
	ErrorClassification string
}

type PlaceOfferResponseFailure struct {
    XMLName xml.Name `xml:"PlaceOfferResponse"`
    Timestamp string   `xml:"Timestamp"`
    Ack string   `xml:"Ack"`
    Errors []Error
}

type UserID struct {
    UserID string
}

type Status struct {
	ConvertedCurrentPrice string
	CurrentPrice string
	HighBidder []UserID
	MinimumToBid string
}

type PlaceOfferResponseSuccess struct {
    XMLName xml.Name `xml:"PlaceOfferResponse"`
    Timestamp string   `xml:"Timestamp"`
    Ack string   `xml:"Ack"`
    Version string `xml:"Version"`
    Build string `xml:"Build"`
    UsageData string `xml:"UsageData"`
    SellingStatus []Status
}

type PlaceOfferResponse struct {
    success *PlaceOfferResponseSuccess
    failure *PlaceOfferResponseFailure
}

func init() {
	PendingBids = make(map[string]*Bid)
	CompletedBids = make(map[string]*CompletedBid)
}
// helpful to undestand timeouts http://code.google.com/p/go-wiki/wiki/Timeouts

//http://golang.org/doc/effective_go.html#allocation_new see 'composite literals'
// incoming time is in rfc1123 GMT time
func NewBid(Itemid string, Amount string, Time string, TFormat string, SecondsBefore string, Logger *log.Logger) (*Bid, error) {

	Logger.Println("entering NewBid()")

	var ebaytime time.Time

	b := Bid{Itemid: Itemid, Amount: 0, Logger: Logger}

	if Itemid == "" || len(b.Itemid) > 50 {
		return nil, errors.New("Itemid is invalid")
	}

	//amount, err := strconv.ParseUint(Amount, 10, 0) ; if err != nil {
	amount, err := strconv.ParseFloat(Amount, 64) ; if err != nil {
		return nil, err
	}
	b.Amount = amount

	if b.Amount < 1 {
		return nil, errors.New("Amount can not be less than 1")
	}
	if b.Amount > 1000 {
		return nil, errors.New("Amount exceeds limit (1000)")
	}

	/*
	// so chop off last fractional second and use time.RFC3339
	Logger.Println("Time: " + Time)
	s := strings.Split(Time, ".")
	ebaytime, err := time.Parse(time.RFC3339, s[0] + "Z") // TODO: add ts offset
	if err != nil {
		b.Logger.Println("time.Parse error: " + fmt.Sprintf("%+v", err))
		return nil, err
	}
	b.AuctionEnd = ebaytime
	*/
	if (TFormat == time.RFC3339) {
		Logger.Println("Time (RFC3339): " + Time)
		var t string
		s := strings.Split(Time, ".")
		if len(s) > 1 {
			// bug: go doesn't handle ebay time with fractions e.g. 2012-09-25T04:36:44.000Z
			// so strip the .000 if it is present and add back the 'Z'
			t = s[0] + "Z"
		} else {
			t = s[0]
		}
		Logger.Println("Fixedup Time (RFC3339): " + t)
		ebaytime, err = time.Parse(time.RFC3339, t)
		if err != nil {
			b.Logger.Println("time.Parse error: " + fmt.Sprintf("%+v", err))
			return nil, err
		}
		b.Logger.Println("Auction end in localtime: ", ebaytime.Local())
		b.AuctionEnd = ebaytime
	} else if (TFormat == time.RFC1123) {
		Logger.Println("Time (RFC1123): " + Time)
		ebaytime, err := time.Parse(time.RFC1123, Time)
		if err != nil {
			fmt.Println("time.Parse error: " + fmt.Sprintf("%+v", err))
			return nil, err
		}
		b.Logger.Println("Auction end in localtime: ", ebaytime.Local())
		b.AuctionEnd = ebaytime
	} else {
		return nil, errors.New("Unknown time format")
	}

	b.Logger.Println("ebaytime: " + ebaytime.String())
	seconds, err := strconv.ParseUint(SecondsBefore, 10, 0); if err != nil {
		return nil, err
	}
	t := time.Second * time.Duration(seconds)
	b.Logger.Println("duration t : " + t.String())
	bidtime := ebaytime.Add(-t)
	b.Logger.Println("bidtime: " + bidtime.String())
	b.BidTime = bidtime

	now := time.Now()
	b.Logger.Println("nowtime: " + now.String())

	/* ebaytime is the future time when bid should be placed */
	b.SleepTime = ebaytime.Sub(now)
	b.SleepTimeS = b.SleepTime.String()
	if b.SleepTime <= 0 {
		return nil, errors.New("Sleep time is less than 0 seconds")
	}
	b.Logger.Println(fmt.Sprintf("--- %+v\n", b))

	// add the pending bid to the global list
	key := b.Itemid
	PendingBids[key] = &b
	go b.SleepThenExecuteBid()

	defer Logger.Println("leaving NewBid()")
	return &b, nil
}

// executes the bid
// should be run in a go routine since the function will sleep
// until the specified time before executing the bid
func (b *Bid) SleepThenExecuteBid() {
	time.Sleep(b.SleepTime)
	b.ExecuteBid()
}

func (b *Bid) ExecuteBid() {
	delete(PendingBids, b.Itemid)
	var out []byte
	n := time.Now()
	b.Logger.Println("Executing bid at: " + n.String())
	if simulateBids == true {
		out = []byte(getRandomResponse())
	} else {
		amtString := fmt.Sprint(float64(b.Amount))
		// TODO: make the path a config variable
		o, err := exec.Command("/usr/local/bidserver/bid.py", b.Itemid, amtString).Output()
		if err != nil {
			b.Logger.Println(err)
			return
		} else {
			b.Logger.Println("./bid.py completed")
			out = o
		}
	}
	b.Logger.Printf("Exec output %s\n", out)
	resp, err := ParsePlaceOfferResponse(string(out))
	if err != nil {
		b.Logger.Println(err)
		return
	}
	b.Logger.Println(fmt.Sprintf("%+v\n", resp))
	if resp.failure != nil {
		str := "Timestamp: " + resp.failure.Timestamp + " Failure ShortMessage: " +
		resp.failure.Errors[0].ShortMessage
		b.Logger.Println(str)
		//CompletedBids = append(CompletedBids, str)
		tmp := CompletedBid{Succeded: false, Timestamp: resp.failure.Timestamp, FailureMessage: resp.failure.Errors[0].ShortMessage}
		CompletedBids[b.Itemid] = &tmp
	} else if resp.success != nil {
		str := "Timestamp: " + resp.success.Timestamp + " Success Highbidder: " +
		resp.success.SellingStatus[0].HighBidder[0].UserID +
		" CurrentPrice: " + resp.success.SellingStatus[0].CurrentPrice;
		b.Logger.Println(str)
		tmp := CompletedBid{	Succeded: true,
					Timestamp: resp.success.Timestamp,
					HighBidder: resp.success.SellingStatus[0].HighBidder[0].UserID,
					CurrentPrice: resp.success.SellingStatus[0].CurrentPrice }
		CompletedBids[b.Itemid] = &tmp
	}
}

func ParsePlaceOfferResponse(response string) (*PlaceOfferResponse, error) {
	resp := new (PlaceOfferResponse)
	resp.success = nil
	resp.failure = nil
	resp_failure := new(PlaceOfferResponseFailure)
	resp_success := new(PlaceOfferResponseSuccess)
	ack := new(PlaceOfferResponseAck)

	err := xml.Unmarshal([]byte(response), &ack)
	if err == nil {
		if strings.EqualFold(ack.Ack, "Failure") {
			err = xml.Unmarshal([]byte(response), &resp_failure)
			resp.failure = resp_failure
		} else if strings.EqualFold(ack.Ack, "Success") {
			err = xml.Unmarshal([]byte(response), &resp_success)
			resp.success = resp_success
		} else {
			err = errors.New("Unknown XML: " + response)
		}
	}

	return resp, err
}

// Simulates bidding by returning successful or failed bids randomly
func SimulateBids() {
    simulateBids = true
}

func getRandomBool() bool {
    one := big.NewInt(1)
    r, err := rand.Int(rand.Reader, big.NewInt(2))
    if err == nil {
        if r.Cmp(one) == 0 {
            return true
        }
    }
    return false
}
func getRandomResponse() string {
    var str string
    if getRandomBool() == true {
        str = `
	        <?xml version="1.0" ?>
	        <PlaceOfferResponse xmlns="urn:ebay:apis:eBLBaseComponents">
	        	<Timestamp>2012-09-03T02:37:30.856Z</Timestamp>
	        	<Ack>Failure</Ack>
	        	<Errors>
	        		<ShortMessage>Auction ended.</ShortMessage>
	        		<LongMessage>This auction has ended.</LongMessage>
	        		<ErrorCode>12243</ErrorCode>
	        		<SeverityCode>Error</SeverityCode>
	        		<ErrorClassification>RequestError</ErrorClassification>
	        	</Errors>
	        	<Version>787</Version>
	        	<Build>E787_CORE_BUNDLED_15273494_R1</Build>
	        </PlaceOfferResponse>
	        `
        } else {
            str = `<?xml version="1.0" ?>
            <PlaceOfferResponse xmlns="urn:ebay:apis:eBLBaseComponents">
              <Timestamp>2012-07-18T01:12:41.617Z</Timestamp>
              <Ack>Success</Ack>
              <Version>779</Version>
              <Build>E779_CORE_BUNDLED_14991381_R1</Build>
              <UsageData>MTMzMDExMDQ5LzE1NTM5Ow**</UsageData>
              <SellingStatus>
                <ConvertedCurrentPrice currencyID="USD">1.0</ConvertedCurrentPrice>
                <CurrentPrice currencyID="USD">1.0</CurrentPrice>
                <HighBidder>
                  <UserID>testuser_olanord</UserID>
                </HighBidder>
                <MinimumToBid currencyID="USD">1.25</MinimumToBid>
              </SellingStatus>
            </PlaceOfferResponse>
	        `
    }
    return str
}
