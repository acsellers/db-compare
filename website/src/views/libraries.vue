<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Card from 'primevue/card';
import Button from 'primevue/button';
import docs from '@/data/docs.json';
import markdownit from 'markdown-it';
import Dialog from 'primevue/dialog';
import ReportCard from '@/components/ReportCard.vue';

const md = markdownit();

interface Library {
  id: string;
  name: string;
  description: string;
}
const libs = ref<Library[]>([]);
onMounted(() => {
  Object.keys(docs.Libraries).forEach((key:string) => {
    libs.value.push({
      id: key,
      name: docs.ReportCards[key].name,
      description: docs.Libraries[key]
    })
  })
})

function viewReportCard(library: any) {
  reportCard.value = docs.ReportCards[library.id];
  showReportCard.value = true;
}
function description(library: any) {
  return md.render(library.description);
}
const reportCard = ref<any>(null);
const showReportCard = ref(false);
</script>

<template>
  <div class="normal-container">
    <h1 class="text-3xl font-bold mb-6">Libraries</h1>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
      <Card v-for="library in libs" :key="library.id">
        <template #title>
          {{ library.name }}
        </template>
        <template #content>
          <div v-html="description(library)"></div>
        </template>
        <template #footer>
          <div class="flex gap-4">
            <Button asChild v-slot="slotProps">
              <RouterLink :to="`/libraries/${library.id}`" :class="slotProps.class">Code Samples</RouterLink>
            </Button>
            <Button @click="viewReportCard(library)">
              Report Card
            </Button>
          </div>
        </template>
      </Card>
    </div>
  </div>
  <Dialog v-model:visible="showReportCard" modal pt:root:class="!border-0 !bg-transparent" pt:mask:class="backdrop-blur-sm">
    <template #container="{ closeCallback }">
      <div class="bg-white">
        <ReportCard :card="reportCard" />
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