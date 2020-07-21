<template>
  <Metric
    :metricname="'reposts'"
    :data="privateState.data"
    :defaultSelected="'minutes'"
    :error="privateState.error"
  />
</template>

<script>
import Metric from "../Metric.vue";
import { resultPageState } from "../../Results.vue";

export default {
  name: "RepostMetric",
  components: {
    Metric
  },
  data() {
    return {
      privateState: {
        data: {},
        error: false
      },
      sharedState: resultPageState.state
    };
  },
  beforeMount: function() {
    this.privateState.data = this.transformData();
  },
  methods: {
    transformData: function() {
      const displayMap = {
        RepostInSeconds: "seconds",
        RepostInMinutes: "minutes",
        RepostInHours: "hours"
      };

      let pairs = {};
      for (const key in this.sharedState.activityData) {
        if (displayMap[key] === undefined) continue;
        pairs[displayMap[key]] = this.sharedState.activityData[key];
      }

      return pairs;
    }
  }
};
</script>