<script setup lang="ts">
import { ref, computed } from 'vue';
import reportCardsData from '@/data/report_cards.json';
import { useRoute } from 'vue-router';
import ReportCard from '@/components/ReportCard.vue';

const route = useRoute();
const reportCards = ref(Object.values(reportCardsData));
const card = computed(() => {
    // Flatten the nested structure where key is the name
    // Actually reportCardsData is { "bob": { ... }, ... }
    // Object.values gives [{ name: "bob", ... }]
    // So finding by name should work if route param matches name.
    return reportCards.value.find((c: any) => c.name === route.params.name);
});

</script>

<template>
    <div class="flex justify-center">
        <div v-if="card" class="bg-white shadow-lg rounded-lg overflow-hidden w-full max-w-4xl border border-gray-200" style="font-family: 'Courier New', Courier, monospace;">
            <ReportCard :card="card" />
        </div>

        <div v-else class="text-center mt-20">
            <h2 class="text-2xl font-bold text-gray-700">Report Card Not Found</h2>
            <p class="text-gray-500">The requested library could not be found.</p>
        </div>
    </div>
</template>

<style scoped>
/* Scoped styles can be minimal since we are using Tailwind classes */
</style>