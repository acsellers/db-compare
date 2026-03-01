<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
    file: string
    library: string
    db: string
    lines: number[][]
    highlights: number[][]
}>()

const emit = defineEmits<{
    (e: 'change', lines: number[][], highlights: number[][]): void
}>()

const fileContent = ref<string[]>([])

const selectedLines = ref<Set<number>>(new Set())
const highlightedLines = ref<Set<number>>(new Set())

watch(() => [props.file, props.library, props.db], async () => {
    if (!props.file || !props.library || !props.db) {
        fileContent.value = []
        return
    }
    try {
        const res = await fetch(`http://localhost:8080/file?library=${props.library}&db=${props.db}&file=${props.file}`)
        if (res.ok) {
            fileContent.value = await res.json()
        } else {
            fileContent.value = []
        }
    } catch (e) {
        fileContent.value = []
    }
}, { immediate: true })

watch(() => props.lines, (newVal) => {
    const newSet = new Set<number>()
    for (const range of newVal || []) {
        for (let i = Math.min(...range); i <= Math.max(...range); i++) {
            newSet.add(i)
        }
    }
    selectedLines.value = newSet
}, { immediate: true, deep: true })

watch(() => props.highlights, (newVal) => {
    const newSet = new Set<number>()
    for (const range of newVal || []) {
        for (let i = Math.min(...range); i <= Math.max(...range); i++) {
            newSet.add(i)
        }
    }
    highlightedLines.value = newSet
}, { immediate: true, deep: true })


function shrinkToRanges(selectedNumbers: Set<number>): number[][] {
    const sorted = Array.from(selectedNumbers).sort((a, b) => a - b);
    const ranges: number[][] = [];
    if (sorted.length === 0) return ranges;
    let start = sorted[0];
    let end = sorted[0];
    for (let i = 1; i < sorted.length; i++) {
        if (sorted[i] === end + 1) {
            end = sorted[i];
        } else {
            ranges.push([start, end]);
            start = sorted[i];
            end = sorted[i];
        }
    }
    ranges.push([start, end]);
    return ranges;
}

const lastClickedLine = ref<number | null>(null)
const lastClickedHighlight = ref<number | null>(null)

function toggleLineSelection(lineNumber: number, event: MouseEvent) {
    const newSet = new Set(selectedLines.value)
    const isChecking = !newSet.has(lineNumber)

    if (event.shiftKey && lastClickedLine.value !== null) {
        const start = Math.min(lineNumber, lastClickedLine.value)
        const end = Math.max(lineNumber, lastClickedLine.value)
        for (let i = start; i <= end; i++) {
            if (isChecking) newSet.add(i)
            else newSet.delete(i)
        }
    } else {
        if (isChecking) newSet.add(lineNumber)
        else newSet.delete(lineNumber)
    }

    lastClickedLine.value = lineNumber
    selectedLines.value = newSet
    emitChange()
}

function toggleHighlightSelection(lineNumber: number, event: MouseEvent) {
    const newSet = new Set(highlightedLines.value)
    const isChecking = !newSet.has(lineNumber)

    if (event.shiftKey && lastClickedHighlight.value !== null) {
        const start = Math.min(lineNumber, lastClickedHighlight.value)
        const end = Math.max(lineNumber, lastClickedHighlight.value)
        for (let i = start; i <= end; i++) {
            if (isChecking) newSet.add(i)
            else newSet.delete(i)
        }
    } else {
        if (isChecking) newSet.add(lineNumber)
        else newSet.delete(lineNumber)
    }

    lastClickedHighlight.value = lineNumber
    highlightedLines.value = newSet
    emitChange()
}

function emitChange() {
    const lRange = shrinkToRanges(selectedLines.value)
    const hRange = shrinkToRanges(highlightedLines.value)
    emit('change', lRange, hRange)
}

function formatRanges(ranges: number[][]): string {
    if (!ranges || ranges.length === 0) return 'None'
    return ranges.map(r => r[0] === r[1] ? `${r[0]}` : `${r[0]}-${r[1]}`).join(', ')
}
</script>

