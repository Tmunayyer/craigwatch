[Back to API Table of Contents](./table_of_contents.MD) <br/>
[Back to README](../../README.md)

**Title**
----
  Returns a JSON array of searches.

* **URL**

  /api/v1/search

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   None

   **Optional:**
 
   None

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** 
    ```javascript
    [
        {
            "ID": 22,
            "Name": "Major Tom",
            "URL": "https://newyork.craigslist.org/search/sss?query=ground%20control\u0026sort=rel",
            "CreatedOn": "2020-07-06T18:51:51.516996-04:00",
            "UnixCutoffDate": 1594056360
        }
    ]
    ```


* **Sample Call:**

```javascript
    async function getSearchList() {
        const response = await fetch("/api/v1/search")
        const list = await response.json()
        console.log("the list:", list)
    }
```