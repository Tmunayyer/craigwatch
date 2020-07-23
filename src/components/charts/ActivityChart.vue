<style scoped>
</style>

<template>
  <SplineChart :seriesData="privateState.data" />
</template>

<script>
import highcharts from "highcharts";
import SplineChart from "../SplineChart.vue";
import { resultPageState } from "../../Results.vue";

import util from "../../utility.js";

export default {
  name: "ActivityChart",
  components: {
    SplineChart,
  },
  beforeMount: function () {
    this.privateState.data = this.transformData();
  },
  data() {
    return {
      privateState: {
        data: [],
      },
      sharedState: resultPageState.state,
    };
  },
  methods: {
    transformData() {
      const formatted = this.sharedState.activityChart.map((point, i) => {
        const date = util.formatDate(point.TopUnixDate * 1000);

        return {
          name: date,
          x: i,
          y: point.Count,
          _label: util.chartDate(date),
        };
      });
      return formatted;
    },
  },
};
</script>