<template>
    <div v-if="!file" class="p-8 text-center bg-surface-100 dark:bg-surface-800 rounded-xl mt-4">
        <p class="text-surface-600 dark:text-surface-300">Please select a file to view and edit its lines.</p>
    </div>

    <div v-else class="flex flex-col gap-4 mt-4">
        <div
            class="flex flex-col md:flex-row gap-4 bg-surface-100 dark:bg-surface-800 p-4 rounded-xl border border-surface-200 dark:border-surface-700">
            <div class="flex-1">
                <h3 class="font-bold text-surface-900 dark:text-surface-50 mb-1 text-sm uppercase tracking-wide">
                    Included Lines</h3>
                <div class="text-surface-700 dark:text-surface-300 font-mono text-sm">
                    {{ formatRanges(props.lines) }}
                </div>
            </div>
            <div class="w-px bg-surface-200 dark:bg-surface-700 hidden md:block"></div>
            <div class="flex-1">
                <h3 class="font-bold text-surface-900 dark:text-surface-50 mb-1 text-sm uppercase tracking-wide">
                    Highlighted Lines</h3>
                <div class="text-surface-700 dark:text-surface-300 font-mono text-sm">
                    {{ formatRanges(props.highlights) }}
                </div>
            </div>
        </div>

        <div
            class="overflow-auto border border-surface-200 dark:border-surface-700 rounded-xl font-mono text-sm max-h-[60vh] bg-surface-0 dark:bg-surface-950 shadow-sm relative">
            <div
                class="sticky top-0 z-10 hidden sm:flex bg-surface-50 dark:bg-surface-900 border-b border-surface-200 dark:border-surface-700 font-bold text-xs uppercase text-surface-500 py-2">
                <div class="w-12 text-center flex-shrink-0">#</div>
                <div class="w-12 text-center flex-shrink-0" title="Include/Crop Line">Line</div>
                <div class="w-12 text-center flex-shrink-0" title="Highlight Line">High</div>
                <div class="px-3 flex-1 flex items-center">Code</div>
            </div>

            <div v-for="(line, idx) in fileContent" :key="idx"
                class="flex flex-row hover:bg-surface-100 dark:hover:bg-surface-800 py-1 transition-colors" :class="{
                    'bg-blue-50/50 dark:bg-blue-900/20': selectedLines.has(idx + 1) && !highlightedLines.has(idx + 1),
                    'bg-yellow-50/50 dark:bg-yellow-900/20': highlightedLines.has(idx + 1) && !selectedLines.has(idx + 1),
                    'bg-green-50/50 dark:bg-green-900/20': selectedLines.has(idx + 1) && highlightedLines.has(idx + 1)
                }">
                <div
                    class="w-12 text-right pr-2 text-surface-400 select-none border-r border-surface-200 dark:border-surface-700 flex-shrink-0 pt-0.5">
                    {{ idx + 1 }}
                </div>
                <div class="w-12 flex justify-center items-center border-r border-surface-200 dark:border-surface-700 flex-shrink-0 cursor-pointer"
                    title="Toggle Line Selection" @click="toggleLineSelection(idx + 1, $event)">
                    <input type="checkbox" :checked="selectedLines.has(idx + 1)"
                        @click.stop.prevent="toggleLineSelection(idx + 1, $event)" class="cursor-pointer" />
                </div>
                <div class="w-12 flex justify-center items-center border-r border-surface-200 dark:border-surface-700 flex-shrink-0 cursor-pointer"
                    title="Toggle Highlight Selection" @click="toggleHighlightSelection(idx + 1, $event)">
                    <input type="checkbox" :checked="highlightedLines.has(idx + 1)"
                        @click.stop.prevent="toggleHighlightSelection(idx + 1, $event)" class="cursor-pointer" />
                </div>
                <div class="px-4 whitespace-pre overflow-x-auto text-surface-800 dark:text-surface-200 flex-1">
                    {{ line === '' ? ' ' : line }}
                </div>
            </div>

            <div v-if="fileContent.length === 0" class="p-8 text-center text-surface-500">
                Loading ... or file is empty.
            </div>
        </div>
    </div>
</template>