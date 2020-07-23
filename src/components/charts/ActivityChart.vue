<style scoped>
</style>

<template>
  <SplineChart :seriesData="privateState.data" />
</template>

<script>
import SplineChart from "../SplineChart.vue";
import { resultPageState } from "../../Results.vue";

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
      const formatted = this.sharedState.activityChart.map((point) => {
        return {
          y: point.Count,
          name: point.Label,
        };
      });
      return formatted;
    },
  },
};
</script>