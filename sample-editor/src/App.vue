<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue"
import Select from 'primevue/select';
import Card from 'primevue/card';
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import Button from 'primevue/button';
import LineSelector from './components/LineSelector.vue';

const examplesData = ref({
  "get_sale": [],
  "create_sale": [],
  "customer_update": ["updates", "temporary_table"],
  "basic_grouping": ["daily_revenue", "customer_sales", "daily_sold_items"],
  "advanced_grouping": ["general_sales_report", "weekly_sales_report"],
  "with_queries": ["general_sales_report", "weekly_sales_report"],
  "sale_search": [],
  "bulk_update": [],
  "json": [],
})

const fileListing = ref<any[]>([])

const selectedLibrary = ref('')
const selectedDb = ref('mysql')
const selectedExample = ref('')
const selectedSubExample = ref('')

interface SampleFile {
  key: string;
  library: string;
  name: string;
  store: {
    mysql: SampleList;
    postgres: SampleList;
    sqlite: SampleList;
  }
}
interface SampleList {
  get_sale: SampleEntry;
  create_sale: SampleEntry;
  customer_update: {
    updates: SampleEntry;
    temporary_table: SampleEntry;
  };
  basic_grouping: {
    daily_revenue: SampleEntry;
    customer_sales: SampleEntry;
    daily_sold_items: SampleEntry;
  };
  advanced_grouping: {
    general_sales_report: SampleEntry;
    weekly_sales_report: SampleEntry;
  };
  with_queries: {
    general_sales_report: SampleEntry;
    weekly_sales_report: SampleEntry;
  };
  sale_search: SampleEntry;
  bulk_update: SampleEntry;
  json: SampleEntry;
}
interface SampleEntry {
  go: SampleFile;
  sql: SampleFile;
}
interface SampleFile {
  file: string;
  lines: number[][];
  highlights: number[][];
}


const dbs = ['mysql', 'postgres', 'sqlite']
const libraries = computed(() => {
  return fileListing.value.map(f => f.library) || []
})
const examples = computed(() => {
  return Object.keys(examplesData.value)
})
const subExamples = computed(() => {
  if (!selectedExample.value) return []
  return (examplesData.value as any)[selectedExample.value] || []
})

const filesForCurrentLibDb = computed(() => {
  if (!selectedLibrary.value || !selectedDb.value) return []
  const libData = fileListing.value.find(f => f.library === selectedLibrary.value)
  if (!libData || !libData.samples_by_db) return []
  return libData.samples_by_db[selectedDb.value] || []
})

const goFiles = computed(() => {
  return filesForCurrentLibDb.value.filter(f => f.endsWith('.go'))
})

const sqlFiles = computed(() => {
  return filesForCurrentLibDb.value.filter(f => f.endsWith('.sql'))
})

onMounted(async () => {
  const res = await fetch('http://localhost:8080/filelist')
  fileListing.value = await res.json()
  if (libraries.value.length > 0) {
    selectedLibrary.value = libraries.value[0]
  }
})

watch(selectedExample, () => {
  if (subExamples.value.length > 0) {
    selectedSubExample.value = subExamples.value[0]
  } else {
    selectedSubExample.value = ''
  }
  loadSelectedFiles()
})

watch(selectedSubExample, () => {
  loadSelectedFiles()
})


const goCode = ref('')
const sqlCode = ref('')
const goLines = ref<number[][]>([])
const sqlLines = ref<number[][]>([])
const goHighlights = ref<number[][]>([])
const sqlHighlights = ref<number[][]>([])

const samples = ref<SampleFile | null>(null)
async function updateLibrarySamples() {
  console.log('updateLibrarySamples', selectedLibrary.value)
  if (!selectedLibrary.value) return
  const res = await fetch(`http://localhost:8080/samples?library=${selectedLibrary.value}`)
  console.log('res', res)
  samples.value = await res.json()
  console.log('samples', samples.value)
}
watch(selectedLibrary, updateLibrarySamples, { immediate: true })

