import { useState } from 'react';
import { Claim } from './types';
import { createSeedClaim, completeAfterUpload, approveClaim, rejectClaim } from './data/claim';
import { AppHeader } from './components/AppHeader';
import { Toast } from './components/Toast';
import { EmptyState } from './components/EmptyState';
import { ClaimIdentity } from './components/ClaimIdentity';
import { DocumentCompleteness } from './components/DocumentCompleteness';
import { AssignmentPanel } from './components/AssignmentPanel';
import { RecommendationPanel } from './components/RecommendationPanel';
import { TimelinePanel } from './components/TimelinePanel';
import { LiveRegion } from './components/LiveRegion';

export default function App() {
  const [claim, setClaim] = useState<Claim | null>(null);
  const [showNotifications, setShowNotifications] = useState(false);
  const [confirmApprove, setConfirmApprove] = useState(false);
  const [toast, setToast] = useState<string | null>(null);

  const handleSeed = () => {
    const seeded = createSeedClaim();
    setClaim(seeded);
    setConfirmApprove(false);
    showToast('Synthetic demo Claim seeded');
  };

  const handleReset = () => {
    setClaim(null);
    setConfirmApprove(false);
    showToast('Workspace reset');
  };

  const handleUpload = () => {
    if (!claim) return;
    const updated = completeAfterUpload(claim);
    setClaim(updated);
    showToast('Police report uploaded successfully');
  };

  const handleApprove = () => {
    if (!claim || !confirmApprove) return;
    const closed = approveClaim(claim);
    setClaim(closed);
    showToast('Recommendation approved. Decision recorded and Claim closed.');
  };

  const handleReject = () => {
    if (!claim) return;
    const rejected = rejectClaim(claim);
    setClaim(rejected);
    showToast('Recommendation rejected. Claim remains open.');
  };

  const showToast = (message: string) => {
    setToast(message);
    setTimeout(() => setToast(null), 4000);
  };

  const isComplete = claim ? claim.documents.every((document) => document.present) : false;
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
      />

      <Toast message={toast} />
      <LiveRegion message={liveMessage} />

      <main className="flex-1 max-w-7xl w-full mx-auto p-4 md:p-6">
        {!claim ? (
          <EmptyState onSeed={handleSeed} />
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <div className="lg:col-span-2 space-y-4">
              <ClaimIdentity claim={claim} />
              <DocumentCompleteness claim={claim} isComplete={isComplete} onUpload={handleUpload} />
              <AssignmentPanel claim={claim} />
              <RecommendationPanel
                claim={claim}
                confirmApprove={confirmApprove}
                onConfirmChange={setConfirmApprove}
                onApprove={handleApprove}
                onReject={handleReject}
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
