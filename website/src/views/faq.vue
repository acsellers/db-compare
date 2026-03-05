<script setup lang="ts">
import Panel from 'primevue/panel'
</script>

<template>
  <div class="normal-container">
    <div class="flex flex-row justify-between mb-4">
      <h1 class="text-4xl">Frequently Asked Questions</h1>
    </div>
    <Panel header="What is this for?">
        <p>
          I wanted to compare how different database mapping libraries performed in developer experience, features, and performance.
          I thought I would wrap this up in a week or two... six weeks later, I'm finally happy with it.
        </p>
    </Panel>
    <Panel header="Is this AI Slop?">
        <p>
            While AI was used to build out some code (like the benchmark data or parts of the website), the sample code and documentation were written by humans. 
            
        </p>
    </Panel>
    <Panel header="What license are the samples under?">
        <p>
            Sample code is under the CC-0 license, which is their maximally permissive license. Otherwise, the content is under a 3 Clause BSD License.
        </p>
    </Panel>
  </div>
</template>