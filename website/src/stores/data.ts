import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import docs from '@/data/docs.json';
import markdownit from 'markdown-it';
import { fromHighlighter } from '@shikijs/markdown-it/core'
import { createHighlighterCore } from 'shiki/core'
import { createOnigurumaEngine } from 'shiki/engine/oniguruma'

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

export interface Sample {
  file?: string;
  query?: string;
  sub_examples?: Record<string, Example>;
}

export interface Library {
  name: string;
  markdown_desc: string;
  website: string;
  repo: string;
  description: string;
  databases: string[];
  license: string;
  features: string[];
  popularity: number;
}

export interface ReportCard {
  name: string;
  website: string;
  repo: string;
  description: string;
  databases: string[];
  license: string;
  features: string[];
  popularity: number;
  grades: Record<string, Grade>;
  matrix: Matrix;
}

export interface Grade {
  level: string;
  notes: string;
}
export interface Matrix {
  databases: Record<string, Grade>;
  features: Record<string, Grade>;
  other: Record<string, Grade>;
}

export interface Benchmark {
  runDate: string;
  items: BenchmarkItem[];
}

export interface BenchmarkItem {
  name: string;
  time: number;
  average: number;
  rating: string;
  notes: string;
}

export interface Feature {
  name: string;
  subjects: Subject[];
}

export interface Subject {
  title: string;
  description: string;
  sub_examples: Example[];
}

export interface Example {
  title: string;
  description: string;
  code: string;
}

export const useDataStore = defineStore('data', () => {
  const libraries = ref<Record<string, Library>>(docs.Libraries)
  const reportCards = ref<Record<string, ReportCard>>(docs.ReportCards)
  //const benchmarks = ref<Record<string, Benchmark>>(docs.Benchmarks)
  const features = ref<Feature[]>(
    subjects.map((subject: any) => {
      return {
        name: subject.name,
        subjects: subject.subjects.map((sub: string) => {
          return docs.Examples[sub]
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

  return { libraries, reportCards, features, subjects, renderMarkdown }
})
