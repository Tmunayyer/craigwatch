<template>
  <ul id="searchList">
    <li v-for="search in searchList" :key="search.ID" v-on:click="redirector(search)">
      <div>Name: {{ search.Name }}</div>
      <div>Created On: {{search.CreatedOn}}</div>
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