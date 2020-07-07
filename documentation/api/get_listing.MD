[Back to API Table of Contents](./table_of_contents.MD) <br/>
[Back to README](../../README.md)

**Title**
----
  Returns a JSON array of searches.

* **URL**

  /api/v1/listing

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `ID=[integer]` <br/>

   `Datetime=[integer]` - Datetime is a unix time stamp. If there is no knowledge of a recent listing, 0 should be sent. If 0 is sent, all stored listings will be returned. If a timestamp is sent, all listings after (exclusive) will be returned.

   **Optional:**
 
   None

* **Success Response:**

    A good request with some new listings after the provided timestamp. <br/>
   **Code:** 200 <br />
    **Content:** 
    ```javascript
    {
        "HasNewListings": true,
        "Listings": [
            {
                "ID": 17,
                "SearchID": 22,
                "DataPID": "7149585315",
                "DataRepostOf": "7047865536",
                "UnixDate": 1594056360,
                "Title": "TASCAM DM3200 Digital Mixer w/ Parametric EQ Effects Pro Tools Control",
                "Link": "https://newyork.craigslist.org/wch/msg/d/carmel-tascam-dm3200-digital-mixer/7149585315.html",
                "Price": 750,
                "Hood": " (Westchester Putnam County)"
            },
            {
                "ID": 28,
                "SearchID": 22,
                "DataPID": "7154598663",
                "DataRepostOf": "",
                "UnixDate": 1594055160,
                "Title": "2017 Yamaha SCR 950 *LEFTOVER SALE!*",
                "Link": "https://newyork.craigslist.org/lgi/mcd/d/plainfield-2017-yamaha-scr-950-leftover/7154598663.html",
                "Price": 6495,
                "Hood": " (Motorsports Nation Plainfield)"
            }
        ];
    }
    ```

    A good request with no new listings after the provided timestamp. <br/>
    **Code:** 200 <br />
    **Content:** 
    ```javascript
        {
            "HasNewListings":false,
            "Listings":null
        } 
    ```

* **Sample Call:**

```javascript
    async function getSearchList() {
        const response = await fetch("/api/v1/listing?ID=22&Datetime=1594051440")
        const list = await response.json()
        console.log("the list:", list)
    }
```

* **Notes:**

  The point of the time stamp is really to allow the same endpoint to serve two purposes. First, to populate a page on an initial load. Second, to minimize the amount of data being transfered by only provided new listings to the requester that have unix timestamps after what has been provided.