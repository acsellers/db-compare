<script setup lang="ts">
import Button from 'primevue/button';
import { useRouter } from 'vue-router';
import Card from 'primevue/card';

const router = useRouter();
</script>

<template>
    <div class="normal-container no-bg">
        <Card>
            <template #title>
                <span class="text-4xl text-primary-700 font-bold py-0 my-0">Go Database Library Breakdown</span>
            </template>
            <template #content>
                <div class="grid grid-cols-1 md:grid-cols-2">
                    <div class="block md:hidden">
                        <img src="@/assets/gopher_riding_server.png" alt="Gopher riding server" class="w-full border-round-2xl shadow-4" style="max-width: 600px;" />
                    </div>
                    <div class="">
                        <h2 class="text-2xl mb-6">Compare popular database libraries for Go</h2>
                        <p class="text-xl mb-6">See working code examples across Grouping, Bulk Loading, Views, Joins, and more for libraries across 9 feature areas.</p>
                        
                    </div>
                    <div class="hidden md:block">
                        <img src="@/assets/gopher_riding_server.png" alt="Gopher riding server" class="w-full border-round-2xl shadow-4" style="max-width: 600px;" />
                    </div>
                </div>
                
            </template>
        </Card>
        <section class="features py-8">
            <div class="container">
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div class="p-3">
                        <div class="shadow-2 bg-white p-4 border-round h-full text-center">
                            <i class="pi pi-chart-line text-primary text-5xl mb-3"></i>
                            <h3 class="text-xl font-bold text-primary-800 mb-2"><RouterLink to="/benchmarks">Performance Benchmarks</RouterLink></h3>
                            <p class="text-600 line-height-3 m-0">
                                Performance comparisons to help you avoid the wrong library.
                            </p>
                        </div>
                    </div>
                    <div class="p-3">
                        <div class="shadow-2 bg-white p-4 border-round h-full text-center">
                            <i class="pi pi-list text-primary text-5xl mb-3"></i>
                            <h3 class="text-xl font-bold text-primary-800 mb-2"><RouterLink to="/matrix">Feature Matrix</RouterLink></h3>
                            <p class="text-600 line-height-3 m-0">
                                Side-by-side comparison of features, support, and ecosystem compatibility.
                            </p>
                        </div>
                    </div>
                    <div class="p-3">
                        <div class="shadow-2 bg-white p-4 border-round h-full text-center">
                            <i class="pi pi-code text-primary text-5xl mb-3"></i>
                            <h3 class="text-xl font-bold text-primary-800 mb-2"><RouterLink to="/examples">Code Samples</RouterLink></h3>
                            <p class="text-600 line-height-3 m-0">
                                Real-world usage examples to demonstrate API design and developer experience.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped>
.p-card {
    background-color: white;
}
.text-900 {
    color: var(--surface-900);
}

.text-600 {
    color: var(--surface-600);
}

.text-primary {
    color: var(--p-primary-color);
}

.shadow-2 {
    box-shadow: 0 4px 10px rgba(0,0,0,0.03), 0 0 2px rgba(0,0,0,0.06), 0 2px 6px rgba(0,0,0,0.06);
}

.border-round {
    border-radius: var(--border-radius);
}

.h-full {
    height: 100%;
}
</style>
