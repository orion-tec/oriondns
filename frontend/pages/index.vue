import { computed, onMounted, reactive, ref, watch } from "vue";

<script setup lang="ts">
import "echarts";

const dashboardData = ref(null);

onMounted(() => {
  fetch(`/api/v1/dashboard/most-used-domains`)
    .then((response) => response.json())
    .then((data) => {
      dashboardData.value = data;
    });
});

const option = computed(() => {
  const dimensions = dashboardData.value
    ? dashboardData.value?.map((i) => i.domain)
    : [];

  const counts = dashboardData.value
    ? dashboardData.value?.map((i) => i.count)
    : [];

  return {
    xAxis: {
      type: "category",
      data: dimensions,
      name: "Domains",
      nameTextStyle: { color: "white" },
      axisLabel: {
        rotate: 30,
        color: "white",
        overflow: "truncate",
        width: "95",
      },
    },
    yAxis: {
      type: "value",
      name: "Count",
      axisLabel: { color: "white" },
      nameTextStyle: { color: "white" },
    },
    series: [{ type: "bar", data: counts }],
    tooltip: {},
    legend: {},
  };
});
</script>

<template>
  <v-card text="Most used domains" class="chart-container">
    <VChart :option="option" />
  </v-card>
</template>

<style scoped>
.chart-container {
  width: 80%;
  height: 50%;
  margin: 30px auto;
  padding: 0 0 50px 0;
}
</style>
