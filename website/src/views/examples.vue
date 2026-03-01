<script setup lang="ts">
import SampleView from '@/components/SampleView.vue'
</script>

<template>
  <div class="normal-container">
    <div class="flex flex-row justify-between mb-4">
      <h1 class="text-4xl">Samples</h1>
    </div>
    <div class="card">
      <SampleView lib="bob" db="mysql" sample="create-sale" displayMode="highlights" />
    </div>
  </div>
</template>