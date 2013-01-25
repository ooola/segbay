package main

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
	"time"
	"strings"
	"errors"
	"math/big"
	"html/template"
	"crypto/rand"
	"segbay/bid"
	"segbay/ebay"
	"github.com/hoisie/web"
)

func getRand(max int64) string {
    r, err := rand.Int(rand.Reader, big.NewInt(max))
    if err == nil {
        return r.String()
    }
    return ""
}

func getTimeInEbayFormat(plusSeconds uint) string {
	now := time.Now()
	// TODO: account for timezone
	_, offset := now.Zone()
	m := now.Add(time.Duration(plusSeconds)*time.Second)
	m = m.Add(-(time.Duration(offset) * time.Second))
	t := m.Format(time.RFC3339)
	s := strings.Split(t, "-")
	return s[0] + "-" + s[1] + "-" + s[2] + ".000Z"
}

func getTimeInRFC1123(plusSeconds uint) string {
	now := time.Now()
	m := now.Add(time.Duration(plusSeconds)*time.Second)
	return m.Format(time.RFC1123)
}

func getTimeInRFC3339(plusSeconds uint) string {
	now := time.Now()
	m := now.Add(time.Duration(plusSeconds)*time.Second)
	return m.Format(time.RFC3339)
}

func validateForm(ctx *web.Context) (itemid string, amount string, time string, secondsbefore string, e error) {
	e = nil
	itemid = ctx.Params["itemid"];
	amount = ctx.Params["amount"];
	time = ctx.Params["time"];
	secondsbefore = ctx.Params["secondsbefore"];

	if itemid == "" {
		e = errors.New("missing itemid ")
	} else if amount == "" {
		e = errors.New("missing amount parameter")
	} else if time == "" {
		e = errors.New("missing time parameter")
	}
	if secondsbefore == "" {
		secondsbefore = "5"
	}
	return itemid, amount, time, secondsbefore, e
}

func newBidFromURLEncodedForm(ctx *web.Context, redirect string) {
	ctx.Server.Logger.Println("entering newBidFromURLEncodedForm()")
	ctx.Server.Logger.Println("ctx.Params: ")
	ctx.Server.Logger.Println(ctx.Params)
	itemid, amount, formtime, secondsbefore, err := validateForm(ctx); if err != nil {
		ctx.Abort(500, err.Error())
		return
	}
	b, err := bid.NewBid(itemid, amount, formtime, time.RFC3339, secondsbefore, ctx.Server.Logger); if err != nil {
		ctx.Abort(500, fmt.Sprintf("%+v", err))
		return
	}
	ctx.Server.Logger.Println(fmt.Sprintf("New Bid: %+v", b))
	if redirect != "" {
		ctx.Server.Logger.Println("Redirecting to: " + redirect)
		ctx.Redirect(302, redirect)
	}
	defer ctx.Server.Logger.Println("leaving newBidFromURLEncodedForm()")
}

func jsonGetEbaytime (ctx *web.Context) {
	t, err := ebay.GetEbayTime(); if err != nil {
		ctx.Abort(500, err.Error())
		return
	}
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		ctx.Abort(500, err.Error())
		return
	}
	//ctx.WriteString(`{ "time": "` + s + `" }`)
	s, err := parsedTime.MarshalJSON()
	if err != nil {
		ctx.Abort(500, err.Error())
		return
	}
	ctx.WriteString(`{ "time": ` + string(s) + ` }`)
}

func jsonGetLocaltime (ctx *web.Context) {
	var s string
	now := time.Now()
	l, err := time.LoadLocation("UTC")
	if err != nil {
		s = "internal error: unabled to load UTC"
	} else {
		tUTC := now.In(l)
		s = tUTC.Format(time.RFC3339)
	}
	ctx.WriteString(`{ "time": "` + string(s) + `" }`)
}

func jsonListBids (ctx *web.Context) {
	data, err := json.Marshal(bid.PendingBids); if err != nil {
		ctx.Abort(500, err.Error())
		return
	}
	ctx.ContentType("json")
	ctx.WriteString(string(data))
}

func jsonGetBid (ctx *web.Context, id string) {
	data, err := json.Marshal(bid.PendingBids[id]); if err != nil {
		ctx.Abort(500, err.Error())
		return
	}
	ctx.ContentType("json")
	ctx.WriteString(string(data))
}

func jsonGetCompletedBids(ctx *web.Context) {
	data, err := json.Marshal(bid.CompletedBids); if err != nil {
		ctx.Abort(500, err.Error())
		return
	}
	ctx.ContentType("json")
	ctx.WriteString(string(data))
}