function loadSelectedFiles() {
  if (!samples.value) return
  let dbData = samples.value.store.mysql
  if (selectedDb.value === 'postgres') {
    dbData = samples.value.store.postgres
  } else if (selectedDb.value === 'sqlite') {
    dbData = samples.value.store.sqlite
  }
  if (!dbData || !selectedExample.value) return
  try {
    let entry: SampleEntry | null = null
    switch (selectedExample.value) {
      case 'get_sale':
        entry = dbData.get_sale
        break
      case 'create_sale':
        entry = dbData.create_sale
        break
      case 'customer_update':
        if (selectedSubExample.value === 'updates') entry = dbData.customer_update.updates
        else if (selectedSubExample.value === 'temporary_table') entry = dbData.customer_update.temporary_table
        break
      case 'basic_grouping':
        if (selectedSubExample.value === 'daily_revenue') entry = dbData.basic_grouping.daily_revenue
        else if (selectedSubExample.value === 'customer_sales') entry = dbData.basic_grouping.customer_sales
        else if (selectedSubExample.value === 'daily_sold_items') entry = dbData.basic_grouping.daily_sold_items
        break
      case 'advanced_grouping':
        if (selectedSubExample.value === 'general_sales_report') entry = dbData.advanced_grouping.general_sales_report
        else if (selectedSubExample.value === 'weekly_sales_report') entry = dbData.advanced_grouping.weekly_sales_report
        break
      case 'with_queries':
        if (selectedSubExample.value === 'general_sales_report') entry = dbData.with_queries.general_sales_report
        else if (selectedSubExample.value === 'weekly_sales_report') entry = dbData.with_queries.weekly_sales_report
        break
      case 'sale_search':
        entry = dbData.sale_search
        break
      case 'bulk_update':
        entry = dbData.bulk_update
        break
      case 'json':
        entry = dbData.json
        break
    }
    if (entry) {
      goCode.value = entry.go?.file || ''
      goLines.value = entry.go?.lines || []
      goHighlights.value = entry.go?.highlights || []
      sqlCode.value = entry.sql?.file || ''
      sqlLines.value = entry.sql?.lines || []
      sqlHighlights.value = entry.sql?.highlights || []
    } else {
      goCode.value = ''
      goLines.value = []
      goHighlights.value = []
      sqlCode.value = ''
      sqlLines.value = []
      sqlHighlights.value = []
    }
  } catch (e) {
    goCode.value = ''
    sqlCode.value = ''
    goLines.value = []
    sqlLines.value = []
    goHighlights.value = []
    sqlHighlights.value = []
  }
}

function updateSelectedFiles() {
  if (!samples.value) {
    goCode.value = ''
    sqlCode.value = ''
    return
  }

  let dbData = samples.value.store.mysql;
  if (selectedDb.value === 'postgres') {
    dbData = samples.value.store.postgres;
  } else if (selectedDb.value === 'sqlite') {
    dbData = samples.value.store.sqlite;
  }

  if (!dbData || !selectedExample.value) {
    goCode.value = ''
    sqlCode.value = ''
    return
  }

  try {
    let entry: SampleEntry | null = null;
    switch (selectedExample.value) {
      case 'get_sale':
        entry = dbData.get_sale;
        break;
      case 'create_sale':
        entry = dbData.create_sale;
        break;
      case 'customer_update':
        if (selectedSubExample.value === 'updates') entry = dbData.customer_update.updates;
        else if (selectedSubExample.value === 'temporary_table') entry = dbData.customer_update.temporary_table;
        break;
      case 'basic_grouping':
        if (selectedSubExample.value === 'daily_revenue') entry = dbData.basic_grouping.daily_revenue;
        else if (selectedSubExample.value === 'customer_sales') entry = dbData.basic_grouping.customer_sales;
        else if (selectedSubExample.value === 'daily_sold_items') entry = dbData.basic_grouping.daily_sold_items;
        break;
      case 'advanced_grouping':
        if (selectedSubExample.value === 'general_sales_report') entry = dbData.advanced_grouping.general_sales_report;
        else if (selectedSubExample.value === 'weekly_sales_report') entry = dbData.advanced_grouping.weekly_sales_report;
        break;
      case 'with_queries':
        if (selectedSubExample.value === 'general_sales_report') entry = dbData.with_queries.general_sales_report;
        else if (selectedSubExample.value === 'weekly_sales_report') entry = dbData.with_queries.weekly_sales_report;
        break;
      case 'sale_search':
        entry = dbData.sale_search;
        break;
      case 'bulk_update':
        entry = dbData.bulk_update;
        break;
      case 'json':
        entry = dbData.json;
        break;
    }

    if (entry) {
      goCode.value = entry.go?.file || ''
      sqlCode.value = entry.sql?.file || ''
    } else {
      goCode.value = ''
      sqlCode.value = ''
    }
  } catch (e) {
    goCode.value = ''
    sqlCode.value = ''
  }
}


