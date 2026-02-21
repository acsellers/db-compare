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
    rating: (l: string) => string;
    notes?: (l: string) => string;
}
const checks = [{
    name: "Database Support",
    subChecks: [{
        key: "postgres",
        name: "Postgres",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].databases.includes('postgres')
        }
    }, {
        key: "mysql",
        name: "MySQL",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].databases.includes('mysql')
        }
    }, {
        key: "sqlite",
        name: "SQLite",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].databases.includes('sqlite')
        }
    }, {
        key: "other",
        name: "Other Databases",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].databases.length > 3
        },
        notes: (l: string) => {
            if (!store.reportCards[l]) { return ""}
            return store.reportCards[l].databases.join(', ')
        }
    }]
}, {
    name: "Mapper Types",
    subChecks: [{
        key: "query_builder",
        name: "Query Builder",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('Query Builder')
        }
    }, {
        key: "mapper",
        name: "Arbitrary Mapper",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('Mapper')
        }
    }, {
        key: "orm",
        name: "ORM Querying",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('ORM')
        }
    }, {
        key: "query_mapper",
        name: "Static Query Mapper",
        rating: (l: string) => {
            if (!store.reportCards[l]) { return false}
            return store.reportCards[l].features.includes('Static Query Mapper')
        }
    }]
}]

</script>

<template>
    <h1 class="text-4xl">Feature Matrix</h1>
</template>