<style scoped>
#search-list {
  box-sizing: border-box;
  width: 100%;
  max-width: 625px;
}

.search-listitem {
  box-shadow: 0 0px 2px 0 rgba(0, 0, 0, 0.3);
  transition: 0.3s;
  padding: 0.5em;
  border-radius: 4px;
  margin-bottom: 0.5em;
}

.search-listitem:hover {
  box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.6);
  cursor: pointer;
}

.listitem-header {
  box-sizing: border-box;
  padding: 0.2em 0.2em 0.2em 0.2em;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-content: center;
  overflow: hidden;
  height: fit-content;
}

.header-name {
  padding: 0.2em 0.2em 0.2em 0.2em;
  min-width: min-content;
  font-weight: bold;
  font-size: 1.2em;
  align-self: center;
}

.header-date {
  padding: 0.2em 0.2em 0.2em 0.2em;
  font-size: 0.8em;
  display: flex;
  flex-direction: column;
  align-items: center;

  overflow: hidden;
  white-space: nowrap;

  min-width: fit-content;
}

.listitem-body {
  font-size: 0.9em;
  padding: 0.2em 0.2em 0.2em 0.2em;
  overflow: hidden;
  white-space: nowrap;

  border-top: 1px solid #a7a7a7;
}

legend {
  font-size: 0.8em;
}

.listitem-body-contents {
  display: flex;
  flex-direction: column;
}

.body-data {
  margin-bottom: 0.4em;
}
</style>

<template>
  <ul id="search-list">
    <Error v-if="error" />
    <li v-for="search in searchList" :key="search.ID" v-on:click="redirector(search)">
      <div class="search-listitem">
        <div class="listitem-header">
          <div class="header-name">{{ search.Name }}</div>
          <div class="header-date">
            <div>Monitored Since</div>
            <div>{{ formatDate(search.CreatedOn) }}</div>
          </div>
        </div>
        <fieldset class="listitem-body">
          <legend>details</legend>
          <div class="listitem-body-contents">
            <div class="body-data">URL: {{search.URL}}</div>
            <div class="body-data">Listings Analyzed: {{search.TotalListings}}</div>
          </div>
        </fieldset>
      </div>
    </li>
  </ul>
</template>

<script>
import Error from "./Error.vue";

export default {
  name: "SearchList",
  components: {
    Error
  },
  data() {
    return {
      searchList: [],
      error: false
    };
  },
  beforeMount: async function() {
    try {
      const searchList = await this.getSearchList();
      this.searchList = searchList;
    } catch (err) {
      this.error = true;
    }
  },
  methods: {
    redirector: function(search) {
      this.$router.push(`/result/${search.ID}`);
    },
    getSearchList: async function() {
      const apiUrl = `/api/v1/search`;
      const body = await this.$http(apiUrl);

      return body;
    },
    formatDate: function(date) {
      var options = {
        weekday: "long",
        year: "numeric",
        month: "long",
        day: "numeric"
      };
      var today = new Date(date);

      return today.toLocaleDateString("en-US");
    }
  }
};
</script>