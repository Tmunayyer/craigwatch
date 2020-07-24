<template>
  <Metric
    :metricname="'post / repost'"
    :data="privateState.data"
    :defaultSelected="'% reposts'"
    :error="privateState.error"
    :label="''"
  />
</template>

<script>
import Metric from "../Metric.vue";
import { resultPageState } from "../../Results.vue";

export default {
  name: "PostVsRepost",
  components: {
    Metric,
  },
  data() {
    return {
      privateState: {
        data: {},
        error: false,
      },
      sharedState: resultPageState.state,
    };
  },
  beforeMount: function () {
    this.privateState = this.transformData();
  },
  methods: {
    transformData: function () {
      const displayMap = {
        TotalCount: "# of posts",
        RepostCount: "# of reposts",
        PercentRepost: "% reposts",
      };

      let pairs = {};
      const { data, error } = this.sharedState.activityMetrics;
      for (const key in data) {
        if (displayMap[key] === undefined) continue;
        pairs[displayMap[key]] = data[key];
      }

      return { data: pairs, error: error };
    },
  },
};
</script>