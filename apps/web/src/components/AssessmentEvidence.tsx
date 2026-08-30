import { Claim } from '../types';
import { Card } from './Card';

interface AssessmentEvidenceProps {
  claim: Claim;
}

function parseExtracted(raw?: string): Record<string, string> {
  if (!raw) return {};
  try {
    const value = JSON.parse(raw);
    if (value && typeof value === 'object') {
      return Object.fromEntries(
        Object.entries(value as Record<string, unknown>).map(([key, val]) => [key, String(val)]),
      );
    }
  } catch {
    return {};
  }
  return {};
}

const LABELS: Record<string, string> = {
  detected_damage: 'Detected damage',
  severity: 'Severity',
  report_id: 'Report ID',
  officer: 'Reporting officer',
  incident_date: 'Incident date',
  location: 'Location',
  liability: 'Liability',
};

export function AssessmentEvidence({ claim }: AssessmentEvidenceProps) {
  const extracted = claim.documents
    .filter((doc) => doc.present && doc.extractedData)
    .map((doc) => ({ name: doc.name, fields: parseExtracted(doc.extractedData) }));

  return (
    <Card className="bg-zinc-50/60">
      <details className="group">
        <summary className="cursor-pointer select-none text-xs font-semibold uppercase tracking-wide text-zinc-600 flex items-center justify-between">
          Assessment evidence
          <span className="text-zinc-400 group-open:hidden">Show</span>
        </summary>

        <div className="mt-4 space-y-4">
          <section>
            <h4 className="text-xs font-medium text-zinc-700 mb-2">Document analysis</h4>
            {extracted.length === 0 ? (
              <p className="text-xs text-zinc-500 italic">No extracted document data yet.</p>
            ) : (
              <div className="space-y-3">
                {extracted.map((doc) => (
                  <div key={doc.name}>
                    <div className="text-xs font-mono text-zinc-500 mb-1">{doc.name}</div>
                    <dl className="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-xs">
                      {Object.entries(doc.fields).map(([key, value]) => (
                        <EvidenceField key={key} label={key} value={value} />
                      ))}
                    </dl>
                  </div>
                ))}
              </div>
            )}
          </section>

          {claim.policy && (
            <section>
              <h4 className="text-xs font-medium text-zinc-700 mb-2">Policy check</h4>
              <dl className="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-xs">
                <dt className="text-zinc-500">Policy</dt>
                <dd className="font-mono text-zinc-800">{claim.policy.policyNumber}</dd>
                <dt className="text-zinc-500">Holder</dt>
                <dd className="text-zinc-800">{claim.policy.policyHolder}</dd>
                <dt className="text-zinc-500">Vehicle</dt>
                <dd className="text-zinc-800">{claim.policy.vehicle}</dd>
                <dt className="text-zinc-500">Status</dt>
                <dd className="text-zinc-800">
                  {claim.policy.status === 'ACTIVE' ? 'In force' : claim.policy.status}
                </dd>
                <dt className="text-zinc-500">Coverage</dt>
                <dd className="font-mono text-zinc-800">
                  {claim.policy.coverage} up to ${claim.policy.maxCoverage.toLocaleString()}
                </dd>
                <dt className="text-zinc-500">Deductible</dt>
                <dd className="font-mono text-zinc-800">${claim.policy.deductible.toLocaleString()}</dd>
              </dl>
            </section>
          )}
        </div>
      </details>
    </Card>
  );
}

function EvidenceField({ label, value }: { label: string; value: string }) {
  return (
    <>
      <dt className="text-zinc-500">{LABELS[label] ?? label}</dt>
      <dd className="text-zinc-800 font-mono">{value}</dd>
    </>
  );
}
