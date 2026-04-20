import { useState, useEffect } from 'preact/hooks';
import { fetchMarketplace } from '../lib/registry';
import { TYPE_BADGE_CLASSES } from '../lib/constants';
import type { Package } from '../lib/types';

function PackageCard({ pkg }: { pkg: Package }) {
  return (
    <a
      href={`/skill-registry/packages/${pkg.name}/`}
      class="block bg-gray-900 border border-gray-800 rounded-lg p-5 hover:border-gray-600 transition-colors"
    >
      <div class="flex items-start justify-between gap-2 mb-2">
        <h3 class="font-semibold text-gray-100 truncate">{pkg.name}</h3>
        <span
          class={`text-xs px-2 py-0.5 rounded border shrink-0 ${TYPE_BADGE_CLASSES[pkg.type] || ''}`}
        >
          {pkg.type}
        </span>
      </div>
      <p class="text-sm text-gray-400 mb-3 line-clamp-2">{pkg.description}</p>
      <div class="flex flex-wrap gap-1.5 mb-3">
        {pkg.tags.map((tag) => (
          <span
            key={tag}
            class="text-xs px-2 py-0.5 rounded-full bg-gray-800 text-gray-400"
          >
            {tag}
          </span>
        ))}
      </div>
      <div class="text-xs text-gray-500">
        &darr; &mdash; installs &middot; &lightning; &mdash; invocations
      </div>
    </a>
  );
}

function SkeletonCard() {
  return (
    <div class="bg-gray-900 border border-gray-800 rounded-lg p-5 animate-pulse">
      <div class="flex justify-between mb-2">
        <div class="h-5 w-32 bg-gray-800 rounded" />
        <div class="h-5 w-12 bg-gray-800 rounded" />
      </div>
      <div class="h-4 w-full bg-gray-800 rounded mb-2" />
      <div class="h-4 w-2/3 bg-gray-800 rounded mb-3" />
      <div class="flex gap-1.5 mb-3">
        <div class="h-5 w-16 bg-gray-800 rounded-full" />
        <div class="h-5 w-14 bg-gray-800 rounded-full" />
      </div>
      <div class="h-3 w-40 bg-gray-800 rounded" />
    </div>
  );
}

export default function PackageBrowser() {
  const [packages, setPackages] = useState<Package[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchText, setSearchText] = useState('');
  const [activeTags, setActiveTags] = useState<Set<string>>(new Set());

  useEffect(() => {
    fetchMarketplace()
      .then((data) => setPackages(data.packages))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const allTags = [...new Set(packages.flatMap((p) => p.tags))].sort();

  const toggleTag = (tag: string) => {
    setActiveTags((prev) => {
      const next = new Set(prev);
      if (next.has(tag)) next.delete(tag);
      else next.add(tag);
      return next;
    });
  };

  const filtered = packages.filter((pkg) => {
    const search = searchText.toLowerCase();
    const matchesSearch =
      !search ||
      pkg.name.toLowerCase().includes(search) ||
      pkg.description.toLowerCase().includes(search) ||
      pkg.tags.some((t) => t.toLowerCase().includes(search));
    const matchesTags =
      activeTags.size === 0 || [...activeTags].every((t) => pkg.tags.includes(t));
    return matchesSearch && matchesTags;
  });

  if (error) {
    return (
      <div class="text-center py-12 text-red-400">
        Failed to load packages: {error}
      </div>
    );
  }

  return (
    <div>
      {!loading && (
        <div class="mb-6 space-y-4">
          <input
            type="text"
            placeholder="Search packages..."
            value={searchText}
            onInput={(e) => setSearchText((e.target as HTMLInputElement).value)}
            class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-gray-100 placeholder-gray-500 focus:outline-none focus:border-gray-500 transition-colors"
          />
          <div class="flex flex-wrap gap-2">
            {allTags.map((tag) => (
              <button
                key={tag}
                onClick={() => toggleTag(tag)}
                class={`text-xs px-3 py-1 rounded-full border transition-colors cursor-pointer ${
                  activeTags.has(tag)
                    ? 'bg-blue-900/50 text-blue-300 border-blue-700'
                    : 'bg-gray-900 text-gray-400 border-gray-700 hover:border-gray-500'
                }`}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>
      )}

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {loading
          ? Array.from({ length: 6 }).map((_, i) => <SkeletonCard key={i} />)
          : filtered.map((pkg) => <PackageCard key={pkg.name} pkg={pkg} />)}
      </div>

      {!loading && filtered.length === 0 && (
        <div class="text-center py-8 text-gray-500">
          No packages match your search.
        </div>
      )}
    </div>
  );
}
