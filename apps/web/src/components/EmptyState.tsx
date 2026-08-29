import { Card } from './Card';

interface EmptyStateProps {
  onSeed: () => void;
}

export function EmptyState({ onSeed }: EmptyStateProps) {
  return (
    <Card className="p-12 text-center max-w-xl mx-auto my-12">
      <h2 className="text-xl font-semibold mb-2">No Claim loaded</h2>
      <p className="text-zinc-500 text-sm mb-6">
        Seed the synthetic motor Claim to start reviewing operational completeness, SLA clocks, and assessment recommendations.
      </p>
      <button
        onClick={onSeed}
        className="bg-teal-700 hover:bg-teal-800 focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 text-white px-6 py-2.5 rounded-md text-sm font-medium transition cursor-pointer"
      >
        Seed demo Claim
      </button>
    </Card>
  );
}
