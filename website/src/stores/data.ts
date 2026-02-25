import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import docs from '@/data/docs.json';
import examples from '@/data/examples.json';
import markdownit from 'markdown-it';
import { fromHighlighter } from '@shikijs/markdown-it/core'
import { createHighlighterCore } from 'shiki/core'
import { createOnigurumaEngine } from 'shiki/engine/oniguruma'
import type { FullLibrary, Feature, Subject } from '@/data/fields';

const subjects = [
  {
    name: 'Create, Read, Update',
    subjects: [
      'get_sale',
      'create_sale',
      'customer_update'
    ],
    sub_examples: {
      customer_update: [
        'updates',
        'temporary_table'
      ]
    }
  },
  {
    name: 'Grouping and Joining',
    subjects: [
      'basic_grouping',
      'advanced_grouping',
      'with_queries'
    ],
    sub_examples: {
      basic_grouping: [
        'daily_revenue',
        'customer_sales',
        'daily_sold_items'
      ],
      advanced_grouping: [
        'general_sales_report',
        'weekly_sales_report'
      ]
    }
  },
  {
    name: 'Advanced',
    subjects: [
      'sale_search',
      'bulk_customers',
      'json'
    ]
  }
]
const examplesRecord = examples as Record<string, Subject>

export const useDataStore = defineStore('data', () => {
  const libraries = ref<Record<string, FullLibrary>>(docs)
  const features = ref<Feature[]>(
    subjects.map((subject: any) => {
      return {
        name: subject.name,
        subjects: subject.subjects.map((sub: string) => {
          return examplesRecord[sub]
        })
      }
    })
  )
  const md = markdownit();
  createHighlighterCore({
    themes: [
      import('@shikijs/themes/catppuccin-latte'),
      import('@shikijs/themes/catppuccin-macchiato')
    ],
    langs: [
      import('@shikijs/langs/go'),
    ],
    engine: createOnigurumaEngine(() => import('shiki/wasm'))
  }).then((highlighter: any) => {
    md.use(fromHighlighter(highlighter, {
      themes: {
        light: 'catppuccin-latte',
        dark: 'catppuccin-macchiato'
      }
    }))
  })

  function renderMarkdown(markdown: string) {
    return md.render(markdown);
  }

  return { libraries, features, subjects, renderMarkdown }
})
