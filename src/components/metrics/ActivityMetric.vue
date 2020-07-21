<template>
  <Metric :metricname="'activity'" :data="data" :defaultSelected="'minutes'" :error="error" />
</template>

<script>
import Metric from "../Metric.vue";
export default {
  name: "ActivityMetric",
  components: {
    Metric
  },
  data() {
    return {
      data: [],
      error: false
    };
  },
  beforeMount: async function() {
    try {
      const data = await this.getData();
      this.data = this.transformData(data);
    } catch (err) {
      this.error = true;
    }
  },
  methods: {
    getData: async function() {
      const data = await this.$http(
        `/api/v1/metric?ID=${this.$route.params.ID}`
      );
      return data;
    },
    transformData: function(obj) {
      const displayMap = {
        InSeconds: "seconds",
        InMinutes: "minutes",
        InHours: "hours"
      };

      let pairs = {};
      for (const key in obj) {
        if (key === "ID") continue;
        pairs[displayMap[key]] = obj[key];
      }

      return pairs;
    }
  }
};
</script>