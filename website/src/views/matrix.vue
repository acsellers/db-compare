<script setup lang="ts">
import { useDataStore } from '../stores/data.ts'

const store = useDataStore();

interface CheckCategory {
    name: string;
    subChecks: CheckItem;
}
interface CheckItem {
    key: string;
    name: string;
    func: (l: string) => boolean;
    notes?: (l: string) => string;
}
const checks = [{
    name: "Database Support",
    subChecks: [{
        key: "postgres",
        name: "Postgres",
        func: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].databases.includes('postgres')
        }
    }, {
        key: "mysql",
        name: "MySQL",
        func: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].databases.includes('mysql')
        }
    }, {
        key: "sqlite",
        name: "SQLite",
        func: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].databases.includes('sqlite')
        }
    }]
}, {
    name: "Mapper Types",
    subChecks: [{
        key: "query_builder",
        name: "Query Builder",
        func: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('Query Builder')
        }
    }, {
        key: "mapper",
        name: "Arbitrary Mapper",
        func: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('Mapper')
        }
    }, {
        key: "orm",
        name: "ORM Querying",
        func: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('ORM')
        }
    }, {
        key: "query_mapper",
        name: "Static Query Mapper",
        func: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('Static Query Mapper')
        }
    }]
}]

</script>

<template>
    <h1 class="text-4xl">Feature Matrix</h1>
</template>