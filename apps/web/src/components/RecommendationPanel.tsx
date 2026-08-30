import { Claim } from '../types';
import { Card } from './Card';
import { AssessmentEvidence } from './AssessmentEvidence';

interface RecommendationPanelProps {
  claim: Claim;
  confirmApprove: boolean;
  onConfirmChange: (checked: boolean) => void;
  onApprove: () => void;
  onReject: () => void;
  busy?: boolean;
}

export function RecommendationPanel({
  claim,
  confirmApprove,
  onConfirmChange,
  onApprove,
  onReject,
  busy
}: RecommendationPanelProps) {
  return (
    <Card>
      <h3 className="font-semibold text-sm text-zinc-800 uppercase tracking-wide mb-3">
        Assessment recommendation &amp; Human approval
      </h3>
      {!claim.recommendation ? (
        <p className="text-sm text-zinc-500 italic">
          Available after document completeness is complete and assignment is finalised.
        </p>
      ) : (
        <div className="space-y-4">
          <div className="bg-zinc-50 border border-zinc-200/80 rounded-md p-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-xs uppercase font-medium text-zinc-500">Proposed Outcome</span>
              <span className="font-mono text-xs font-semibold bg-teal-100 text-teal-800 px-2.5 py-0.5 rounded">
                {claim.recommendation.outcome}
              </span>
            </div>
            <ul className="list-disc list-inside text-sm text-zinc-700 space-y-1 mb-3">
              {claim.recommendation.reasons.map((reason, index) => (
                <li key={index}>{reason}</li>
              ))}
            </ul>
            <div className="flex items-center justify-between text-xs">
              <div className="text-zinc-500">
                Confidence:{' '}
                <span className="font-mono font-semibold text-zinc-800">
                  {claim.recommendation.confidence * 100}%
                </span>
              </div>
              <div className="text-zinc-500 flex items-center gap-1">
                <span className="font-mono text-[10px] uppercase">Powered by Gemini</span>
              </div>
            </div>
          </div>

          <AssessmentEvidence claim={claim} />

            {claim.closed ? (
              <div className="p-4 bg-zinc-900 text-white rounded-md flex items-center justify-between">
                <div>
                  <div className="text-xs text-zinc-400">Recorded Decision</div>
                  <div className="font-bold text-base">{claim.decision}</div>
                </div>
                <span className="text-xs bg-teal-800 text-teal-100 px-2.5 py-1 rounded">
                  {claim.reportLabel ?? 'Assessment PDF Stored'}
                </span>
              </div>
            ) : claim.recommendation.outcome === 'MANUAL_REVIEW' ? (
              <div className="p-3 bg-amber-50 border border-amber-200 rounded-md space-y-2">
                <div className="text-xs font-semibold text-amber-800 uppercase">Mandatory human review required</div>
                <p className="text-xs text-amber-700 leading-relaxed">
                  Fraud signals detected. Auto-approval is blocked by orchestration policy. Review investigative findings before issuing a decision.
                </p>
                <div className="pt-2">
                  <button
                    onClick={onReject}
                    disabled={busy}
                    className="w-full border border-red-300 text-red-700 hover:bg-red-50 focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 py-2.5 px-4 rounded-md text-sm font-medium transition cursor-pointer disabled:opacity-60"
                  >
                    Reject Claim (Fraud Override)
                  </button>
                </div>
              </div>
            ) : (
              <>
                {claim.decision && (
                  <p className="text-xs text-amber-700 bg-amber-50 border border-amber-200 rounded px-3 py-2 mb-1">
                    Decision recorded: {claim.decision}. Claim remains open; you can still Approve.
                  </p>
                )}
                <div className="space-y-3 pt-2 border-t border-zinc-100">
                  <div className="flex items-start gap-2">
                    <input
                      id="confirm-decision"
                      type="checkbox"
                      checked={confirmApprove}
                      onChange={(event) => onConfirmChange(event.target.checked)}
                      className="mt-1 h-4 w-4 rounded border-zinc-300 text-teal-700 focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2"
                    />
                    <label htmlFor="confirm-decision" className="text-xs text-zinc-600 leading-relaxed cursor-pointer">
                      I confirm this creates a Decision and closes this Claim.
                    </label>
                  </div>

                  <div className="flex items-center gap-3">
                    <button
                      onClick={onApprove}
                      disabled={!confirmApprove || busy}
                      className={`flex-1 py-2.5 px-4 rounded-md text-sm font-medium transition cursor-pointer focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 ${
                        confirmApprove
                          ? 'bg-teal-700 hover:bg-teal-800 text-white shadow-xs'
                          : 'bg-zinc-200 text-zinc-400 cursor-not-allowed'
                      }`}
                    >
                      Approve
                    </button>
                    <button
                      onClick={onReject}
                      disabled={busy}
                      className="flex-1 border border-red-300 text-red-700 hover:bg-red-50 focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 py-2.5 px-4 rounded-md text-sm font-medium transition cursor-pointer disabled:opacity-60"
                    >
                      Reject
                    </button>
                  </div>
                </div>
              </>
            )}
        </div>
      )}
    </Card>
  );
}
