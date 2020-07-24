[Back to API Table of Contents](./table_of_contents.MD) <br/>
[Back to README](../../README.md)

**Title**
----
  Returns a JSON object containing metrics for a specific search.

* **URL**

  /api/v1/metric

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `ID=[integer]` <br/>

   **Optional:**
 
   None

* **Success Response:**

   **Code:** 200 <br />
    **Content:** 
    ```javascript
    {
        "Label": "2020-07-22 16:06:09",
        "TopUnixDate": 1595448369,
        "BotUnixDate": 1595444769,
        "Count": 0,
    }
    ```

* **Sample Call:**

```javascript
    async function getSearchList() {
        const response = await fetch("/api/v1/activityChart?ID=22")
        const metrics = await response.json()
        console.log("the activity:", activity)
    }
```

* **Notes:**