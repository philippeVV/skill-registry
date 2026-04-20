export interface Package {
  name: string;
  type: 'skill' | 'rule' | 'knowledge';
  version: string;
  description: string;
  tags: string[];
  author: string;
  license: string;
  path: string;
  files: string[];
  artifact_hash: string;
  notes?: string;
}

export interface Marketplace {
  packages: Package[];
}

export interface PackageStats {
  installs: number;
  invocations: number;
}

export interface StatsSummary {
  name: string;
  installs: number;
  invocations: number;
}
