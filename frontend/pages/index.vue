<script setup lang="ts">
import { computed, ref } from "vue";

import "echarts";
import { getMostUsedDomains } from "~/services/dashboard";

const selectedRange = ref<string>("Today");
const selectedCategories = ref<string[]>([]);
const isAllCategoriesSelected = ref<boolean>(true);

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
      v-if="statusMostUsedDomains === 'success'"
      elevation="4"
      height="300"
      width="100%"
    >
      <VChart :option="option" />
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
