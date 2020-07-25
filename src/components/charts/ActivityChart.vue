<style scoped>
</style>

<template>
  <SplineChart :seriesData="privateState.data" :error="privateState.error" />
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
    this.privateState = this.transformData();
  },
  data() {
    return {
      privateState: {
        data: [],
        error: false,
      },
      sharedState: resultPageState.state,
    };
  },
  methods: {
    transformData() {
      const { data, error } = this.sharedState.activityChart;
      const formatted = data.map((point, i) => {
        const date = util.chartDate(point.TopUnixDate * 1000);

        return {
          name: util.formatDate(point.TopUnixDate * 1000),
          x: i,
          y: point.Count,
          _label: date,
        };
      });
      return { data: formatted, error: error };
    },
  },
};
</script>