<script setup lang="ts">
import "echarts";
import { getMostUsedDomains } from "~/services/dashboard";

import { computed, ref } from "vue";

const selectedRange = ref<any>("Today");
const selectedCategories = ref<string[]>([]);

const { data: mostUsedDomains, status: statusMostUsedDomains } = await useAsyncData(
  () => getMostUsedDomains(selectedRange.value, selectedCategories.value),
  {
    server: false,
    watch: [selectedRange, selectedCategories],
  },
);

const option = computed(() => {
  const dimensions = mostUsedDomains.value ? mostUsedDomains.value?.map((i: any) => i.domain) : [];
  const counts = mostUsedDomains.value ? mostUsedDomains.value?.map((i: any) => i.count) : [];

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
      width="100%"
      label="Range"
      :items="['Last month', 'Last 2 weeks', 'Last week', 'Last 3 days', 'Yesterday', 'Today']"
      variant="underlined"
      v-model="selectedRange"
    ></v-select>
    <v-select
      width="100%"
      label="Category"
      :items="getDomainCategories()"
      variant="underlined"
      v-model="selectedCategories"
      multiple
      chips
    ></v-select>
  </div>
  <div class="dashboard-container">
    <v-sheet
      elevation="4"
      height="300"
      width="80%"
    >
      <VChart
        :option="option"
        v-if="statusMostUsedDomains === 'success'"
      />
    </v-sheet>
    <v-sheet
      elevation="4"
      height="300"
      width="80%"
    >
      <VChart
        :option="option"
        v-if="statusMostUsedDomains === 'success'"
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
