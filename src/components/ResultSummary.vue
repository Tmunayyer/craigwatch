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
export default {
  name: "ResultSummary",
  props: ["searchDetails"]
};
</script>