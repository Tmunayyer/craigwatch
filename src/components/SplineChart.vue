<style scoped>
.chart-container {
  box-sizing: border-box;

  margin-bottom: 1em;
  overflow: hidden;

  background-color: #f3f2f2;
  border: 1px solid #a7a7a7;

  width: 100%;
  max-width: 600px;

  padding: 1em;
  padding-top: unset;
  border-radius: 4px;

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-content: center;
}
</style>

<template>
  <fieldset class="chart-container">
    <legend>activity</legend>
    <Chart :options="chartOptions" />
  </fieldset>
</template>

<script>
import { Chart } from "highcharts-vue";

export default {
  name: "ActivityChart",
  components: {
    Chart,
  },
  props: {
    seriesData: Array,
    error: Boolean,
  },
  beforeMount() {
    this.chartOptions.yAxis.tickInterval = this.setYAxisRange();
  },
  data() {
    return {
      chartOptions: {
        chart: {
          type: "spline",
          height: (9 / 16) * 100 + "%",
          backgroundColor: "#f3f2f2",
          marginLeft: 64,
          marginBottom: 38,
        },
        title: {
          text: "Past 48 Hours",
          floating: false,
          margin: 10,
          style: {
            fontSize: "1.2em",
          },
        },
        series: [
          {
            type: "spline",
            data: this.seriesData,
          },
        ],
        plotOptions: {
          spline: {
            color: "#1A1919",
          },
          series: {
            marker: {
              enabled: false,
            },
          },
        },
        yAxis: {
          title: {
            text: "# of listings",
          },
        },
        xAxis: {
          title: {
            text: "hour",
            x: -30,
          },
          labels: {
            style: {
              fontSize: "0.7em",
            },
            formatter: (() => {
              const seriesData = this.seriesData;
              return function () {
                const index = this.value;
                const value = seriesData[index]._label;
                return seriesData[index]._label;
              };
            })(),
          },
        },
        credits: {
          enabled: false,
        },
        legend: {
          enabled: false,
        },
        tooltip: {
          useHTML: true,
          formatter: function () {
            return `
              <span>
                <b>Date:</b> ${this.key}
                <br/>
                <b>Count:</b> ${this.y}
              </span>
            `;
          },
        },
      },
    };
  },
  methods: {
    setYAxisRange() {
      const { seriesData } = this;
      let max = 0;
      seriesData.forEach((elem) => {
        if (elem.y > max) {
          max = elem.y;
        }
      });

      return Math.round(max / 4);
    },
  },
};
</script>