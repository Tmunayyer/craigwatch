<style module>
.search-form {
  box-sizing: border-box;
  max-width: 250px;

  display: flex;
  flex-direction: column;

  /* margin: 5px 5px 5px 5px; */
}

.search-field {
  margin-bottom: 0.5em;
}

.search-label {
  width: 100%;
  font-weight: 600;
}

.search-input {
  box-sizing: border-box;
  width: 100%;
  /* padding: 0.3em; */
}

.submit-button {
  width: fit-content;
  align-self: flex-end;
}
</style>

<template>
  <div class="search-form">
    <div class="search-field">
      <div class="search-label">search name:</div>
      <input class="search-input" v-model="name" type="text" />
    </div>

    <div class="search-field">
      <div class="search-label">craigslist search url:</div>
      <input class="search-input" v-model="url" type="text" />
    </div>

    <button class="submit-button" v-on:click="handleSubmit">monitor listings</button>
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
