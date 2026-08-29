import { Warning } from '@phosphor-icons/react/Warning';
import { Claim } from '../types';
import { Card } from './Card';

interface ClaimIdentityProps {
  claim: Claim;
}

export function ClaimIdentity({ claim }: ClaimIdentityProps) {
  return (
    <Card>
      <div className="flex flex-wrap items-center justify-between gap-4 border-b border-zinc-100 pb-4 mb-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="font-mono text-base font-semibold text-zinc-900">{claim.id}</span>
            <span className="text-xs uppercase bg-zinc-100 text-zinc-700 px-2 py-0.5 rounded font-mono font-medium">
              {claim.line}
            </span>
            {claim.closed && (
              <span className="text-xs bg-zinc-900 text-white px-2 py-0.5 rounded font-medium">Closed</span>
            )}
          </div>
          <p className="text-xs text-zinc-500 mt-1">
            Policy ID: <span className="font-mono">{claim.policyId}</span>
          </p>
        </div>

        <div className="flex items-center gap-3">
          <div className="text-right">
            <div className="text-xs text-zinc-500 uppercase font-medium">Current Stage</div>
            <div className="text-sm font-semibold text-zinc-800">{claim.stage}</div>
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between bg-zinc-50 border border-zinc-200/80 rounded-md p-3">
        <div className="flex items-center gap-2">
          <Warning size={20} className={claim.sla.atRisk ? 'text-amber-600' : 'text-teal-700'} weight="fill" />
          <div>
            <span className="text-xs font-medium text-zinc-600 uppercase">Stage SLA Clock</span>
            <div className="text-sm font-medium text-zinc-800">Deadline: {claim.sla.deadline}</div>
          </div>
        </div>
        <div className="text-right">
          <div className="font-mono text-lg font-bold text-zinc-900">
            {claim.sla.remainingMinutes < 0 ? `-${Math.abs(claim.sla.remainingMinutes)}m` : `${claim.sla.remainingMinutes}m`}
          </div>
          {claim.sla.atRisk && (
            <span className="text-xs bg-amber-50 text-amber-700 border border-amber-200 px-2 py-0.5 rounded font-medium inline-block">
              At risk
            </span>
          )}
        </div>
      </div>
    </Card>
  );
}
