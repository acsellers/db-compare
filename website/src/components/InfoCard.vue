<script setup lang="ts">
import markdownit from 'markdown-it';
const md = markdownit();
import Card from 'primevue/card';
import Button from 'primevue/button';
import { RouterLink } from 'vue-router';

const emit = defineEmits(['view-report-card']);
const props = defineProps({
    library: {
        type: Object,
        required: true
    }
})

function description(library: any) {
  return md.render(library.description);
}
function viewReportCard(library: any) {
  emit('view-report-card', library);
}
</script>

<template>
    <Card>
        <template #content>
          <div class="float-right">
            <span class="text-sm text-gray-500">{{ library.popularity }} <i class="pi pi-heart"></i></span>
          </div>
          <div class="markdown">
            <div v-html="description(library)"></div>
          </div>
        </template>
        <template #footer>
          <div class="flex justify-end gap-4">
            <Button asChild v-slot="slotProps" severity="secondary">
              <RouterLink :to="`/libraries/${library.id}`" :class="slotProps.class">Code Samples</RouterLink>
            </Button>
            <Button @click="viewReportCard(library)" severity="info">
              Report Card
            </Button>
          </div>
        </template>
    </Card>
</template>