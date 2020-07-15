<style module>
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
</style>

<template>
  <div class="result-header">
    <Error v-if="error" />
    <div class="result-header-lead">Results For:</div>
    <div class="result-header-name">{{ searchDetails.Name }}</div>
    <a
      class="result-header-url"
      v-bind:href="searchDetails.URL"
      target="_blank"
    >{{ searchDetails.URL }}</a>
  </div>
</template>

<script>
import Error from "./Error.vue";
export default {
  name: "ResultSummary",
  components: {
    Error
  },
  props: ["searchID"],
  data() {
    return {
      searchDetails: {},
      error: false
    };
  },
  beforeMount: async function() {
    try {
      let initDetails = await this.getSearchDetails();
      this.searchDetails = initDetails;
    } catch (err) {
      this.error = true;
    }
  },
  methods: {
    getSearchDetails: async function() {
      const details = await this.$http(
        `/api/v1/search?ID=${this.$props.searchID}`
      );

      return details;
    }
  }
};
</script>