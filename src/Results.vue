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
}

/* 
  Media queries to handle dynamic page layout. Greater than
  450 px should have 2 column. less than should have 1.

  this is a little tricky since I wanted the column items to be "out of order"
  or the middle items (metrics) to be on the right when > 450 but for them to be
  placed before the list under 450 px
 */

@media screen and (min-width: 902px) {
  .metric-container {
    margin: 0 auto;
  }

  /* 
    Set up a container that will space evenly. Columns will then be given
    a 50% width so no more than 2 items will show up on a given column.

    Make sure to wrap the row so that anything after the "right" column
    will appear on the left side
   */
  .column-container {
    display: flex;
    flex-direction: row;
    align-content: space-between;
    flex-flow: row wrap;
  }

  /* 
    column widths

    make the right column height 0 so trailing "left" column items will be
    directly under things in left column
   */

  .column.pre-left {
    padding-right: 1em;
    width: 45%;
  }

  .column.right {
    padding-left: 1em;
    width: 45%;
    height: 0px;
    position: sticky;
    top: 15px;
  }

  .column.post-left {
    padding-right: 1em;
    width: 45%;
  }
}

@media screen and (max-width: 901px) {
  .metric-container {
    width: 100%;

    width: 450px;
    max-width: 450px;
  }

  .column-container {
    display: unset;
    flex-direction: unset;
    align-content: unset;
    flex-flow: unset;
  }

  .column.pre-left {
    width: 100%;
  }

  .column.right {
    width: 100%;
    height: fit-content;
  }

  .column.post-left {
    width: 100%;
  }
}
</style>

<template>
  <div v-if="!loading" class="page-container">
    <div class="column-container">
      <div class="column pre-left">
        <ResultSummary v-bind:searchID="searchID" />
      </div>
      <div class="column right">
        <div class="metric-container">
          <ListingMetric />
          <RepostMetric />
          <PostVsRepost />
          <ActivityChart />
        </div>
      </div>
      <div class="column post-left">
        <br />
        <ResultList v-bind:searchID="searchID" />
      </div>
    </div>
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
    activityChart: {
      data: [],
      error: false,
    },
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

    try {
      const activityMetrics = await this.$http.fetchRetry(
        `/api/v1/metric?ID=${ID}`,
        {},
        // retry if no data is returned
        (data) => data === undefined
      );
      resultPageState.setActivityMetrics({
        data: activityMetrics,
        error: false,
      });
    } catch (err) {
      reusltPageState.setActivityMetrics({ data: {}, error: true });
    }

    try {
      const activityChart = await this.$http.fetchRetry(
        `/api/v1/activityChart?ID=${ID}`,
        {},
        (data) => data.length === 0
      );
      resultPageState.setActivityChart({ data: activityChart, error: false });
    } catch (err) {
      resultPageState.setActivityChart({ data: [], error: true });
    }

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
