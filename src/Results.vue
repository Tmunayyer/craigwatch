<style module>
.page-container {
  box-sizing: border-box;
  width: 100%;
}

.metric-container {
  display: flex;
  flex-direction: row;

  justify-content: space-between;
  flex-wrap: wrap;

  width: 450px;
}

@media screen and (max-width: 450px) {
  .metric-container {
    width: 100%;
  }
}
</style>

<template>
  <div v-if="!loading" class="page-container">
    <ActivityChart />
    <ResultSummary v-bind:searchID="searchID" />
    <div class="metric-container">
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

import ActivityChart from "./components/charts/ActivityChart.vue";

export const resultPageState = {
  state: {
    activityMetrics: {},
    activityChart: [],
  },
  setActivityMetrics: function (data) {
    this.state.activityMetrics = data;
  },
  setActivityChart: function (data) {
    this.state.activityChart = data;
  },
};

export default {
  name: "Results",
  components: {
    ResultSummary,
    ResultList,
    ListingMetric,
    RepostMetric,
    PostVsRepost,
    ActivityChart,
  },
  beforeMount: async function () {
    const { ID } = this.$route.params;

    const activityMetrics = await this.$http.fetchRetry(
      `/api/v1/metric?ID=${ID}`,
      {},
      // retry if no data is returned
      (data) => data === undefined
    );
    resultPageState.setActivityMetrics(activityMetrics);

    const activityChart = await this.$http.fetchRetry(
      `/api/v1/activityChart?ID=${ID}`,
      {},
      (data) => data.length === 0
    );
    console.log("the activityChart:", activityChart);
    resultPageState.setActivityChart(activityChart);

    this.loading = false;
  },
  data() {
    return {
      searchID: this.$route.params.ID,
      loading: true,
    };
  },
};
</script>
