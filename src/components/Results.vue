<template>
  <div>
    <ul>
      <div>Search Results For:</div>
      <div>Name:</div>
      <div>URl:</div>

      <div id="resultsBox">
        <div v-for="(listing, index) in resultList" v-bind:key="listing.URL">
          <Listing v-bind:listing="listing" v-bind:index="index"></Listing>
        </div>
      </div>
    </ul>
  </div>
</template>

<script>
import Listing from "./Listing.vue";

export default {
  name: "Results",
  components: {
    Listing
  },
  data() {
    return {
      resultList: [],
      newListings: [],
      searchID: this.$route.params.ID,
      // UnixTimestamp is the cutoff to use when requesting only new listings, should be 0 on load
      unixTimestamp: 0
    };
  },
  beforeMount: async function() {
    const resultList = await this.getResultList();

    this.resultList = resultList.Listings;
  },
  methods: {
    getResultList: async function getResultList() {
      const response = await fetch(
        `/api/v1/listing?ID=${this.searchID}&Datetime=${this.unixTimestamp}`
      );
      const list = await response.json();

      return list;
    },
    update: () => {
      async function getSearchList() {
        const response = await fetch(
          `/api/v1/listing?ID=${this.searchId}&Datetime=${Date.now()}`
        );
        return await response.json();
        //post new listings or j the price somewhere too?
      }
      var results = getSearchList();
      if (results.HasNewListings) {
        this.listings = this.newListings.concat(results.Listings);
      }
    }
  }
};
</script>
