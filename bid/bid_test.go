package bid

import (
	"testing"
	"strings"
	"log"
	"os"
)

func TestNewBid(t *testing.T) {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	_, err := NewBid("", "100", "Not valid time", "5", l)
	if err == nil {
		t.Error("NewBid failed to detect invalid (too short) itemid")
	}

	_, err = NewBid("0123456789012345678901234567890123456789012345678901", "100", "Not valid time", "5", l)
	if err == nil {
		t.Errorf("NewBid failed to detect invalid (too long) itemid")
	}

	_, err = NewBid("0123456789", "1XYZ", "Not valid time", "5", l)
	if err == nil {
		t.Errorf("NewBid failed to detect invalid amount")
	}

	_, err = NewBid("0123456789", "1", "Mon, 17 Sep 2012 23:45:30 EDT", "5", l)
	if err == nil {
		t.Errorf("NewBid failed to detect a valid amount")
	}

	_, err = NewBid("0123456789", "1001", "Mon, 17 Sep 2012 23:45:30 EDT", "5", l)
	if err == nil {
		t.Errorf("NewBid failed to detect a valid amount")
	}

	_, err = NewBid("0123456789", "-1", "Mon, 17 Sep 2012 23:45:30 EDT", "5", l)
	if err == nil {
		t.Errorf("NewBid failed to detect a valid amount")
	}

	_, err = NewBid("0123456789", "1", "ASDF 33-2-2 2012-04-25T07:03:31.768Z", "5", l)
	if err == nil {
		t.Errorf("NewBid failed to detect an inalid time")
	}

	_, err = NewBid("0123456789", "1", "2012-04-25T07:03:31.768Z", "-1", l)
	if err == nil {
		t.Errorf("NewBid failed to detect a valid amount")
	}

	_, err = NewBid("0123456789", "1", "Mon, 17 Sep 2013 23:45:30 EDT", "5", l)
	if err != nil {
		t.Errorf("NewBid failed to detect a bid")
		t.Errorf(err.Error())
	}


	/*
	if b != nil || err == nil {
		t.Errorf("NewBid failed detect invalid time string")
		//t.Errorf(err)
	}
	*/
}

/*
func TestExecuteBid(t *testing.T) {
    b, err := NewBid("0123456789", "1", "2023-08-25T07:03:31.768Z", "5")
	if err != nil {
		t.Errorf("NewBid failed to detect a bid")
		t.Errorf(err.Error())
	}
    b.ExecuteBid()
}
*/

func TestParsePlaceOfferReponse(t *testing.T) {
    str := `
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

    resp, err := ParsePlaceOfferResponse(str)
	if resp.failure == nil || resp.success != nil {
		t.Errorf("ParsePlaceOfferResponse failed to parse a valid failure response")
		t.Errorf(err.Error())
	}
    if strings.EqualFold(resp.failure.Ack, "Failure") == false {
		t.Errorf("ParsePlaceOfferResponse failed to parse a Failure")
		t.Errorf(err.Error())

    }

    str = `
	<ADSFADShhh?xml version="1.0" ?>
	<PlaceOffe///jasdfj;;adfskrResponse xmlns="urn:ebay:apis:eBLBaseComponents">
		<Timestamp>2012-09-03T02:37:30.856Z</Timestamp>
		<Ack>Failure</Ack>
		<Errors>
			<ShortMessage#@@#$@#JADSFJAFSdll>Auction ended.</ShortMessage>
			<LongMessage>This auction has ended.</LongMessage>
			<ErrorCode>12243</ErrorCode>
			<SeverityCode>Error</SeverityCode>
			<ErrorClassification>RequestError</ErrorClassification>
		</Errors>
		<Version>787</Version>
		<Build>E787_CORE_BUNDLED_15273494_R1</Build>
	</PlaceOfferResponse>
	`

    resp, err = ParsePlaceOfferResponse(str)
	if resp.failure != nil || resp.success != nil {
		t.Errorf("ParsePlaceOfferResponse failed to generate an error on an invalid error response")
		t.Errorf(err.Error())
	}

    str = `<?xml version="1.0" ?>
    <PlaceOfferResponse xmlns="urn:ebay:apis:eBLBaseComponents">
      <Timestamp>
        2012-07-18T01:12:41.617Z
      </Timestamp>
      <Ack>Success</Ack>
      <Version>
        779
      </Version>
      <Build>
        E779_CORE_BUNDLED_14991381_R1
      </Build>
      <UsageData>
        MTMzMDExMDQ5LzE1NTM5Ow**
      </UsageData>
      <SellingStatus>
        <ConvertedCurrentPrice currencyID="USD">
          1.0
        </ConvertedCurrentPrice>
        <CurrentPrice currencyID="USD">
          1.0
        </CurrentPrice>
        <HighBidder>
          <UserID>
            testuser_olanord
          </UserID>
        </HighBidder>
        <MinimumToBid currencyID="USD">
          1.25
        </MinimumToBid>
      </SellingStatus>
    </PlaceOfferResponse>
	`

    resp, err = ParsePlaceOfferResponse(str)
	if resp.success == nil {
		t.Errorf("ParsePlaceOfferResponse generated an error on a valid success response")
        if err != nil {
		    t.Errorf(err.Error())
        } else {
		    t.Errorf("err is nil!!!")
        }
	}
    if strings.EqualFold(resp.success.Ack, "Success") != true {
		t.Errorf("ParsePlaceOfferResponse generated an error on a valid success response")
        t.Errorf("Ack: %v\n", resp.success.Ack)
        t.Errorf("Timestamp: %v\n", resp.success.Timestamp)
    }

    str = `<?xml version="1.0" ?>
    <PlaceOfferResponse xmlns="urn:ebay:apis:eBLBaseComponents">
      <Timestamp>
        2012-07-18T01:12:41.617Z
      </Timestamp>
      <Ack>Unknown</Ack>
    </PlaceOfferResponse>
	`

    resp, err = ParsePlaceOfferResponse(str)
	if err == nil {
		t.Errorf("ParsePlaceOfferResponse did not generate an error on a invalid response")
	}
}

func TestGetRandomBool(t *testing.T) {
    t.Log("Getting a five random bools")
    t.Log(getRandomBool())
    t.Log(getRandomBool())
    t.Log(getRandomBool())
    t.Log(getRandomBool())
    t.Log(getRandomBool())
}
