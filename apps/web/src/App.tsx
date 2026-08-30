import { useEffect, useState } from 'react';
import { Claim } from './types';
import {
  resetDemo,
  uploadDocument,
  runAssessment,
  submitApproval,
  toClaimView,
  listClaims,
  getClaim,
  type TriageItem,
} from './api/claims';
import { AppHeader } from './components/AppHeader';
import { Toast } from './components/Toast';
import { EmptyState } from './components/EmptyState';
import { TriageStrip } from './components/TriageStrip';
import { ClaimIdentity } from './components/ClaimIdentity';
import { DocumentCompleteness } from './components/DocumentCompleteness';
import { AssignmentPanel } from './components/AssignmentPanel';
import { SurveyPanel } from './components/SurveyPanel';
import { RecommendationPanel } from './components/RecommendationPanel';
import { TimelinePanel } from './components/TimelinePanel';
import { LiveRegion } from './components/LiveRegion';

export default function App() {
  const [claim, setClaim] = useState<Claim | null>(null);
  const [claimId, setClaimId] = useState<number | null>(null);
  const [triage, setTriage] = useState<TriageItem[]>([]);
  const [busy, setBusy] = useState(false);
  const [showNotifications, setShowNotifications] = useState(false);
  const [confirmApprove, setConfirmApprove] = useState(false);
  const [toast, setToast] = useState<string | null>(null);

  const showToast = (message: string) => {
    setToast(message);
    setTimeout(() => setToast(null), 4000);
  };

  const loadTriage = async () => {
    try {
      setTriage(await listClaims());
    } catch {
      setTriage([]);
    }
  };

  useEffect(() => {
    void loadTriage();
  }, []);

  const handleSeed = async () => {
    setBusy(true);
    try {
      const seeded = await resetDemo();
      setClaim(toClaimView(seeded));
      setClaimId(seeded.id);
      setConfirmApprove(false);
      showToast('Demo Claims created from backend');
      await loadTriage();
    } catch (error) {
      showToast(error instanceof Error ? error.message : 'Failed to seed Claim');
    } finally {
      setBusy(false);
    }
  };

  const handleSelectClaim = async (id: number) => {
    setBusy(true);
    try {
      const selected = await getClaim(id);
      setClaim(toClaimView(selected));
      setClaimId(selected.id);
      setConfirmApprove(false);
      showToast(`Loaded ${selected.claim_number}`);
    } catch (error) {
      showToast(error instanceof Error ? error.message : 'Failed to load Claim');
    } finally {
      setBusy(false);
    }
  };

  const handleReset = () => {
    setClaim(null);
    setClaimId(null);
    setConfirmApprove(false);
    showToast('Workspace reset');
  };

  const handleUpload = async () => {
    if (claimId === null) return;
    setBusy(true);
    try {
      await uploadDocument(claimId);
      const assessed = await runAssessment(claimId);
      setClaim(toClaimView(assessed));
      showToast('Police report uploaded and assessment completed');
    } catch (error) {
      showToast(error instanceof Error ? error.message : 'Failed to upload document');
    } finally {
      setBusy(false);
    }
  };

  const closeClaim = async (action: 'APPROVE' | 'REJECT' | 'REJECT_FRAUD') => {
    if (claimId === null || !claim) return;
    const officerId = claim.assignment ? Number(claim.assignment.winnerId) : 0;
    setBusy(true);
    try {
      const closed = await submitApproval(claimId, action, officerId);
      setClaim(toClaimView(closed));
      const msg = action === 'APPROVE'
        ? 'Recommendation approved. Decision recorded and Claim closed.'
        : action === 'REJECT_FRAUD'
        ? 'Fraud claim rejected. Decision recorded and Claim closed.'
        : 'Recommendation rejected. Claim remains open.';
      showToast(msg);
    } catch (error) {
      showToast(error instanceof Error ? error.message : 'Failed to record decision');
    } finally {
      setBusy(false);
    }
  };

  const handleApprove = () => {
    if (!confirmApprove) return;
    void closeClaim('APPROVE');
  };

  const handleReject = () => {
    const action = claim?.fraudSignal ? 'REJECT_FRAUD' : 'REJECT';
    void closeClaim(action);
  };

  const isComplete = claim ? claim.documentCompleteness === 'COMPLETE' : false;
  const latestNotification =
    claim && claim.notifications.length > 0 ? claim.notifications[claim.notifications.length - 1] : '';
  const liveMessage = claim
    ? `${claim.sla.atRisk ? 'Stage SLA clock at risk. ' : ''}${latestNotification}`
    : '';

  return (
    <div className="min-h-dvh flex flex-col bg-zinc-100 font-sans text-zinc-900">
      <AppHeader
        claim={claim}
        notificationsOpen={showNotifications}
        onToggleNotifications={() => setShowNotifications(!showNotifications)}
        onSeed={handleSeed}
        onReset={handleReset}
        busy={busy}
      />

      <Toast message={toast} />
      <LiveRegion message={liveMessage} />

      <main className="flex-1 max-w-7xl w-full mx-auto p-4 md:p-6">
        <TriageStrip items={triage} selectedId={claimId} onSelect={handleSelectClaim} />

        {!claim ? (
          <EmptyState onSeed={handleSeed} busy={busy} />
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <div className="lg:col-span-2 space-y-4">
              <ClaimIdentity claim={claim} />
              <DocumentCompleteness claim={claim} isComplete={isComplete} onUpload={handleUpload} busy={busy} />
              <AssignmentPanel claim={claim} />
              {claim.survey?.required && <SurveyPanel claim={claim} />}
              <RecommendationPanel
                claim={claim}
                confirmApprove={confirmApprove}
                onConfirmChange={setConfirmApprove}
                onApprove={handleApprove}
                onReject={handleReject}
                busy={busy}
              />
            </div>

            <div className="space-y-4">
              <TimelinePanel claim={claim} />
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
