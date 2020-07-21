<style module>
.page-container {
  box-sizing: border-box;
  width: 100%;
}
</style>

<template>
  <div class="page-container">
    <ResultSummary v-bind:searchID="searchID" />
    <div>
      <ActivityMetric />
    </div>
    <br />
    <ResultList v-bind:searchID="searchID" />
  </div>
</template>

<script>
import ResultSummary from "./components/ResultSummary.vue";
import ResultList from "./components/ResultList.vue";
import ActivityMetric from "./components/metrics/ActivityMetric.vue";

export default {
  name: "Results",
  components: {
    ResultSummary,
    ResultList,
    ActivityMetric
  },
  data() {
    return {
      searchID: this.$route.params.ID
    };
  }
};
</script>
