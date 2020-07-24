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
    this.privateState = this.setPrivateState();
  },
  methods: {
    setPrivateState: function () {
      const displayMap = {
        InSeconds: "seconds",
        InMinutes: "minutes",
        InHours: "hours",
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