import { Claim } from '../types';
import { Card } from './Card';

interface TimelinePanelProps {
  claim: Claim;
}

export function TimelinePanel({ claim }: TimelinePanelProps) {
  return (
    <Card className="sticky top-4">
      <h3 className="font-semibold text-sm text-zinc-800 uppercase tracking-wide mb-4">Claim event timeline</h3>
      {claim.events.length === 0 ? (
        <p className="text-sm text-zinc-500 italic">
          No Claim events yet. Items appear here after Intake, upload, assignment, assessment, and Human approval.
        </p>
      ) : (
        <ol className="relative border-l border-zinc-200 ml-2 space-y-4">
          {claim.events.map((claimEvent) => (
            <li key={claimEvent.id} className="ml-4">
              <div className="absolute w-2 h-2 bg-teal-700 rounded-full -left-1.5 mt-1.5 border border-white" />
              <div className="text-xs font-mono text-zinc-400">{claimEvent.at}</div>
              <div className="text-sm font-semibold text-zinc-800">{claimEvent.action}</div>
              <div className="text-xs text-zinc-600">Actor: {claimEvent.actor}</div>
              {claimEvent.detail && <div className="text-xs text-zinc-500 mt-0.5">{claimEvent.detail}</div>}
            </li>
          ))}
        </ol>
      )}
    </Card>
  );
}
