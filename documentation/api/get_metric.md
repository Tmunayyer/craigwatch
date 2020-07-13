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
        ID: 48,
        InSeconds: 168,
        InMinutes: 2.8, 
        InHours: 0.05
    }
    ```

* **Sample Call:**

```javascript
    async function getSearchList() {
        const response = await fetch("/api/v1/metric?ID=22")
        const metrics = await response.json()
        console.log("the metrics:", metric)
    }
```

* **Notes:**