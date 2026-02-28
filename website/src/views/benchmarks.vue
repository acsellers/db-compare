<script setup lang="ts">
import { ref, computed } from 'vue';
import { useDataStore } from '@/stores/data';
import benchmarkData from '@/data/benchmarks.json';
import Chart from 'primevue/chart';
import SelectButton from 'primevue/selectbutton';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';

const store = useDataStore();

const metricOptions = [
  { label: 'Overall', value: 'overall' },
  { label: 'Sale Mix', value: 'sale_mix' },
  { label: 'Customer Updates', value: 'customer_updates' },
  { label: 'Reporting', value: 'reporting' },
  { label: 'Search', value: 'search' },
  { label: 'Bulk Import', value: 'bulk_import' }
];
const selectedMetric = ref('overall');

const modeOptions = [
  { label: 'Vs Std. Library', value: 'vs_stdlib' },
  { label: 'Absolute Time (ms)', value: 'time' }
];
const selectedMode = ref('vs_stdlib');

const libraries = Object.values(benchmarkData) as any[];

const chartData = computed(() => {
  const labels = [];
  const data = [];
  const bgColors = [];
  
  // Sort libraries by the current metric and mode to make the chart look nice
  const sortedLibs = [...libraries].sort((a, b) => {
    return getValue(b) - getValue(a); // Sort descending so best is at the top/right depending on chart, actually ascending is best for times
  });
  
  // Let's sort ascending (lower is better)
  sortedLibs.sort((a, b) => getValue(a) - getValue(b));

  for (const lib of sortedLibs) {
    const val = getValue(lib);
    const name = store.libraries[lib.key]?.info?.name || (lib.key === 'stdlib' ? 'Standard Library' : lib.key);
    labels.push(name);
    
    // adjust for stdlib in vs_stdlib mode 
    if (lib.key === 'stdlib' && selectedMode.value === 'vs_stdlib') {
      data.push(1);
    } else {
      data.push(val);
    }
    // Highlight stdlib
    bgColors.push(lib.key === 'stdlib' ? 'rgba(75, 192, 192, 0.8)' : 'rgba(54, 162, 235, 0.8)');
  }

  return {
    labels: labels,
    datasets: [
      {
        label: selectedMode.value === 'vs_stdlib' ? 'Multiplier vs Stdlib' : 'Time in ms',
        data: data,
        backgroundColor: bgColors,
        borderWidth: 1,
        borderRadius: 4
      }
    ]
  };
});

const chartOptions = computed(() => {
  return {
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: selectedMode.value === 'vs_stdlib' ? 'Multiplier (Lower is Better)' : 'Time in ms (Lower is Better)'
        }
      }
    }
  };
});

function getValue(lib: any) {
  if (selectedMetric.value === 'overall') {
    if (selectedMode.value === 'vs_stdlib') {
      return lib.overall;
    } else {
      // absolute time for overall is sum of items
      return lib.items.reduce((sum: number, item: any) => sum + item.time, 0);
    }
  } else {
    const item = lib.items.find((i: any) => i.name === selectedMetric.value);
    if (!item) return 0;
    if (selectedMode.value === 'vs_stdlib') {
      return item.vs_stdlib;
    } else {
      return item.time;
    }
  }
}

const tableData = computed(() => {
  return libraries.map(lib => {
    // Generate object containing all values
    const tableRow: Record<string, any> = {
      key: lib.key,
      name: store.libraries[lib.key]?.info?.name || (lib.key === 'stdlib' ? 'Standard Library' : lib.key),
      overall: selectedMode.value === 'vs_stdlib' ? (lib.key === 'stdlib' ? 1 : lib.overall) : lib.items.reduce((sum: number, item: any) => sum + item.time, 0),
    };
    
    // Add each metric
    for (const metric of metricOptions) {
      if (metric.value === 'overall') continue;
      
      const item = lib.items.find((i: any) => i.name === metric.value);
      if (item) {
        tableRow[metric.value] = selectedMode.value === 'vs_stdlib' ? (lib.key === 'stdlib' ? 1 : item.vs_stdlib) : item.time;
      } else {
        tableRow[metric.value] = 0;
      }
    }
    return tableRow;
  }).sort((a, b) => a[selectedMetric.value] - b[selectedMetric.value]);
});

function formatValue(value: number) {
  if (value === undefined || value === null) return '-';
  if (selectedMode.value === 'vs_stdlib') {
    return value.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2}) + 'x';
  } else {
    return value.toLocaleString(undefined, {maximumFractionDigits: 0}) + ' ms';
  }
}
</script>

<template>
  <div class="normal-container p-4">
    <h1 class="text-4xl font-bold mb-6 text-gray-800">Benchmarks</h1>
    <p class="mb-4 text-gray-600">
      Compare the performance of different Go database libraries across a variety of workloads.
      Use the controls below to switch between specific workloads and metric views.
    </p>

    <div class="bg-white p-6 shadow-md rounded-xl mb-8">
      <div class="flex flex-col xl:flex-row justify-between mb-6 gap-4">
        <SelectButton v-model="selectedMetric" :options="metricOptions" optionLabel="label" optionValue="value" :allowEmpty="false" />
        <SelectButton v-model="selectedMode" :options="modeOptions" optionLabel="label" optionValue="value" :allowEmpty="false" />
      </div>
      
      <div class="chart-container h-[50vh] min-h-[400px]">
        <Chart type="bar" :data="chartData" :options="chartOptions" class="h-full w-full" />
      </div>
    </div>

    <div class="bg-white p-6 shadow-md rounded-xl">
      <h2 class="text-2xl font-semibold mb-4 text-gray-800">Performance Data</h2>
      <DataTable :value="tableData" stripedRows responsiveLayout="scroll" :paginator="false" class="p-datatable-sm">
        <Column field="name" header="Library Name" sortable>
          <template #body="slotProps">
            <span class="font-medium" :class="{'text-primary': slotProps.data.key === 'stdlib'}">
              {{ slotProps.data.name }}
            </span>
          </template>
        </Column>
        <Column v-for="metric in metricOptions" :key="metric.value" :field="metric.value" :header="metric.label" sortable>
          <template #body="slotProps">
            <span>{{ formatValue(slotProps.data[metric.value]) }}</span>
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<style scoped>
.normal-container {
  max-width: 1200px;
  margin: 0 auto;
}
.chart-container {
  position: relative;
}
</style>
