[Back to API Table of Contents](./table_of_contents.MD) <br/>
[Back to README](../../README.md)

**Title**
----
  Save a new search to be monitored by the application.

* **URL**

  /api/v1/search

* **Method:**

    `POST`
  
*  **URL Params**

   None

* **Data Params**

   **Required:**
 
   `Name=[string]` <br/>
   `URL=[string]`

* **Success Response:**

    **Code:** 200 <br />
    **Content:**
    ```javascript
    [
        {
            "ID": 22,
            "Name": "Major Tom",
            "URL": "https://newyork.craigslist.org/search/sss?query=ground%20control\u0026sort=rel",
            "Confirmed": false,
            "Interval": 0,
            "CreatedOn": "2020-07-06T18:51:51.516996-04:00",
            "UnixCutoffDate": 1594056360
        }
    ]
    ```
 
* **Error Response:**

    Most commonly due to the body missing a field or properties misnamed.

   **Code:** 400 BAD REQUEST <br />
   **Content:** `could not decode user payload`

* **Sample Call:**

```javascript
    async function getSearchList() {
        const data = JSON.stringify({
            URL: "https://newyork.craigslist.org/search/sss?query=ground%20control\u0026sort=rel",
            Name: 
        })

        const response = await fetch("/api/v1/search", {
            method: "POST",
            body: data
        })
        const savedSearch = await response.json()
        console.log("the savedSearch:", savedSearch)
    }
```

* **Notes:**

  This will spawn a polling service automatically. It may take some time for listings to load for the new search URL.