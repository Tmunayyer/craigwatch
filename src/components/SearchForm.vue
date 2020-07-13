<style module>
.search-form {
  display: flex;
  flex-direction: row;
}
</style>

<template>
  <div class="search-form">
    <h3>search name:</h3>
    <input v-model="name" type="text" />
    <br />
    <h3>craigslist search url:</h3>
    <input v-model="url" type="text" />
    <br />
    <button v-on:click="handleSubmit">find listings</button>
  </div>
</template>

<script>
export default {
  name: "SearchForm",
  data() {
    return {
      name: "",
      url: ""
    };
  },
  methods: {
    redirector: function(search) {
      this.$router.push(`/result/${search.ID}`);
    },
    handleSubmit: async function() {
      const apiUrl = `/api/v1/search`;
      const apiOptions = {
        method: "POST",
        body: JSON.stringify({ Name: this.name, URL: this.url })
      };

      const response = await fetch(apiUrl, apiOptions);
      const body = await response.json();

      this.redirector(body);
    }
  }
};
</script>
