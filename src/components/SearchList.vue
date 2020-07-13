<style module>
.search-listitem {
  box-shadow: 0 0px 2px 0 rgba(0, 0, 0, 0.3);
  transition: 0.3s;
  padding: 0.5em;
  border-radius: 4px;
}

/* On mouse-over, add a deeper shadow */
.search-listitem:hover {
  box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.6);
  cursor: pointer;
}
</style>

<template>
  <ul id="searchList">
    <li v-for="search in searchList" :key="search.ID" v-on:click="redirector(search)">
      <div class="search-listitem">
        <div>Name: {{ search.Name }}</div>
        <div>Created On: {{search.CreatedOn}}</div>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  name: "SearchList",
  data() {
    return {
      searchList: []
    };
  },
  beforeMount: async function() {
    const searchList = await this.getSearchList();
    this.searchList = searchList;
  },
  methods: {
    redirector: function(search) {
      this.$router.push(`/result/${search.ID}`);
    },
    getSearchList: async function() {
      const apiUrl = `/api/v1/search`;
      const response = await fetch(apiUrl);
      const body = await response.json();

      return body;
    }
  }
};
</script>