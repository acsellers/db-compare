<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useDataStore } from '@/stores/data'

const props = defineProps({
    lib: {
        type: String,
        required: true
    },
    db: {
        type: String,
        required: true
    },
    sample: {
        type: String,
        required: true
    },
    displayMode: {
        type: String, // "lines" | "highlights" | "src"
        required: true
    }
})

const store = useDataStore()

interface SampleData {
    go_file: string;
    sql_file: string;
    go_src: string;
    sql_src: string;
    go_lines: string[];
    sql_lines: string[];
    go_highlights: string[];
    sql_highlights: string[];
}

const sampleData = ref<SampleData | null>(null)
const loading = ref(true)

watch(() => [props.lib, props.db, props.sample], async () => {
    loading.value = true
    try {
        const res = await fetch(`/samples/${props.lib}-${props.db}_${props.sample}.json`)
        if (res.ok) {
            sampleData.value = await res.json()
        } else {
            sampleData.value = null
        }
    } catch (e) {
        sampleData.value = null
    } finally {
        loading.value = false
    }
}, { immediate: true })

function getCode(lang: 'go' | 'sql') {
    if (!sampleData.value) return ''

    let codeArray: string[] | string = []

    if (props.displayMode === 'src') {
        return lang === 'go' ? sampleData.value.go_src : sampleData.value.sql_src
    } else if (props.displayMode === 'highlights') {
        codeArray = lang === 'go' ? sampleData.value.go_highlights : sampleData.value.sql_highlights
    } else {
        codeArray = lang === 'go' ? sampleData.value.go_lines : sampleData.value.sql_lines
    }

    if (Array.isArray(codeArray)) {
        return codeArray.join('\n...\n')
    }
    return codeArray
}

const goHtml = computed(() => {
    const code = getCode('go')
    if (!code || !store.highlighterInstance) return ''
    return store.highlighterInstance.codeToHtml(code, {
        lang: 'go',
        theme: 'catppuccin-latte' // Adjust to dark mode if we add a theme switcher
    })
})

const sqlHtml = computed(() => {
    const code = getCode('sql')
    if (!code || !store.highlighterInstance) return ''
    return store.highlighterInstance.codeToHtml(code, {
        lang: 'sql',
        theme: 'catppuccin-latte'
    })
})
</script>

<template>
    <div v-if="loading" class="text-surface-500 py-4 text-center">Loading sample...</div>
    <div v-else-if="!sampleData" class="text-surface-500 py-4 text-center">No sample available for this backend.</div>
    <div v-else class="flex flex-col gap-6">
        <div v-if="sampleData.go_file && goHtml">
            <h4 class="text-sm font-semibold text-surface-600 mb-2">{{ sampleData.go_file }}</h4>
            <div class="overflow-x-auto rounded-xl border border-surface-200">
                <div v-html="goHtml" class="p-4 text-sm font-mono whitespace-pre shiki-container"></div>
            </div>
        </div>

        <div v-if="sampleData.sql_file && sqlHtml">
            <h4 class="text-sm font-semibold text-surface-600 mb-2">{{ sampleData.sql_file }}</h4>
            <div class="overflow-x-auto rounded-xl border border-surface-200">
                <div v-html="sqlHtml" class="p-4 text-sm font-mono whitespace-pre shiki-container"></div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Minor shiki adjustments to ensure it looks good */
:deep(.shiki-container pre) {
    margin: 0;
    background: transparent !important;
}
</style>