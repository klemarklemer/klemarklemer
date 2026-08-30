import type { TriageItem } from '../api/claims';

interface TriageStripProps {
  items: TriageItem[];
  selectedId: number | null;
  onSelect: (id: number) => void;
}

export function TriageStrip({ items, selectedId, onSelect }: TriageStripProps) {
  if (items.length === 0) return null;

  return (
    <section aria-label="Claims needing attention" className="mb-4">
      <h2 className="sr-only">Claims needing attention</h2>
      <ul className="flex flex-wrap gap-2">
        {items.map((item) => {
          const active = item.id === selectedId;
          return (
            <li key={item.id}>
              <button
                onClick={() => onSelect(item.id)}
                aria-pressed={active}
                className={`flex items-center gap-2 px-3 py-2 rounded-md text-sm border transition focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 ${
                  active
                    ? 'bg-teal-50 border-teal-300 text-teal-900'
                    : 'bg-white border-zinc-200 text-zinc-700 hover:bg-zinc-50'
                }`}
              >
                <span className="font-mono font-semibold">{item.claimNumber}</span>
                <span className="text-xs text-zinc-500">{item.stage}</span>
                {item.atRisk && (
                  <span
                    className="w-2 h-2 rounded-full bg-amber-500"
                    title="SLA at risk"
                    aria-label="SLA at risk"
                  />
                )}
                {item.fraud && (
                  <span className="text-xs font-medium bg-red-100 text-red-700 px-1.5 py-0.5 rounded">
                    Fraud
                  </span>
                )}
              </button>
            </li>
          );
        })}
      </ul>
    </section>
  );
}
