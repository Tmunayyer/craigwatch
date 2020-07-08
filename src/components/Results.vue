<template>
<div> 
  <ul>Results:
  <div id="resultsBox">
    <div v-if="newListings.length"> {{ newListings.length }} new listings found</div>
    <button v-on:click="()=>{listings}" v-if="newListings.length">Load</button>
    <div v-for="(listing, index) in listings" v-bind:key="listing.url">
    <Listing v-bind:listing="listing" v-bind:index="index" ></Listing>
    </div>
  </div>
  </ul>
  </div>
</template>

<script>
 import Listing from './Listing.vue';

  export default {
    name: 'Results',
    data() {
      return {
      listings: [
        {name: 'post1', content: 'posting content1', date: '12/19/19', id: '1000293'}
        ],
        newListings: [],
        searchId: 1
        }
    },
    methods: {
      update: () => {
          async function getSearchList() {
            const response = await fetch(`/api/v1/listing?ID=${this.searchId}&Datetime=${Date.now()}`)
            return await response.json()
            //post new listings or j the price somewhere too?
        }
        var results = getSearchList()
        if(results.HasNewListings){
           this.listings = this.newListings.concat(results.Listings)
        }
        
      }
    },
    components: {
      Listing
    }

  }
</script>
