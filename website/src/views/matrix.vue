<script setup lang="ts">
import { ref, computed } from 'vue';
import { useDataStore } from '../stores/data.ts'
import DataTable from 'primevue/datatable';
import Column from 'primevue/column'; 
import goldIcon from '@/assets/gold.svg';
import silverIcon from '@/assets/silver.svg';
import bronzeIcon from '@/assets/bronze.svg';
import todoIcon from '@/assets/todo.svg';
import naIcon from '@/assets/na.svg';
import failIcon from '@/assets/fail.svg';
import type { FullLibrary } from '../data/fields.ts';

const store = useDataStore();

interface CheckCategory {
    name: string;
    subChecks: CheckItem;
}
interface CheckItem {
    key: string;
    name: string;
    rating: (l: string) => string;
    notes?: (l: string) => string;
}

function genSubCheck(key: string, name: string, src: string) {
    return {
        key: key,
        name: name,
        rating: (l: FullLibrary) => {
            if (!l.features[src][key]) { return false}
            return l.features[src][key].level
        },
        notes: (l: FullLibrary) => {
            if (!l.features[src][key]) { return ""}
            return l.features[src][key].notes
        }
    }
}
const checks = [{
    name: "Database Support",
    subChecks: [
        genSubCheck("postgres", "Postgres", "databases"),
        genSubCheck("mysql", "MySQL", "databases"),
        genSubCheck("sqlite", "SQLite", "databases"),
        genSubCheck("other", "Other Databases", "databases")
    ]
}, {
    name: "Mapper Types",
    subChecks: [
        genSubCheck("Query Builder", "Query Builder", "features"),
        genSubCheck("Generic Mapper", "Generic Mapper", "features"),
        genSubCheck("ORM", "ORM", "features"),
        genSubCheck("Query Mapper", "Query Mapper", "features")
    ]
}, {
    name: "Library Ecosystem",
    subChecks: [
        genSubCheck("activity", "Development Activity", "other"),
        genSubCheck("documentation", "Documentation", "other"),
        genSubCheck("examples", "Examples", "other"),
        genSubCheck("migration", "Migrations", "other")
    ]
}]
const featureNames = computed(() => {
    return [
        ...checks.map((c) => {
            return [
                ...c.subChecks.map((s) => {return {grouping: c.name, name: s.name, check: s}})
            ]
        }).flat()
    ]
})
const selectedLibraries = ref(Object.values(store.libraries).map((l: FullLibrary) => l.info.key));

const tableData = computed(() => {
    return featureNames.value.map((feature) => {
        var l = <Record<string, any>>{
            grouping: feature.grouping,
            name: feature.name,
        }
        Object.values(store.libraries).forEach((lib) => {
            l[lib.info.key] = {
                name: lib.info.name,
                grouping: feature.grouping,
                check: feature.check,
                rating: feature.check.rating(lib),
                notes: feature.check.notes ? feature.check.notes(lib) : '',
            }  
        })
        return l
    })
})
const slotInfo = (check: any, lib: FullLibrary) => {
    if (check) {
        return check.rating(lib)
    }
    return ""
}
const getIcon = (rating: string | boolean) => {
    if (!rating) { return naIcon; }
    switch (rating.toString().toLowerCase()) {
        case 'gold': return goldIcon;
        case 'silver': return silverIcon;
        case 'bronze': return bronzeIcon;
        case 'fail': return failIcon;
        default: return naIcon;
    }
}
</script>

<template>
    <div class="normal-container">
        <div class="flex flex-row justify-between mb-4">
            <h1 class="text-4xl">Feature Matrix</h1>
        </div>
        <div class="flex flex-row justify-between">
            <div class="flex flex-col">
                <DataTable :value="tableData" groupRowsBy="grouping" rowGroupMode="subheader" tableStyle="min-width: 50rem">
                    <Column field="name">
                        <template #body="slotProps">
                            {{ slotProps.data.name }}
                        </template>
                        <template #header>
                            <div class="flex flex-row justify-start">
                                <span class="text-2xl">Feature</span>
                            </div>
                        </template>
                    </Column>
                    <Column v-for="(key, index) in selectedLibraries" :key="index" :field="key">
                        <template #body="slotProps">
                            <div class="flex flex-row justify-center">
                                <img class="w-12 h-12" :src="getIcon(slotProps.data[key].rating)" :alt="slotProps.data[key].rating" v-tooltip.top="slotProps.data[key].notes" />
                            </div>
                        </template>
                        <template #header>
                            <div class="flex flex-row justify-center">
                                <span class="text-2xl">{{ key }}</span>
                            </div>
                        </template>
                    </Column>
                    <Column class="hidden"></Column>
                    <template #groupheader="slotProps">
                        <span class="font-bold text-xl">{{ slotProps.data.grouping }}</span>
                    </template>
                </DataTable>
            </div>
        </div>
    </div>
</template>