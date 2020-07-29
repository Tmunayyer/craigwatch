<style scoped>
.result-list {
  padding: 0.2em 0.2em 0.2em 0.2em;
  margin-bottom: 2em;

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

.listitem-title {
  font-size: 1em;
  font-weight: bold;
  padding: 0.2em 0.2em 0.2em 0.2em;
}

.listitem-date {
  font-size: 0.8em;
  padding: 0.2em 0.2em 0.2em 0.2em;
}

.listitem-body {
  padding: 0.2em 0.2em 0.2em 0.2em;
  overflow: hidden;
}

.price {
  font-weight: bold;
}

.button {
  float: right;
}
</style>

<template>
  <div class="result-list">
    <Error v-if="error" />
    <ul>
      <li v-for="(listing) in resultListPreview" v-bind:key="listing.URL">
        <div class="result-listitem">
          <div>
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
    <button v-on:click="showMore">show more</button>
  </div>
</template>

<script>
import Error from "./Error.vue";
import { spinnerState } from "./Spinner.vue";

import util from "../utility.js";

export default {
  name: "ResultList",
  components: {
    Error,
  },
  props: ["searchID"],
  data() {
    return {
      resultListPreview: [],
      resultList: [],
      error: false,
      // UnixTimestamp is the cutoff to use when requesting only new listings, should be 0 on load
      unixDate: 0,
      polling: null,

      // variable to track page size
      pageSize: 15,
    };
  },
  beforeMount: async function () {
    try {
      const initResultList = await this.getResultList(true);
      this.resultList = initResultList;
      this.resultListPreview = initResultList.slice(0, this.pageSize);

      if (this.resultList.length) {
        this.unixDate = this.resultList[0].UnixDate;
      }
    } catch (err) {
      this.error = true;
    }
  },
  mounted: function () {
    this.startPolling();
  },
  beforeDestroy: function () {
    this.stopPolling();
  },
  methods: {
    formatDate: util.formatDate,
    getResultList: async function (retry) {
      const { fetch, fetchRetry } = this.$http;

      let url = `/api/v1/listing?ID=${this.searchID}&Datetime=${this.unixDate}`;

      let result;
      if (retry) {
        result = await fetchRetry(url, {}, (data) => !data.HasNewListings);
      } else {
        result = await fetch(url);
      }

      if (result.HasNewListings) {
        return result.Listings;
      } else {
        return [];
      }
    },
    startPolling: function () {
      // update list every 60 seconds
      this.polling = setInterval(async () => {
        // let these go without setting error
        const updatedResultList = await this.getResultList();

        if (updatedResultList.length) {
          this.unixDate = updatedResultList[0].UnixDate;
          this.resultList = updatedResultList.concat(this.resultList);
          this.resultListPreview = updatedResultList.concat(
            this.resultListPreview
          );
        }
      }, 60000);
    },
    stopPolling: function () {
      // the next two lined must always be grouped together
      clearInterval(this.polling);
      this.polling = null;
    },
    showMore: function () {
      const newPageSize = this.resultListPreview.length + 15;

      this.resultListPreview = this.resultList.slice(0, newPageSize);
    },
  },
};
</script>