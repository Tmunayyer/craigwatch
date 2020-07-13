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
import Listing from "./components/Listing.vue";

let _RESULT_FETCH_INTERVAL; // to track the interval for cleanup later

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
      unixDate: 0
    };
  },
  beforeMount: async function() {
    let initResultList = await this.getResultList();

    if (initResultList.HasNewListings) {
      this.unixDate = initResultList.Listings[0].UnixDate;
      this.resultList = initResultList.Listings;
    } else {
      setTimeout(async () => {
        // on new search created, the backend takes a second to get listings.
        // because of this, retry after 3 seconds and then spawn interval.
        initResultList = await this.getResultList();

        if (initResultList.HasNewListings) {
          this.unixDate = initResultList.Listings[0].UnixDate;
          this.resultList = initResultList.Listings;
        }
      }, 3000);
    }

    // update list every 60 seconds
    _RESULT_FETCH_INTERVAL = setInterval(async () => {
      const updatedResultList = await this.getResultList();

      if (updatedResultList.HasNewListings) {
        this.unixDate = updatedResultList.Listings[0].UnixDate;
        this.resultList = updatedResultList.Listings.concat(
          resultList.Listings
        );
      }
    }, 60000);
  },
  beforeDestroy: function() {
    clearInterval(_RESULT_FETCH_INTERVAL);
  },
  methods: {
    getResultList: async function getResultList() {
      const response = await fetch(
        `/api/v1/listing?ID=${this.searchID}&Datetime=${this.unixDate}`
      );
      const list = await response.json();

      return list;
    }
  }
};
</script>
