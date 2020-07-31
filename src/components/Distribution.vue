<style scoped>
.chart-container {
  box-sizing: border-box;

  margin-bottom: 1em;
  overflow: hidden;

  background-color: #f3f2f2;
  border: 1px solid #a7a7a7;

  width: 100%;
  max-width: 450px;

  padding: 1em;
  padding-top: unset;
  border-radius: 4px;
  margin: 0 auto;

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-content: center;
}
</style>

<template>
  <fieldset class="chart-container">
    <legend>price</legend>
    <Chart :options="chartOptions" />
  </fieldset>
</template>

<script>
import { Chart } from "highcharts-vue";
import Highcharts from "highcharts";
import bellcurve from "highcharts/modules/histogram-bellcurve.js";
bellcurve(Highcharts);

// prettier-ignore
var data = [3.5, 3, 3.2, 3.1, 3.6, 3.9, 3.4, 3.4, 2.9, 3.1, 3.7, 3.4, 3, 3, 4,
    4.4, 3.9, 3.5, 3.8, 3.8, 3.4, 3.7, 3.6, 3.3, 3.4, 3, 3.4, 3.5, 3.4, 3.2,
    3.1, 3.4, 4.1, 4.2, 3.1, 3.2, 3.5, 3.6, 3, 3.4, 3.5, 2.3, 3.2, 3.5, 3.8, 3,
    3.8, 3.2, 3.7, 3.3, 3.2, 3.2, 3.1, 2.3, 2.8, 2.8, 3.3, 2.4, 2.9, 2.7, 2, 3,
    2.2, 2.9, 2.9, 3.1, 3, 2.7, 2.2, 2.5, 3.2, 2.8, 2.5, 2.8, 2.9, 3, 2.8, 3,
    2.9, 2.6, 2.4, 2.4, 2.7, 2.7, 3, 3.4, 3.1, 2.3, 3, 2.5, 2.6, 3, 2.6, 2.3,
    2.7, 3, 2.9, 2.9, 2.5, 2.8, 3.3, 2.7, 3, 2.9, 3, 3, 2.5, 2.9, 2.5, 3.6,
    3.2, 2.7, 3, 2.5, 2.8, 3.2, 3, 3.8, 2.6, 2.2, 3.2, 2.8, 2.8, 2.7, 3.3, 3.2,
    2.8, 3, 2.8, 3, 2.8, 3.8, 2.8, 2.8, 2.6, 3, 3.4, 3.1, 3, 3.1, 3.1, 3.1, 2.7,
    3.2, 3.3, 3, 2.5, 3, 3.4, 3];

export default {
  name: "Distribution",
  components: {
    Chart,
  },
  props: {
    seriesData: Object, // AveragePrice, SampleSize, DataSet
    error: Boolean,
  },
  beforeMount: function () {
    this.chartOptions.series[1].data = this.seriesData.DataSet;
  },
  data() {
    return {
      chartOptions: {
        title: {
          text: "distribution sample of n",
          floating: false,
          margin: 10,
          style: {
            fontSize: "1.2em",
          },
        },

        chart: {
          height: (9 / 16) * 100 + "%",
          backgroundColor: "#f3f2f2",
          marginLeft: 64,
          marginBottom: 38,
        },

        xAxis: [
          {
            title: {
              text: "Data",
            },
            alignTicks: false,
          },
          {
            title: {
              text: "Bell curve",
            },
            alignTicks: false,
            opposite: true,
          },
        ],

        yAxis: [
          {
            title: { text: "Data" },
          },
          {
            title: { text: "Bell curve" },
            opposite: true,
          },
        ],

        legend: {
          enabled: false,
        },

        credits: {
          enabled: false,
        },

        series: [
          {
            name: "Bell curve",
            type: "bellcurve",
            xAxis: 1,
            yAxis: 1,
            baseSeries: 1,
            zIndex: -1,
          },
          {
            name: "Data",
            type: "scatter",
            data: undefined,
            accessibility: {
              exposeAsGroupOnly: true,
            },
            marker: {
              radius: 1.5,
            },
          },
        ],
      },
    };
  },
};
</script>