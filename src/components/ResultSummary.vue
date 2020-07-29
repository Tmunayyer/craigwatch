<style scoped>
.result-header {
  margin-bottom: 1em;
  overflow: hidden;

  box-sizing: border-box;

  background-color: #f3f2f2;
  border: 1px solid #a7a7a7;

  height: fit-content;
  width: 100%;
  max-width: 450px;

  padding: 1em;
  border-radius: 4px;
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

@media screen and (max-width: 450px) {
  .result-header {
    width: 100%;
  }
}
</style>

<template>
  <fieldset class="result-header">
    <legend>results for</legend>
    <Error v-if="error" />
    <div class="result-header-name">{{ searchDetails.Name }}</div>
    <a
      class="result-header-url"
      v-bind:href="searchDetails.URL"
      target="_blank"
    >{{ searchDetails.URL }}</a>
  </fieldset>
</template>

<script>
import Error from "./Error.vue";
export default {
  name: "ResultSummary",
  components: {
    Error,
  },
  props: ["searchID"],
  data() {
    return {
      searchDetails: {},
      error: false,
    };
  },
  beforeMount: async function () {
    try {
      let initDetails = await this.getSearchDetails();
      this.searchDetails = initDetails;
    } catch (err) {
      this.error = true;
    }
  },
  methods: {
    getSearchDetails: async function () {
      const details = await this.$http.fetch(
        `/api/v1/search?ID=${this.$props.searchID}`
      );

      return details;
    },
  },
};
</script>