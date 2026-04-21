import { useState, useEffect } from 'preact/hooks';
import { marked } from 'marked';
import { fetchMarketplace, fetchReadme, fetchPackageStats } from '../lib/registry';
import { TYPE_BADGE_CLASSES } from '../lib/constants';
import type { Package, PackageStats } from '../lib/types';

function Skeleton() {
  return (
    <div class="animate-pulse space-y-6">
      <div class="h-8 w-64 bg-gray-800 rounded" />
      <div class="flex gap-4">
        <div class="h-5 w-20 bg-gray-800 rounded" />
        <div class="h-5 w-20 bg-gray-800 rounded" />
        <div class="h-5 w-20 bg-gray-800 rounded" />
      </div>
      <div class="h-10 w-80 bg-gray-800 rounded" />
      <div class="space-y-3">
        <div class="h-4 w-full bg-gray-800 rounded" />
        <div class="h-4 w-5/6 bg-gray-800 rounded" />
        <div class="h-4 w-4/6 bg-gray-800 rounded" />
      </div>
    </div>
  );
}

function StatValue({ label, value }: { label: string; value: string | number }) {
  return (
    <div class="text-center">
      <div class="text-2xl font-bold text-gray-100">{value}</div>
      <div class="text-xs text-gray-500 uppercase tracking-wide">{label}</div>
    </div>
  );
}

export default function PackageDetail({ packageName }: { packageName: string }) {
  const [pkg, setPkg] = useState<Package | null>(null);
  const [readme, setReadme] = useState<string | null>(null);
  const [stats, setStats] = useState<PackageStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([
      fetchMarketplace().then((data) => {
        const found = data.packages.find((p) => p.name === packageName);
        if (!found) throw new Error(`Package "${packageName}" not found`);
        setPkg(found);
        return found;
      }).then((found) => fetchReadme(packageName, found.type).then(setReadme)),
      fetchPackageStats(packageName).then(setStats),
    ])
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [packageName]);

  if (loading) return <Skeleton />;

  if (error || !pkg) {
    return (
      <div class="text-center py-12 text-red-400">
        {error || 'Package not found'}
      </div>
    );
  }

  const readmeHtml = readme ? marked.parse(readme) : null;

  return (
    <div class="space-y-8">
      <a
        href="/skill-registry/"
        class="inline-block text-sm text-gray-400 hover:text-gray-200 transition-colors"
      >
        &larr; Back to packages
      </a>

      <div class="flex items-start gap-3">
        <h1 class="text-3xl font-bold">{pkg.name}</h1>
        <span
          class={`text-sm px-2.5 py-0.5 rounded border mt-1 ${TYPE_BADGE_CLASSES[pkg.type] || ''}`}
        >
          {pkg.type}
        </span>
      </div>

      <p class="text-lg text-gray-400">{pkg.description}</p>

      <div class="flex flex-wrap gap-6 text-sm text-gray-400">
        <span>v{pkg.version}</span>
        <span>by {pkg.author}</span>
        <span>{pkg.license}</span>
      </div>

      <div class="flex flex-wrap gap-2">
        {pkg.tags.map((tag) => (
          <span
            key={tag}
            class="text-xs px-2.5 py-1 rounded-full bg-gray-800 text-gray-400"
          >
            {tag}
          </span>
        ))}
      </div>

      <div class="bg-gray-900 border border-gray-700 rounded-lg px-4 py-3 font-mono text-sm text-gray-300">
        skr install {pkg.name}
      </div>

      <div class="flex gap-8 py-4 border-y border-gray-800">
        <StatValue label="Installs" value={stats?.installs ?? '\u2014'} />
        <StatValue label="Invocations" value={stats?.invocations ?? '\u2014'} />
      </div>

      {pkg.notes && (
        <div class="bg-amber-950/30 border border-amber-900/50 rounded-lg px-4 py-3 text-sm text-amber-200">
          {pkg.notes}
        </div>
      )}

      {readmeHtml && (
        <div
          class="prose prose-invert max-w-none"
          dangerouslySetInnerHTML={{ __html: readmeHtml as string }}
        />
      )}
    </div>
  );
}
