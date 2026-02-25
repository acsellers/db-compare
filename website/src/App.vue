<script setup lang="ts">
import { ref } from "vue";
import Menubar from 'primevue/menubar';
import { useRouter } from 'vue-router';

const router = useRouter();

const items = ref([
    {
        label: 'Libraries',
        icon: 'pi pi-book',
        command: () => {
            router.push('/libraries');
        }
    },
    {
        label: 'Criteria',
        icon: 'pi pi-list',
        command: () => {
            router.push('/features');
        }
    },
    {
        label: 'Feature Matrix',
        icon: 'pi pi-table',
        command: () => {
            router.push('/matrix');
        }
    },
    {
        label: 'Samples',
        icon: 'pi pi-code',
        command: () => {
            router.push('/examples');
        }
    },
    {
        label: 'Benchmarks',
        icon: 'pi pi-chart-bar',
        command: () => {
            router.push('/benchmarks');
        }
    },

    {
        label: 'Coming Soon',
        icon: 'pi pi-clock',
        command: () => {
            router.push('/coming-soon');
        }
    }
]);
</script>

<template>
    <div class="layout-wrapper">
        <Menubar :model="items" class="navbar">
            <template #start>
                <span class="text-xl font-bold text-primary-700 mr-4 cursor-pointer" @click="router.push('/')">Go DB
                    Breakdown</span>
            </template>
            <template #end>
                <a href="https://github.com/acsellers/db-compare" class="text-primary-700 mr-4 cursor-pointer"><i
                        class="pi pi-github" style="font-size: 1.5rem;"></i></a>
            </template>
        </Menubar>
        <div class="layout-content mt-4">
            <RouterView />
        </div>
    </div>
</template>