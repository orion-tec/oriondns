<script setup lang="ts">
import { Line } from "vue-chartjs";
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  CategoryScale,
  LinearScale,
  Colors,
  PointElement,
  LineElement,
  scales,
} from "chart.js";
import type { callback } from "chart.js/helpers";

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  CategoryScale,
  LinearScale,
  Colors,
  PointElement,
  LineElement,
);

const now = new Date();
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
const yesterday = new Date(today);
yesterday.setDate(yesterday.getDate() - 1);

const tomorrow = new Date(today);
tomorrow.setDate(tomorrow.getDate() + 1);

const renderDate = (date: Date) => {
  return `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()}`;
};

const labels = [yesterday, today, tomorrow];

const chartData = {
  labels,
  datasets: [
    {
      label: "Random data",
      data: [40, 20, 50],
      tension: 0.5,
      fill: false,
      borderColor: "red",
    },
  ],
};

const chartOptions = {
  responsive: true,
  scales: {
    x: {
      display: true,
      title: {
        display: true,
        text: "Date",
      },
      ticks: {
        callback: (value: number | string) => {
          return renderDate(labels[value]);
        },
      },
    },
  },
  plugins: {
    title: {
      display: true,
      text: "Chart.js Line Chart",
    },
    legend: {
      labels: {
        color: "white",
      },
      title: {
        text: "Legend",
      },
    },
  },
};
</script>

<template>
  <div class="chart-container">
    <Line id="my-chart-id" :options="chartOptions" :data="chartData" />
  </div>
</template>

<style scoped>
.chart-container {
  width: 50%;
  height: 50%;
  margin: 0 auto;
}
</style>
