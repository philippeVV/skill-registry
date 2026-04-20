import { useState, useEffect } from 'preact/hooks';
import { fetchMarketplace, fetchStatsSummary } from '../lib/registry';
import { TYPE_TEXT_CLASSES } from '../lib/constants';
import type { Package, StatsSummary } from '../lib/types';

interface AuthorEntry {
  name: string;
  packageCount: number;
  totalInstalls: number | null;
  totalInvocations: number | null;
}

function SkeletonRows() {
  return (
    <>
      {Array.from({ length: 5 }).map((_, i) => (
        <tr key={i} class="animate-pulse">
          <td class="py-3 pr-4"><div class="h-4 w-6 bg-gray-800 rounded" /></td>
          <td class="py-3 pr-4"><div class="h-4 w-40 bg-gray-800 rounded" /></td>
          <td class="py-3 pr-4"><div class="h-4 w-16 bg-gray-800 rounded" /></td>
          <td class="py-3"><div class="h-4 w-16 bg-gray-800 rounded" /></td>
        </tr>
      ))}
    </>
  );
}

export default function Leaderboard() {
  const [packages, setPackages] = useState<Package[]>([]);
  const [stats, setStats] = useState<StatsSummary[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([
      fetchMarketplace().then((data) => setPackages(data.packages)),
      fetchStatsSummary().then(setStats),
    ])
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const getStats = (name: string) => stats?.find((s) => s.name === name);

  const rankedPackages = [...packages].sort((a, b) => {
    const sa = getStats(a.name);
    const sb = getStats(b.name);
    if (sa && sb) {
      const scoreA = a.type === 'skill' ? sa.invocations : sa.installs;
      const scoreB = b.type === 'skill' ? sb.invocations : sb.installs;
      return scoreB - scoreA;
    }
    if (sa) return -1;
    if (sb) return 1;
    return a.name.localeCompare(b.name);
  });

  const authors: AuthorEntry[] = (() => {
    const map = new Map<string, AuthorEntry>();
    for (const pkg of packages) {
      const entry = map.get(pkg.author) || {
        name: pkg.author,
        packageCount: 0,
        totalInstalls: null,
        totalInvocations: null,
      };
      entry.packageCount++;
      const s = getStats(pkg.name);
      if (s) {
        entry.totalInstalls = (entry.totalInstalls ?? 0) + s.installs;
        entry.totalInvocations = (entry.totalInvocations ?? 0) + s.invocations;
      }
      map.set(pkg.author, entry);
    }
    return [...map.values()].sort((a, b) => {
      const scoreA = (a.totalInstalls ?? 0) + (a.totalInvocations ?? 0);
      const scoreB = (b.totalInstalls ?? 0) + (b.totalInvocations ?? 0);
      if (scoreA !== scoreB) return scoreB - scoreA;
      return b.packageCount - a.packageCount;
    });
  })();

  if (error) {
    return (
      <div class="text-center py-12 text-red-400">
        Failed to load leaderboard: {error}
      </div>
    );
  }

  return (
    <div class="space-y-12">
      <section>
        <h2 class="text-2xl font-bold mb-6">Top Packages</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-gray-500 border-b border-gray-800">
                <th class="py-2 pr-4 w-12">#</th>
                <th class="py-2 pr-4">Package</th>
                <th class="py-2 pr-4">Type</th>
                <th class="py-2 pr-4 text-right">Installs</th>
                <th class="py-2 text-right">Invocations</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <SkeletonRows />
              ) : (
                rankedPackages.map((pkg, i) => {
                  const s = getStats(pkg.name);
                  return (
                    <tr key={pkg.name} class="border-b border-gray-800/50 hover:bg-gray-900/50">
                      <td class="py-3 pr-4 text-gray-500">{i + 1}</td>
                      <td class="py-3 pr-4">
                        <a
                          href={`/skill-registry/packages/${pkg.name}/`}
                          class="text-gray-100 hover:text-white transition-colors"
                        >
                          {pkg.name}
                        </a>
                      </td>
                      <td class={`py-3 pr-4 ${TYPE_TEXT_CLASSES[pkg.type] || 'text-gray-400'}`}>
                        {pkg.type}
                      </td>
                      <td class="py-3 pr-4 text-right text-gray-300">
                        {s ? s.installs.toLocaleString() : '\u2014'}
                      </td>
                      <td class="py-3 text-right text-gray-300">
                        {s ? s.invocations.toLocaleString() : '\u2014'}
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 class="text-2xl font-bold mb-6">Top Authors</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-gray-500 border-b border-gray-800">
                <th class="py-2 pr-4 w-12">#</th>
                <th class="py-2 pr-4">Author</th>
                <th class="py-2 pr-4 text-right">Packages</th>
                <th class="py-2 text-right">Score</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <SkeletonRows />
              ) : (
                authors.map((author, i) => {
                  const score =
                    author.totalInstalls != null || author.totalInvocations != null
                      ? ((author.totalInstalls ?? 0) + (author.totalInvocations ?? 0)).toLocaleString()
                      : '\u2014';
                  return (
                    <tr key={author.name} class="border-b border-gray-800/50 hover:bg-gray-900/50">
                      <td class="py-3 pr-4 text-gray-500">{i + 1}</td>
                      <td class="py-3 pr-4 text-gray-100">{author.name}</td>
                      <td class="py-3 pr-4 text-right text-gray-300">{author.packageCount}</td>
                      <td class="py-3 text-right text-gray-300">{score}</td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
}
