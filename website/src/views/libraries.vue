<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import Button from 'primevue/button';
import { useDataStore } from '@/stores/data';
import Dialog from 'primevue/dialog';
import ReportCard from '@/components/ReportCard.vue';
import InfoCard from '@/components/InfoCard.vue';
import type { FullLibrary } from '@/data/fields';

interface Library {
  id: string;
  name: string;
  description: string;
  popularity: number;
}
const store = useDataStore();
const databases = ref('any');
const databaseOptions = [
  { label: 'Any Database', value: 'any' },
  { label: 'Postgres', value: 'postgres' },
  { label: 'MySQL', value: 'mysql' },
  { label: 'Sqlite', value: 'sqlite' }
]
const sortBy = ref('popularity');
const sortOptions = [
  { label: 'Github Stars', value: 'popularity' },
  { label: 'Name', value: 'name' },
  { label: 'Activity', value: 'activity' }
]
const libType = ref('all');
const libTypeOptions = [
  { label: 'All', value: 'all' },
  { label: 'ORM', value: 'orm' },
  { label: 'Query Builder', value: 'query_builder' },
  { label: 'Generic Mapper', value: 'mapper' },
  { label: 'Generated Mapper', value: 'generated_mapper' },
  { label: 'Generated ORM', value: 'generated_orm' }
]
const libs = computed(() => {
  let list = Object.values(store.libraries);
  if (databases.value !== 'any') {
    list = list.filter((lib: any) => lib.databases.includes(databases.value));
  }
  if (libType.value !== 'all') {
    list = list.filter((lib: any) => lib.types.includes(libType.value));
  }
  if (sortBy.value === 'popularity') {
    list = list.sort((a: any, b: any) => b.popularity - a.popularity);
  } else if (sortBy.value === 'name') {
    list = list.sort((a: any, b: any) => a.name.localeCompare(b.name));
  } else if (sortBy.value === 'activity') {
    list = list.sort((a: any, b: any) => b.activity - a.activity);
  }
  return list;
})

function viewReportCard(library: FullLibrary) {
  selectedLibrary.value = library;
  showReportCard.value = true;
}

const selectedLibrary = ref<FullLibrary | null>(null);
const showReportCard = ref(false);
</script>

<template>
  <div class="normal-container">
    <div class="flex flex-row justify-between">
      <h1 class="text-3xl font-bold mb-6">Libraries</h1>
      <div>
      </div>
    </div>
    <div class="flex flex-col gap-3">
      <InfoCard v-for="library in libs" :key="library.info.key" :library="library" @view-report-card="viewReportCard" />
    </div>
  </div>
  <Dialog v-model:visible="showReportCard" modal pt:root:class="!border-0 !bg-transparent"
    pt:mask:class="backdrop-blur-sm">
    <template #container="{ closeCallback }">
      <div class="bg-white max-h-[90vh] overflow-y-auto">
        <ReportCard :library="selectedLibrary" />
      </div>
      <div class="close-button fixed -right-4 -top-4">
        <Button @click="closeCallback" icon="pi pi-times" rounded severity="secondary" />
      </div>
    </template>
  </Dialog>
</template>
<style scoped>
a {
  text-decoration: none;
}
</style>