function handleGoChange(lines: number[][], highlights: number[][]) {
  console.log('handleGoChange', lines, highlights)
  // set the file, lines, and highlights for the samples record for the selected example/subexample
  if (!samples.value) return
  let dbData = samples.value.store.mysql;
  if (selectedDb.value === 'postgres') {
    dbData = samples.value.store.postgres;
  } else if (selectedDb.value === 'sqlite') {
    dbData = samples.value.store.sqlite;
  }
  goLines.value = lines
  goHighlights.value = highlights
  switch (selectedExample.value) {
    case 'get_sale':
      dbData.get_sale.go.file = goCode.value
      dbData.get_sale.go.lines = lines
      dbData.get_sale.go.highlights = highlights
      break
    case 'create_sale':
      dbData.create_sale.go.file = goCode.value
      dbData.create_sale.go.lines = lines
      dbData.create_sale.go.highlights = highlights
      break
    case 'customer_update':
      switch (selectedSubExample.value) {
        case 'updates':
          dbData.customer_update.updates.go.file = goCode.value
          dbData.customer_update.updates.go.lines = lines
          dbData.customer_update.updates.go.highlights = highlights
          break
        case 'temporary_table':
          dbData.customer_update.temporary_table.go.file = goCode.value
          dbData.customer_update.temporary_table.go.lines = lines
          dbData.customer_update.temporary_table.go.highlights = highlights
          break
      }
      break
    case 'basic_grouping':
      switch (selectedSubExample.value) {
        case 'daily_revenue':
          dbData.basic_grouping.daily_revenue.go.file = goCode.value
          dbData.basic_grouping.daily_revenue.go.lines = lines
          dbData.basic_grouping.daily_revenue.go.highlights = highlights
          break
        case 'customer_sales':
          dbData.basic_grouping.customer_sales.go.file = goCode.value
          dbData.basic_grouping.customer_sales.go.lines = lines
          dbData.basic_grouping.customer_sales.go.highlights = highlights
          break
        case 'daily_sold_items':
          dbData.basic_grouping.daily_sold_items.go.file = goCode.value
          dbData.basic_grouping.daily_sold_items.go.lines = lines
          dbData.basic_grouping.daily_sold_items.go.highlights = highlights
          break
      }
      break
    case 'advanced_grouping':
      switch (selectedSubExample.value) {
        case 'general_sales_report':
          dbData.advanced_grouping.general_sales_report.go.file = goCode.value
          dbData.advanced_grouping.general_sales_report.go.lines = lines
          dbData.advanced_grouping.general_sales_report.go.highlights = highlights
          break
        case 'weekly_sales_report':
          dbData.advanced_grouping.weekly_sales_report.go.file = goCode.value
          dbData.advanced_grouping.weekly_sales_report.go.lines = lines
          dbData.advanced_grouping.weekly_sales_report.go.highlights = highlights
          break
      }
      break
    case 'with_queries':
      switch (selectedSubExample.value) {
        case 'general_sales_report':
          dbData.with_queries.general_sales_report.go.file = goCode.value
          dbData.with_queries.general_sales_report.go.lines = lines
          dbData.with_queries.general_sales_report.go.highlights = highlights
          break
        case 'weekly_sales_report':
          dbData.with_queries.weekly_sales_report.go.file = goCode.value
          dbData.with_queries.weekly_sales_report.go.lines = lines
          dbData.with_queries.weekly_sales_report.go.highlights = highlights
          break
      }
      break
    case 'sale_search':
      dbData.sale_search.go.file = goCode.value
      dbData.sale_search.go.lines = lines
      dbData.sale_search.go.highlights = highlights
      break
    case 'bulk_update':
      dbData.bulk_update.go.file = goCode.value
      dbData.bulk_update.go.lines = lines
      dbData.bulk_update.go.highlights = highlights
      break
    case 'json':
      dbData.json.go.file = goCode.value
      dbData.json.go.lines = lines
      dbData.json.go.highlights = highlights
      break
  }
}
function handleSqlChange(lines: number[][], highlights: number[][]) {
  if (!samples.value) return
  let dbData = samples.value.store.mysql;
  if (selectedDb.value === 'postgres') {
    dbData = samples.value.store.postgres;
  } else if (selectedDb.value === 'sqlite') {
    dbData = samples.value.store.sqlite;
  }
  switch (selectedExample.value) {
    case 'get_sale':
      dbData.get_sale.sql.file = sqlCode.value
      dbData.get_sale.sql.lines = lines
      dbData.get_sale.sql.highlights = highlights
      break
    case 'create_sale':
      dbData.create_sale.sql.file = sqlCode.value
      dbData.create_sale.sql.lines = lines
      dbData.create_sale.sql.highlights = highlights
      break
    case 'customer_update':
      switch (selectedSubExample.value) {
        case 'updates':
          dbData.customer_update.updates.sql.file = sqlCode.value
          dbData.customer_update.updates.sql.lines = lines
          dbData.customer_update.updates.sql.highlights = highlights
          break
        case 'temporary_table':
          dbData.customer_update.temporary_table.sql.file = sqlCode.value
          dbData.customer_update.temporary_table.sql.lines = lines
          dbData.customer_update.temporary_table.sql.highlights = highlights
          break
      }
      break
    case 'basic_grouping':
      switch (selectedSubExample.value) {
        case 'daily_revenue':
          dbData.basic_grouping.daily_revenue.sql.file = sqlCode.value
          dbData.basic_grouping.daily_revenue.sql.lines = lines
          dbData.basic_grouping.daily_revenue.sql.highlights = highlights
          break
        case 'customer_sales':
          dbData.basic_grouping.customer_sales.sql.file = sqlCode.value
          dbData.basic_grouping.customer_sales.sql.lines = lines
          dbData.basic_grouping.customer_sales.sql.highlights = highlights
          break
        case 'daily_sold_items':
          dbData.basic_grouping.daily_sold_items.sql.file = sqlCode.value
          dbData.basic_grouping.daily_sold_items.sql.lines = lines
          dbData.basic_grouping.daily_sold_items.sql.highlights = highlights
          break
      }
      break
    case 'advanced_grouping':
      switch (selectedSubExample.value) {
        case 'general_sales_report':
          dbData.advanced_grouping.general_sales_report.sql.file = sqlCode.value
          dbData.advanced_grouping.general_sales_report.sql.lines = lines
          dbData.advanced_grouping.general_sales_report.sql.highlights = highlights
          break
        case 'weekly_sales_report':
          dbData.advanced_grouping.weekly_sales_report.sql.file = sqlCode.value
          dbData.advanced_grouping.weekly_sales_report.sql.lines = lines
          dbData.advanced_grouping.weekly_sales_report.sql.highlights = highlights
          break
      }
      break
    case 'with_queries':
      switch (selectedSubExample.value) {
        case 'general_sales_report':
          dbData.with_queries.general_sales_report.sql.file = sqlCode.value
          dbData.with_queries.general_sales_report.sql.lines = lines
          dbData.with_queries.general_sales_report.sql.highlights = highlights
          break
        case 'weekly_sales_report':
          dbData.with_queries.weekly_sales_report.sql.file = sqlCode.value
          dbData.with_queries.weekly_sales_report.sql.lines = lines
          dbData.with_queries.weekly_sales_report.sql.highlights = highlights
          break
      }
      break
    case 'sale_search':
      dbData.sale_search.sql.file = sqlCode.value
      dbData.sale_search.sql.lines = lines
      dbData.sale_search.sql.highlights = highlights
      break
    case 'bulk_update':
      dbData.bulk_update.sql.file = sqlCode.value
      dbData.bulk_update.sql.lines = lines
      dbData.bulk_update.sql.highlights = highlights
      break
    case 'json':
      dbData.json.sql.file = sqlCode.value
      dbData.json.sql.lines = lines
      dbData.json.sql.highlights = highlights
      break
  }
}
async function saveSamples() {
  if (!selectedLibrary.value || !samples.value) return

  try {
    const res = await fetch(`http://localhost:8080/save?library=${selectedLibrary.value}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(samples.value)
    })

    if (res.ok) {
      alert('Saved successfully!')
    } else {
      alert('Failed to save')
    }
  } catch (e) {
    alert('Failed to save')
    console.error(e)
  }
}
</script>

<template>
  <div class="min-h-screen bg-surface-50 dark:bg-surface-950 p-8">
    <div class="max-w-6xl mx-auto space-y-6">
      <div class="flex flex-row justify-between items-center">
        <h1 class="text-3xl font-bold text-surface-900 dark:text-surface-0">Sample Editor</h1>
        <Button label="Save Changes" icon="pi pi-save" @click="saveSamples" :disabled="!samples" severity="success" />
      </div>

      <Card>
        <template #content>
          <div class="">
            <div class="flex flex-row gap-2">
              <div class="flex flex-col gap-2">
                <label for="library" class="font-medium text-surface-700 dark:text-surface-300">Library</label>
                <Select id="library" v-model="selectedLibrary" :options="libraries" placeholder="Select Library"
                  class="w-full" @change="updateLibrarySamples" />
              </div>
              <div class="flex flex-col gap-2">
                <label for="db" class="font-medium text-surface-700 dark:text-surface-300">Database</label>
                <Select id="db" v-model="selectedDb" :options="dbs" placeholder="Select Database" class="w-full" />
              </div>
              <div class="flex flex-col gap-2">
                <label for="example" class="font-medium text-surface-700 dark:text-surface-300">Example</label>
                <Select id="example" v-model="selectedExample" :options="examples" placeholder="Select Example"
                  class="w-full" />
              </div>
              <div class="flex flex-col gap-2" v-if="subExamples.length > 0">
                <label for="subexample" class="font-medium text-surface-700 dark:text-surface-300">Sub-Example</label>
                <Select id="subexample" v-model="selectedSubExample" :options="subExamples"
                  placeholder="Select Sub-Example" class="w-full" />
              </div>
            </div>
          </div>
        </template>
      </Card>

      <div v-if="selectedLibrary && selectedDb && selectedExample" class="">
        <Tabs value="go">
          <TabList>
            <Tab value="go">Go</Tab>
            <Tab value="sql">SQL</Tab>
          </TabList>
          <TabPanels>
            <TabPanel value="go">
              <div class="flex flex-col gap-2">
                <label for="go" class="font-medium text-surface-700 dark:text-surface-300">Go Files</label>
                <Select id="go" v-model="goCode" :options="goFiles" placeholder="Select Go File" class="w-full" />
              </div>
              <LineSelector @change="handleGoChange" :file="goCode" :library="selectedLibrary" :db="selectedDb"
                :lines="goLines" :highlights="goHighlights" />
            </TabPanel>
            <TabPanel value="sql">
              <div class="flex flex-col gap-2">
                <label for="sql" class="font-medium text-surface-700 dark:text-surface-300">SQL Files</label>
                <Select id="sql" v-model="sqlCode" :options="sqlFiles" placeholder="Select SQL File" class="w-full" />
              </div>
              <LineSelector @change="handleSqlChange" :file="sqlCode" :library="selectedLibrary" :db="selectedDb"
                :lines="sqlLines" :highlights="sqlHighlights" />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </div>
    </div>
  </div>
</template>
