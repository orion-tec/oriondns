<script setup lang="ts">
import { computed, ref } from "vue";

import "echarts";
import { getMostUsedDomains } from "~/services/dashboard";

const selectedRange = ref<string>("Today");
const selectedCategories = ref<string[]>([]);

const { data: mostUsedDomains, status: statusMostUsedDomains } = await useAsyncData(
  () =>
    getMostUsedDomains({
      categories: selectedCategories.value,
      range: selectedRange.value,
    }),
  {
    server: false,
    watch: [selectedRange, selectedCategories],
  },
);

const option = computed(() => {
  const dimensions = mostUsedDomains.value ? mostUsedDomains.value?.map((i) => i.domain) : [];
  const counts = mostUsedDomains.value ? mostUsedDomains.value?.map((i) => i.count) : [];

  return {
    xAxis: {
      type: "category",
      data: dimensions,
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
  <div class="filter-container">
    <v-select
      v-model="selectedRange"
      width="100%"
      label="Range"
      :items="['Last month', 'Last 2 weeks', 'Last week', 'Last 3 days', 'Yesterday', 'Today']"
      variant="underlined"
    />
    <v-select
      v-model="selectedCategories"
      width="100%"
      label="Category"
      :items="getDomainCategories()"
      variant="underlined"
      multiple
      chips
    />
  </div>
  <div class="dashboard-container">
    <v-sheet
      elevation="4"
      height="300"
      width="80%"
    >
      <VChart
        v-if="statusMostUsedDomains === 'success'"
        :option="option"
      />
    </v-sheet>
    <v-sheet
      elevation="4"
      height="300"
      width="80%"
    >
      <VChart
        v-if="statusMostUsedDomains === 'success'"
        :option="option"
      />
    </v-sheet>
  </div>
</template>

<style scoped>
.dashboard-container {
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  align-items: center;
  gap: 20px;
}

.filter-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  gap: 20px;
  width: 80%;
  margin: auto;
}
</style>
