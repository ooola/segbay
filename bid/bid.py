#!/usr/bin/env python
import httplib
import sys
from xml.dom.minidom import parse, parseString

#serverUrl = "https://api.sandbox.ebay.com/ws/api.dll"
serverUrl = "api.sandbox.ebay.com:443"

requestXML1 = """<?xml version="1.0" encoding="utf-8"?>
<PlaceOfferRequest xmlns="urn:ebay:apis:eBLBaseComponents">
<ErrorLanguage>en_US</ErrorLanguage>
<EndUserIP>192.168.255.255</EndUserIP>
<ItemID>%s</ItemID>
<Offer>
<Action>Bid</Action>
<MaxBid>%s</MaxBid>
<Quantity>1</Quantity>
</Offer>
<RequesterCredentials>
<eBayAuthToken>AgAAAA**AQAAAA**aAAAAA**Q5P/Tw**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4GhCZGCpQWdj6x9nY+seQ**ytsBAA**AAMAAA**8zSBbAyBODh312UY1gAABEQC1Xk94fV2P16wOnwsEK+rbjVBfGthp732lgdxaz0IxFXLqmiGXzLWP7i5YsRjij6JNAXOqLJg8/9mtauZnYWCoUK5j2DewDIA67/yzVAbH2/9siQOgrOyfZ0i9xZYxpfvQcA4ymsGzKwfjd5abRjWvFJDZNSRaJKm3et+w3M/J54cqkTYtqKvqpq0NhbYbT81H0of1AsKwEIXQb7fwIV9P256X3u+GsJwo2vyzuBejAyJPAv5RHXryPAiEjrxWW+rSEQNOzcMH/O9P05EFBLfirO4WEZ8rSvvA/RnZQlCurrvSECu1axMqpOYCPDDo9Y+KhUcb0/hZvwyacX+OvKYkffz++kAS4QBjYIK+FzCzBXYb1WRPeTnh7EcCwECPwTDQOerxUdO4JUFpxGmIoM/IThbxE2wIzGMYd9tlcZmrM+hClZO2jGDxU9qvQxTAw1QXcLv4RtoC2mzXtKHm6qZrO4ZXmJbfPlttA6gNrducURUZ5McsUnGYRIqsInNvUVu1vU6U+DOf/W/vEm4DcP6kurOSlkbbRNHbTMdZ7HhsPIZqjRdW0A2gNtbyiCNauo9VLUpurZFoKfV5PVpx/OrrzIZeAZ0fzdp8Tg6g7LGS8At5hzoFJ62GNgfxCDHjetGij37rb/LjCpfe10eS1iALZ34K9n1kZ8H9+t5VMo6LBGPij9UT6SSDjWH9V4EcF1WjxLCryl0PjYIOlx6oeTfRRbZ2C+uv5ho3Q3ow4EP</eBayAuthToken>
</RequesterCredentials>
<WarningLevel>High</WarningLevel>
</PlaceOfferRequest>"""

requestXML2 = """<?xml version="1.0" encoding="utf-8"?>
<PlaceOfferRequest xmlns="urn:ebay:apis:eBLBaseComponents">
<ErrorLanguage>en_US</ErrorLanguage>
<EndUserIP>192.168.255.255</EndUserIP>
<ItemID>%s</ItemID>
<Offer>
<Action>Bid</Action>
<MaxBid>%s</MaxBid>
<Quantity>1</Quantity>
</Offer>
<RequesterCredentials>
<eBayAuthToken>AgAAAA**AQAAAA**aAAAAA**1FliUA**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4GhCZSAoQmdj6x9nY+seQ**ytsBAA**AAMAAA**B3PuJCPhfGP/AhMndKIH347qZ6yBVPl2YNldRcoGYe43ETIWQrZ2r6c6k7hNpP4zwe1tvYSwgTqy7ZNZwtkDQFznEN2FeCPNO/CddHKWFMm4+bhhT1xHA6Xp3ywTRJ8D1uBHgj1xxNGvcGmnYHLslUa4a21nO6h9p1Zcu2PoJzcI6kqG8w9S3uK8QxGsR+8qhVmcnGu4bJkCfatSrvQll2AUq+IPDSN/WyvxIfL1hGl2NR+GOjRxKiltqzAAtZg0N78aw6DMRGwkvzhKaGWrd8h9L2keV3HDgIsCiuFlHGC2s7yLzFoYHFaKk/gWjP1rAfXL/4bfk3+1CUQGR8v3umMupMZgJf2bEv/0vxDlI4R98WVtGni/hTGjiy0/EhO0MfeZfHY7SUzz4Dje8u7e9FOd3unXqDQzJSK0DmBSqpVg6/3G97Q9GJKhvcc7CMYz4Q5W0iJHo6StCvwOaWCv/GXxhJN2AM7rzuV+IZXwYpy7yAS1hPGBv2pUscv1L6n6Kk4VKwFuw9LKHbzqQM+Fl0SF3pwRStYTpcl+wM52oltx3vvk7l0WqyMszomHNzm2NN0lP8B9l3msSFJF3p/tZvN4WgHcWGmK5lUXOnHk97WVlrpM69RmjRfUb5i+xzqUHYxaWbNDx5p9FeB1PqDDyKvuElVz5fAdpSFzcALBLSwQqG99uiOOMzg7RQYb71Nrs782UsZnV7iUQljkShIwSYhuyLx5wy9S9W0fR1Vwe5vObBRw5J1mGqCNUZBc20vy</eBayAuthToken>
</RequesterCredentials>
<WarningLevel>High</WarningLevel>
</PlaceOfferRequest>"""

httpHeaders = { "X-EBAY-API-COMPATIBILITY-LEVEL" : "781",
                "X-EBAY-API-DEV-NAME" : "8f9566df-77a3-4489-832a-56952547d200",
                "X-EBAY-API-APP-NAME" : "segbaye3c-c714-4d11-9312-1eb4b80e36c",
                "X-EBAY-API-CERT-NAME" : "9089079c-6f28-4524-8ecf-3ae458ccca96",
                "X-EBAY-API-SITEID" : "0",
                "X-EBAY-API-CALL-NAME" : "PlaceOffer" }

def callEbay(itemid, amount):
    r = requestXML1 % (itemid, amount)

    # specify the connection to the eBay Sandbox environment
    connection = httplib.HTTPSConnection(serverUrl)

    # specify a POST with the results of generateHeaders and generateRequest
    connection.request("POST" , "/ws/api.dll", r, httpHeaders)
    #connection.request("POST" , serverUrl, requestXML, httpHeaders)
    response = connection.getresponse()

    # if response was unsuccessful, output message
    if response.status != 200:
            print "Error sending request:" + response.reason
    else: #response successful
        # store the response data and close the connection
        data = response.read()
        connection.close()

        # parse the response data into a DOM
        response = parseString(data)
        print response.toxml()
        #print response.toprettyxml(indent="  ")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print "usage: %s <itemid> <amount>" % (sys.argv[0])
        exit(1)
    callEbay(sys.argv[1], sys.argv[2])