// creates a new Bid from a JSON encoded BidForm HTTP POST'ed
func jsonPostBid (ctx *web.Context) {
	ctx.Server.Logger.Println("entering jsonPostBid()")
	ctx.Server.Logger.Println("ctx.Params: %+v", ctx.Params)
	ctx.Server.Logger.Println("ctx.Request: %+v", ctx.Request)
	ctx.Server.Logger.Println("ctx.Header: %+v", ctx.Request.Header.Get("Content-Type"))
	ctype := ctx.Request.Header.Get("Content-Type")
	if strings.EqualFold(ctype, "application/x-www-form-urlencoded") {
		ctx.Server.Logger.Println("Decoded application/x-www-form-urlencoded BidForm: ")
		newBidFromURLEncodedForm(ctx, "")
	} else if strings.EqualFold(ctype, "application/json") {
		var f bid.BidForm
		json.NewDecoder(ctx.Request.Body).Decode(&f)
		ctx.Server.Logger.Println("Decoded application/json BidForm: " + fmt.Sprintf("%+v", f))
		b, err := bid.NewBid(f.Itemid, f.Amount, f.Time, time.RFC1123, f.SecondsBefore, ctx.Server.Logger); if err != nil {
			ctx.Abort(500, err.Error())
			return
		}
		// echo the bid back
		jsonGetBid(ctx, b.Itemid)
	} else {
		ctx.Abort(500, "Unknown Content-Type: " + ctype)
	}
	defer ctx.Server.Logger.Println("leaving jsonPostBid()")
}

const uiTemplateHTML = `
<html>
<head>
<title>BidServer</title>
<script src="http://code.jquery.com/jquery-latest.min.js" type="text/javascript"></script>
<script type="text/javascript">

$(document).ready(function() {

$.getJSON('/bids/', function(data) {
  var items = [];
  console.log("printing data");
  console.log(data);

  $.each(data, function(index, bid) {
    //alert(index + ': ' + bid.Amount); 
    //$("#pending").append(document.createTextNode("<p>" + bid.Amount + "</p>"));
    //$("#pending").appendChild(document.createTextNode("<p>" + bid.Amount + "</p>"));
    //document.getElementById("pending").appendChild(document.createTextNode("<p>" + bid.Amount + "</p>"));
    var newbid = '<p>Itemid: ' + bid.Itemid +
                 ' Amount: ' + bid.Amount +
                 ' SleepTime: ' + bid.SleepTimeS +
                 ' BidTime: ' + bid.BidTime +
                 ' AuctionEnd: ' + bid.AuctionEnd +
		 '</p>';
    $("#pending").append(newbid);
    //$("#pending").append('<p>' + bid.Amount + '</p>');
    //$("<img/>").attr("alt", item.itemid).appendTo("#pending");
  });
});

$.getJSON('/completedbids/', function(data) {
  var items = [];
  console.log("printing data");
  console.log(data);

  $.each(data, function(index, bid) {
    var completedbid = '<p>Succeded: ' + bid.Succeded +
                 ' Itemid: ' + index +
                 ' Timestamp: ' + bid.Timestamp +
                 ' High Bidder: ' + bid.HighBidder +
                 ' Current Price: ' + bid.CurrentPrice +
                 ' Failure Message: ' + bid.FailureMessage +
		 '</p>';
    $("#completed").append(completedbid);
  });
});

});

</script>                                                               
</head>
<body>
<h2>Add New Bid</h2>
<form action="/ui/" method="POST">

<label for="a">Item ID</label>
<input id="a" type="text" name="itemid" value="{{.Itemid}}"/>
<br>
<label for="b">Amount ($)</label>
<input id="b" type="text" name="amount" value="{{.Amount}}"/>
<br>
<label for="c">Auction end time (Local Time))</label>
<input id="c" type="text" name="time" value="{{.Time}}"/>
<br>
<label for="d">Seconds before auction end time</label>
<input id="d" type="text" name="secondsbefore" value="{{.SecondsBefore}}"/>
<br>
<input type="submit" name="Submit" value="Submit"/>
</form> 
<h2>Pending Bids</h2>
<div id="pending"></div>
<h2>Completed Auctions</h2>
<div id="completed"></div>
</body>
</html>
`

func uiGet(ctx *web.Context) {

	// create a new form with some random values for testing
	b := bid.BidForm{Itemid: getRand(10000),
			 Time: getTimeInRFC3339(60),
			 Amount: getRand(100),
			 SecondsBefore: getRand(10)}

	var uiTemplate = template.Must(template.New("index").Parse(uiTemplateHTML))

	if err := uiTemplate.Execute(ctx.ResponseWriter, b); err != nil {
		ctx.Abort(500, err.Error())
	}
}

func uiPost(ctx *web.Context) {
	newBidFromURLEncodedForm(ctx, "/ui/")
}

func main() {
	//bid.SimulateBids()
	logger := log.New(os.Stdout, "bidserver: ", log.Ldate|log.Ltime)

	/* json interface */
	web.Get("/ebaytime/", jsonGetEbaytime)
	web.Get("/localtime/", jsonGetLocaltime)
	web.Get("/completedbids/", jsonGetCompletedBids)
	web.Get("/bids/", jsonListBids)
	web.Get("/bids/(.*)", jsonGetBid)
	web.Post("/bids/", jsonPostBid)

	/* main (simple) form ui */
	web.Get("/ui/", uiGet)
	web.Post("/ui/", uiPost)

	web.Get("/json", func(ctx *web.Context) string {
		ctx.ContentType("json")
		data, _ := json.Marshal(ctx.Request)
		return string(data)
	})

	web.SetLogger(logger)
	web.Run("127.0.0.1:9999")
}
