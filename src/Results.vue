<style module>
.page-container {
  box-sizing: border-box;
  width: 100%;
}
</style>

<template>
  <div class="page-container">
    <ResultSummary v-bind:searchID="searchID" />
    <br />
    <ResultList v-bind:searchID="searchID" />
  </div>
</template>

<script>
import ResultSummary from "./components/ResultSummary.vue";
import ResultList from "./components/ResultList.vue";

export default {
  name: "Results",
  components: {
    ResultSummary,
    ResultList
  },
  data() {
    return {
      searchID: this.$route.params.ID
    };
  }
};
</script>
