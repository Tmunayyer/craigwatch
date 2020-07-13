<style module>
.page-container {
  box-sizing: border-box;
  width: 100%;
}
.result-header {
  padding: 0.2em 0.2em 0.2em 0.2em;
  margin-bottom: 1em;
  overflow: hidden;
}

.result-header-lead {
  padding: 0.2em 0.2em 0.2em 0.2em;
}

.result-header-name {
  padding: 0.2em 0.2em 0.2em 0.2em;
  font-size: 1.2em;
  font-weight: bold;
}

.result-header-url {
  padding: 0.2em 0.2em 0.2em 0.2em;
  text-decoration: none;
  font-size: 0.8em;
}

.result-list {
  padding: 0.2em 0.2em 0.2em 0.2em;

  box-sizing: border-box;
  width: 100%;
  max-width: 625px;
}

.result-listitem {
  box-shadow: 0 0px 2px 0 rgba(0, 0, 0, 0.3);
  transition: 0.3s;
  padding: 0.5em;
  border-radius: 4px;
  margin-bottom: 0.5em;
}

/* .result-listitem:hover {
  box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.6);
  cursor: pointer;
} */

.listitem-title {
  font-size: 1em;
  font-weight: bold;
}

.listitem-date {
  font-size: 0.8em;
}

.listitem-body {
  padding: 0.2em 0.2em 0.2em 0.2em;
  overflow: hidden;
}

.price {
  font-weight: bold;
}
</style>

<template>
  <div class="page-container">
    <div class="result-header">
      <div class="result-header-lead">Results For:</div>
      <div class="result-header-name">{{ this.searchDetails.Name }}</div>
      <a
        class="result-header-url"
        v-bind:href="this.searchDetails.URL"
        target="_blank"
      >{{ this.searchDetails.URL }}</a>
    </div>
    <ul class="result-list">
      <li v-for="(listing) in resultList" v-bind:key="listing.URL">
        <div class="result-listitem">
          <div class="listitem-header">
            <div class="listitem-title">{{ listing.Title }}</div>
            <div class="listitem-date">{{ formatDate(listing.UnixDate*1000) }}</div>
          </div>
          <hr />
          <div class="listitem-body">
            <div class="price">${{ listing.Price }}</div>
            <br />
            <a
              class="result-header-url"
              v-bind:href="listing.Link"
              target="_blank"
            >{{ listing.Link }}</a>
          </div>
        </div>
      </li>
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
      searchDetails: {},
      // UnixTimestamp is the cutoff to use when requesting only new listings, should be 0 on load
      unixDate: 0
    };
  },
  beforeMount: async function() {
    let initDetails = await this.getSearchDetails();
    this.searchDetails = initDetails;

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
    getResultList: async function() {
      const response = await fetch(
        `/api/v1/listing?ID=${this.searchID}&Datetime=${this.unixDate}`
      );
      const list = await response.json();

      return list;
    },
    getSearchDetails: async function() {
      const response = await fetch(`/api/v1/search?ID=${this.searchID}`);
      const details = await response.json();

      return details;
    },
    formatDate: function(unixTimestamp) {
      var options = {
        year: "numeric",
        month: "numeric",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit"
      };
      var today = new Date(unixTimestamp);

      return today.toLocaleDateString("en-US", options);
    }
  }
};
</script>
