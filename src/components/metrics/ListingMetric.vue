<template>
  <Metric
    :metricname="'listing'"
    :data="privateState.data"
    :defaultSelected="'minutes'"
    :error="privateState.error"
  />
</template>

<script>
import Metric from "../Metric.vue";
import { resultPageState } from "../../Results.vue";

export default {
  name: "ListingMetric",
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
    this.privateState.data = this.transformData();
  },
  methods: {
    transformData: function () {
      const displayMap = {
        InSeconds: "seconds",
        InMinutes: "minutes",
        InHours: "hours",
      };

      let pairs = {};
      for (const key in this.sharedState.activityMetrics) {
        if (displayMap[key] === undefined) continue;
        pairs[displayMap[key]] = this.sharedState.activityMetrics[key];
      }

      return pairs;
    },
  },
};
</script>