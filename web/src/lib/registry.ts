import type { Marketplace, PackageStats, StatsSummary } from './types';

const DEFAULT_REGISTRY_URL = 'https://raw.githubusercontent.com/philippeVV/skill-registry/main';

function getRegistryUrl(): string {
  const url = import.meta.env.PUBLIC_REGISTRY_URL || DEFAULT_REGISTRY_URL;
  return url.replace(/\/$/, '');
}

function getStatsUrl(): string | null {
  const url = import.meta.env.PUBLIC_STATS_URL;
  return url ? url.replace(/\/$/, '') : null;
}

export async function fetchMarketplace(): Promise<Marketplace> {
  const res = await fetch(`${getRegistryUrl()}/marketplace.json`);
  if (!res.ok) throw new Error(`Failed to fetch marketplace: ${res.status}`);
  return res.json();
}

export async function fetchReadme(name: string): Promise<string | null> {
  try {
    const res = await fetch(`${getRegistryUrl()}/packages/${name}/README.md`);
    if (!res.ok) return null;
    return res.text();
  } catch {
    return null;
  }
}

export async function fetchPackageStats(name: string): Promise<PackageStats | null> {
  const statsUrl = getStatsUrl();
  if (!statsUrl) return null;
  try {
    const res = await fetch(`${statsUrl}/packages/${name}`);
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export async function fetchStatsSummary(): Promise<StatsSummary[] | null> {
  const statsUrl = getStatsUrl();
  if (!statsUrl) return null;
  try {
    const res = await fetch(`${statsUrl}/summary`);
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}
