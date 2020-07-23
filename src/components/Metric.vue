<style scoped>
.metric-container {
  margin-bottom: 1em;
  overflow: hidden;

  background-color: #f3f2f2;
  border: 1px solid #a7a7a7;

  width: fit-content;
  max-width: 375px;

  padding: 1em;
  border-radius: 4px;

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-content: center;
}

legend {
  margin-right: 2em;
}

.metric {
  text-align: center;

  font-size: 2em;
}

.measurement {
  font-size: 0.7em;
  margin-top: -1em;
  margin-bottom: 1em;
}

select {
  font-size: 1em;
}
</style>

<template>
  <fieldset class="metric-container">
    <legend>{{metricname}}</legend>

    <template v-if="error">
      <div class="metric">-</div>
    </template>
    <template v-else>
      <div class="measurement">
        <label>{{ computedLabel }}</label>
        <select name="granularity" v-model="selected">
          <option v-for="(_, name) in data" :key="name" :value="name">{{name}}</option>
        </select>
      </div>

      <div class="metric">{{data[selected]}}</div>
    </template>
  </fieldset>
</template>

<script>
export default {
  name: "Metric",
  props: ["label", "metricname", "data", "defaultSelected", "error"],
  data() {
    return {
      selected: this.$props.defaultSelected,
    };
  },
  computed: {
    computedLabel: function () {
      if (this.$props.label !== undefined) {
        return this.$props.label;
      } else {
        return "per";
      }
    },
  },
};
</script>