<script setup lang="ts">
const option = ref({
  dataset: {
    dimensions: ["Product", "2015", "2016", "2017"],
    source: [
      {
        Product: "Matcha Latte",
        2015: 54,
        2016: 42,
        2017: 23,
      },
    ],
  },
  xAxis: { type: "category" },
  yAxis: {},
  series: [{ type: "bar" }],
});
</script>

<template>
  <div class="chart-container">
    <VChart :option="option" />
  </div>
</template>

<style scoped>
.chart-container {
  width: 50%;
  height: 50%;
  margin: 0 auto;
}
</style>
