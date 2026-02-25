<script setup lang="ts">
import Card from 'primevue/card';
import Button from 'primevue/button';
import { RouterLink } from 'vue-router';
import { useDataStore } from '@/stores/data';

const store = useDataStore();
const emit = defineEmits(['view-report-card']);
const props = defineProps({
  library: {
    type: Object,
    required: true
  }
})

function description(library: any) {
  return store.renderMarkdown(library.info.markdown_desc);
}
function viewReportCard(library: any) {
  emit('view-report-card', library);
}
</script>

<template>
  <Card>
    <template #content>
      <div class="float-right">
        <span class="text-sm text-gray-500">{{ library.info.popularity }} <i class="pi pi-heart"></i></span>
      </div>
      <div class="markdown">
        <div v-html="description(library)"></div>
      </div>
    </template>
    <template #footer>
      <div class="flex justify-end gap-4">
        <Button asChild v-slot="slotProps" severity="secondary">
          <RouterLink :to="`/libraries/${library.info.key}`" :class="slotProps.class">Code Samples</RouterLink>
        </Button>
        <Button @click="viewReportCard(library)" severity="info">
          Report Card
        </Button>
      </div>
    </template>
  </Card>
</template>