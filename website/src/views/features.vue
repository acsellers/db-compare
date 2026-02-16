<script setup lang="ts">
import { useDataStore } from '@/stores/data';
import Panel from 'primevue/panel';

const store = useDataStore();


function toKey(s: any) {
  let key = s['title'] || s['name'];
  return key.toLowerCase().
    replaceAll(' ', '_').
    replaceAll('/', '_')

}
function renderMarkdown(markdown: string) {
  return store.renderMarkdown(markdown);
}
</script>

<template>
  <div class="normal-container">
    <h1 class="text-3xl font-bold mb-6">Evaluation Criteria</h1>
    <div class="flex flex-col gap-3">
      <p>
        In order to evaluate libraries across a number of different
        features, I've broken down the features into these criteria. Each
        section has an example of an ideal implementation. Some subjects
        will have multiple sub-criteria, but they all get sort of averaged
        together to get a final score.
      </p>
      <div v-for="feature in store.features" :key="feature.name">
        <h2 class="text-2xl font-bold mb-2">{{ feature.name }}</h2>
        <Panel v-for="subject in feature.subjects" :key="toKey(subject)" :header="subject.title" toggleable collapsed>
          <div class="markdown" v-html="renderMarkdown(subject.description)"></div>
          <div v-for="example in subject.sub_examples" :key="toKey(example)" class="mt-5">
            <h4 class="text-xl font-bold mb-2">{{ example.title }}</h4>
            <div class="markdown" v-html="renderMarkdown(example.description)"></div>
          </div>
        </Panel>
      </div>
    </div>
  </div>
</template>
<style scoped>
a {
  text-decoration: none;
}
</style>