<style scoped>
</style>

<template>
  <Distribution :seriesData="privateState.data" :error="privateState.error" />
</template>

<script>
import Distribution from "../Distribution.vue";
import { resultPageState } from "../../Results.vue";

import util from "../../utility.js";

export default {
  name: "PriceDistribution",
  components: {
    Distribution,
  },
  beforeMount: function () {
    // technically unsafe but there wont be state changes after initial load
    const { data, error } = resultPageState.state.priceDistribution;
    this.privateState = {
      data: data,
      error: error,
    };
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
};
</script>