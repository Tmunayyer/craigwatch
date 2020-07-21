<style module>
.page-container {
  box-sizing: border-box;
  width: 100%;
}
</style>

<template>
  <div v-if="!loading" class="page-container">
    <ResultSummary v-bind:searchID="searchID" />
    <div>
      <ListingMetric />
      <RepostMetric />
      <PostVsRepost />
    </div>
    <br />
    <ResultList v-bind:searchID="searchID" />
  </div>
</template>

<script>
import ResultSummary from "./components/ResultSummary.vue";
import ResultList from "./components/ResultList.vue";
import ListingMetric from "./components/metrics/ListingMetric.vue";
import RepostMetric from "./components/metrics/RepostMetric.vue";
import PostVsRepost from "./components/metrics/PostVsRepost.vue";

export const resultPageState = {
  state: {
    activityData: {}
  },
  setActivityData: function(data) {
    this.state.activityData = data;
  }
};

export default {
  name: "Results",
  components: {
    ResultSummary,
    ResultList,
    ListingMetric,
    RepostMetric,
    PostVsRepost
  },
  beforeMount: async function() {
    const { ID } = this.$route.params;

    const activityData = await this.$http.fetchRetry(
      `/api/v1/metric?ID=${ID}`,
      {},
      // retry if no data is returned
      data => data === undefined
    );
    resultPageState.setActivityData(activityData);

    this.loading = false;
  },
  data() {
    return {
      searchID: this.$route.params.ID,
      loading: true
    };
  }
};
</script>
