<template>
<div>
  <div class="inputs">
  <h3>craigslist search url:</h3>
  <input type="text"/>
  <br>
  <h3>search name:</h3>
  <input type="text"/>
  <br>
  <button onClick="">find listings</button>
  </div>
</div>
</template>

<script>
  
  export default {
    name: 'CreateSearch',
    data() {
      return {
        listings: null,
        searchTitlesOnly: false,
        hasImage: false,
        postedToday: false,
        bundleDuplicates: false,
        includeNearbyAreas: false
      };
  },
  methods: {
    inititalSearch : () => {
      async function getSearchList() {
            const response = await fetch(`/api/v1/listing?ID=${this.searchId}&Datetime=0`)
            return await response.json()
            //post new listings or j the price somewhere too?
        }
        getSearchList()
    }
  }
  }
</script>
