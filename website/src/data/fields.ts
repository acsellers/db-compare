export interface FullLibrary {
    info: LibraryInfo;
    benchmarks: Benchmark;
    grades: Record<string, Grade>;
    features: Matrix;
}

export interface LibraryInfo {
    key: string;
    name: string;
    short_description: string;
    markdown_desc: string;
    website: string;
    repo: string;
    databases: string[];
    license: string;
    features: string[];
    popularity: number;
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

export interface Grade {
    level: string;
    notes: string;
}


export interface Matrix {
    databases: Record<string, Grade>;
    features: Record<string, Grade>;
    other: Record<string, Grade>;
}

export interface Sample {
    file?: string;
    query?: string;
    sub_examples?: Record<string, Example>;
}


export interface Feature {
    name: string;
    subjects: Subject[];
}

export interface Subject {
    title: string;
    description: string;
    sub_examples: Example[] | null;
}

export interface Example {
    title: string;
    description: string;
    code: string;
}