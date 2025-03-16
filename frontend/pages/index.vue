<script setup lang="ts">
import { computed, ref } from "vue";

import "echarts";
import { getMostUsedDomains, getServerUsageByTimeRange } from "~/services/dashboard";

const selectedRange = ref<TimeRange>("Today");
const selectedCategories = ref<string[]>([]);
const isAllCategoriesSelected = ref<boolean>(true);

const { data: mostUsedDomains, status: statusMostUsedDomains } = await useAsyncData(
  () =>
    getMostUsedDomains({
      categories: selectedCategories.value,
      from: timeRangeValues.value.from,
      to: timeRangeValues.value.to,
    }),
  {
    server: false,
    watch: [selectedRange, selectedCategories],
  },
);

const { data: serverUsageByTimeRange, status: statusServerUsageByTimeRange } = await useAsyncData(
  () =>
    getServerUsageByTimeRange({
      categories: selectedCategories.value,
      from: timeRangeValues.value.from,
      to: timeRangeValues.value.to,
    }),
  {
    server: false,
    watch: [selectedRange, selectedCategories],
  },
);

const timeRangeValues = computed(() => {
  const { from, to } = getDateFromRange(selectedRange.value);
  return { from: from.getTime(), to: to.getTime() };
});

const serverUsageByTimeRangeOption = computed(() => {
  const dimensions = serverUsageByTimeRange.value
    ? serverUsageByTimeRange.value?.map((i) => formatDate(dateAtCurrentTZ(i.timeRange)))
    : [];
  const counts = serverUsageByTimeRange.value
    ? serverUsageByTimeRange.value?.map((i) => i.count)
    : [];

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
    series: [{ type: "line", data: counts }],
    tooltip: {},
    legend: {},
  };
});

const mostUsedDomainsOption = computed(() => {
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

watch(
  isAllCategoriesSelected,
  () => {
    if (isAllCategoriesSelected.value) {
      selectedCategories.value = getDomainCategories();
    } else {
      selectedCategories.value = [];
    }
  },
  {
    immediate: true,
  },
);
</script>

<template>
  <div class="filter-container">
    <v-select
      v-model="selectedRange"
      :items="['Last month', 'Last 2 weeks', 'Last week', 'Last 3 days', 'Yesterday', 'Today']"
      label="Range"
      variant="underlined"
      width="100%"
    />
    <v-checkbox
      v-model="isAllCategoriesSelected"
      label="Show all categories"
      min-width="250px"
    />
    <v-select
      v-model="selectedCategories"
      :disabled="isAllCategoriesSelected"
      :items="getDomainCategories()"
      label="Category"
      multiple
      variant="underlined"
      width="100%"
    >
      <template #selection="{ item, index }">
        <v-chip v-if="index < 4">
          <span>{{ item.title }}</span>
        </v-chip>
        <span
          v-if="index === 4"
          class="text-grey text-caption align-self-center"
        >
          (+{{ selectedCategories.length - 4 }} outros)
        </span>
      </template>
    </v-select>
  </div>
  <div class="dashboard-container">
    <v-sheet
      v-if="statusMostUsedDomains === 'success' || statusMostUsedDomains === 'pending'"
      elevation="4"
      height="300"
      width="100%"
    >
      <VChart
        autoresize
        :loading="statusMostUsedDomains === 'pending'"
        :option="mostUsedDomainsOption"
      />
    </v-sheet>
    <v-sheet
      v-if="
        statusServerUsageByTimeRange === 'success' || statusServerUsageByTimeRange === 'pending'
      "
      elevation="4"
      height="300"
      width="100%"
    >
      <VChart
        autoresize
        :loading="statusServerUsageByTimeRange === 'pending'"
        :option="serverUsageByTimeRangeOption"
      />
    </v-sheet>
  </div>
</template>

<style scoped>
.dashboard-container {
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 20px;
  width: 80%;
  margin: auto;
}

.filter-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  gap: 20px;
  width: 80%;
  margin: auto;
  margin-top: 30px;
}

@media (max-width: 768px) {
  .filter-container {
    flex-direction: column;
    width: 90%;
  }

  .dashboard-container {
    width: 90%;
  }
}
</style>
