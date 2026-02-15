<script setup lang="ts">
import { ref, onMounted } from 'vue';
import docs from '@/data/docs.json';
import { Subjects } from '@/data/subjects';
import markdownit from 'markdown-it';

const md = markdownit();

const examples = ref<Record<string, string[]>>({});
onMounted(() => {
    Subjects.forEach((subject: any) => {
        examples.value[subject.name] = [];
        subject.subjects.forEach((sub: string) => {
            console.log(sub, docs.Examples[sub])
            examples.value[subject.name].push(md.render(docs.Examples[sub]));
        })
    })
})
</script>

<template>
  <div class="normal-container">
    <h1 class="text-3xl font-bold mb-6">Subject Breakdown</h1>
    <div class="flex flex-col gap-3">
      <p>In order to evaluate libraries across a number of different features, I've broken down the features into these "subjects".</p>
      <div class="markdown" v-for="(example, subject) in examples" :key="subject">
        <h2>{{ subject }}</h2>
        <div v-for="example in example" :key="example" v-html="example"></div>
      </div>
    </div>
  </div>
</template>
<style scoped>
a {
    text-decoration: none;
}
</style>