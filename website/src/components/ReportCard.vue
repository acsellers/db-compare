<script setup lang="ts">
// Import images as URLs for vite
import goldIcon from '@/assets/gold.svg';
import silverIcon from '@/assets/silver.svg';
import bronzeIcon from '@/assets/bronze.svg';
import todoIcon from '@/assets/todo.svg';
import failIcon from '@/assets/fail.svg';
import { Subjects } from '@/data/subjects';

const props = defineProps({
    card: {
        type: Object,
        required: true
    }
})

const getIcon = (subject: string) => {
    console.log(subject, props.card.grades[subject])
    let grade = props.card.grades[subject];
    if (!grade) { return todoIcon; }
    if (typeof grade === 'object') {
        grade = grade.level;
    }
    if (!grade) return todoIcon;
    switch (grade.toLowerCase()) {
        case 'gold': return goldIcon;
        case 'silver': return silverIcon;
        case 'bronze': return bronzeIcon;
        case 'fail': return failIcon;
        default: return todoIcon;
    }
}
const getAlt = (subject: string) => {
    if (typeof props.card.grades[subject] === 'object') {
        return props.card.grades[subject].level;
    }
    return "Not Graded";
}
const getNotes = (subject: string) => {
    if (typeof props.card.grades[subject] === 'object') {
        return props.card.grades[subject].notes;
    }
    return "Not Graded";
}
const formatSubject = (subject: string) => {
    return subject.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
}
const gradeSections = Subjects;
</script>

<template>
    <!-- Header -->
    <div class="bg-gray-800 text-white p-4 border-b-4 border-blue-500">
        <div class="flex justify-between items-center">
            <div>
                <h1 class="text-5xl font-bold text-blue-400 uppercase tracking-wide">{{ card.name }}</h1>
            </div>
            <div class="text-right">
                <h2 class="text-2xl font-bold">Report Card</h2>
                <div class="flex gap-4 mt-2 justify-end text-sm">
                    <a :href="card.website" target="_blank" class="hover:text-blue-300 underline">Website</a>
                    <a :href="card.repo" target="_blank" class="hover:text-blue-300 underline">Repo</a>
                </div>
            </div>
        </div>
        <p class="mt-4 text-gray-300 italic">{{ card.description }}</p>
    </div>

    <!-- Content -->
    <div class="">
        <div class="overflow-x-auto">
            <table class="w-full table-auto">
                <thead>
                    <tr class="bg-gray-200 text-gray-600 uppercase text-sm leading-normal">
                        <th class="py-3 px-6 text-left">Subject</th>
                        <th class="py-3 px-6 text-center">Grade</th>
                        <th class="py-3 px-6 text-left">Notes</th>
                    </tr>
                </thead>
                <tbody class="text-gray-600 text-sm font-light">
                    <template v-for="section in gradeSections" :key="section.name">
                        <tr class="border-b border-gray-200 hover:bg-gray-100">
                            <td colspan="3" class="pt-6 pb-1 px-6 text-center whitespace-nowrap">
                                <span class="font-bold text-xl text-gray-600">{{ section.name }}</span>
                            </td>
                        </tr>
                        <tr v-for="subject in section.subjects" :key="subject" class="border-b border-gray-200 hover:bg-gray-100">
                            <td class="py-3 px-3 text-left whitespace-nowrap">
                                <span class="font-bold text-xl text-gray-700">{{ formatSubject(subject) }}</span>
                            </td>
                            <td class="py-3 px-6 text-center">
                                <div class="flex items-center justify-center">
                                    <img :src="getIcon(subject)" :alt="getAlt(subject)" class="w-9 h-9" :title="getAlt(subject)" />
                                </div>
                            </td>
                            <td class="py-3 px-6 text-left">
                                <span v-if="getNotes(subject)" class="italic text-gray-500 text-base">
                                    {{ getNotes(subject) }}
                                </span>
                                <span v-else class="text-gray-300">-</span>
                            </td>
                        </tr>
                    </template>
                </tbody>
            </table>
        </div>
    </div>

    <!-- Footer -->
    <div class="bg-gray-50 p-4 text-center border-t border-gray-200 text-xs text-gray-500 uppercase tracking-widest">
        Go Database Library Breakdown - 
        <RouterLink to="/features">Subject Breakdown</RouterLink>
    </div>
</template>