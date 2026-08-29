import type { Claim, ClaimEvent } from '../types';

function event(id: string, at: string, actor: string, action: string, detail?: string): ClaimEvent {
  return { id, at, actor, action, detail };
}

export function createSeedClaim(): Claim {
  return {
    id: 'CLM-2026-0817',
    policyId: 'POL-88231',
    line: 'MOTOR',
    stage: 'Intake',
    documents: [
      { id: 'doc-claim-form', name: 'Claim form', required: true, present: true },
      { id: 'doc-id-proof', name: 'ID proof', required: true, present: true },
      { id: 'doc-police-report', name: 'Police report', required: true, present: false }
    ],
    assignment: null,
    recommendation: null,
    decision: null,
    closed: false,
    events: [
      event('evt-seed', '2026-08-29 09:10 UTC', 'system', 'Seed demo Claim', 'Synthetic motor Claim created'),
      event('evt-intake', '2026-08-29 09:11 UTC', 'intake-agent', 'Intake', 'Document completeness recorded: incomplete')
    ],
    notifications: ['Document completeness incomplete: police report missing'],
    sla: {
      deadline: '2026-08-29 09:20 UTC',
      remainingMinutes: -3,
      atRisk: true
    }
  };
}

export function completeAfterUpload(claim: Claim): Claim {
  const documents = claim.documents.map((document) =>
    document.id === 'doc-police-report' ? { ...document, present: true } : document
  );

  const assignment = {
    scores: [
      { officerId: 'off-1', officerName: 'A. Rahman', skill: 0.92, workload: 0.78, score: 0.81 },
      { officerId: 'off-2', officerName: 'J. Lim', skill: 0.85, workload: 0.35, score: 0.9 },
      { officerId: 'off-3', officerName: 'S. Patel', skill: 0.71, workload: 0.2, score: 0.76 }
    ],
    winnerId: 'off-2',
    winnerName: 'J. Lim'
  };

  const recommendation = {
    outcome: 'APPROVE' as const,
    reasons: [
      'Policy POL-88231 covers third-party motor damage for the reported incident date.',
      'Police report matches the Claim identifiers and incident description.',
      'Document completeness is complete for the current Stage.'
    ],
    confidence: 0.87
  };

  const events = [
    ...claim.events,
    event('evt-upload', '2026-08-29 09:18 UTC', 'claims-officer', 'Upload', 'Police report uploaded'),
    event('evt-reintake', '2026-08-29 09:18 UTC', 'intake-agent', 'Intake', 'Document completeness recorded: complete'),
    event('evt-classify', '2026-08-29 09:19 UTC', 'intake-agent', 'Classification', 'MOTOR / MEDIUM; Survey not required'),
    event('evt-assignment', '2026-08-29 09:20 UTC', 'assignment-agent', 'Assignment', 'Assigned to J. Lim (score 0.90)'),
    event('evt-sla', '2026-08-29 09:21 UTC', 'sla-agent', 'SLA reminder', 'Stage SLA clock at risk'),
    event('evt-assessment', '2026-08-29 09:25 UTC', 'assessment-agent', 'Assessment recommendation', 'APPROVE with confidence 0.87')
  ];

  return {
    ...claim,
    stage: 'Assessment',
    documents,
    assignment,
    recommendation,
    events,
    notifications: [...claim.notifications, 'Stage SLA clock at risk']
  };
}

export function approveClaim(claim: Claim): Claim {
  const events = [
    ...claim.events,
    event('evt-human-approval', '2026-08-29 09:30 UTC', 'claims-officer', 'Human approval', 'Recommendation approved'),
    event('evt-decision', '2026-08-29 09:30 UTC', 'system', 'Decision', 'APPROVE'),
    event('evt-close', '2026-08-29 09:31 UTC', 'system', 'Claim closed', 'Assessment report PDF stored')
  ];
  return {
    ...claim,
    stage: 'Closed',
    decision: 'APPROVE',
    closed: true,
    events,
    notifications: [...claim.notifications, 'Decision recorded: APPROVE']
  };
}

export function rejectClaim(claim: Claim): Claim {
  const events = [
    ...claim.events,
    event('evt-human-approval-reject', '2026-08-29 09:30 UTC', 'claims-officer', 'Human approval', 'Recommendation rejected'),
    event('evt-decision-reject', '2026-08-29 09:30 UTC', 'system', 'Decision', 'REJECT')
  ];
  return {
    ...claim,
    decision: 'REJECT',
    events,
    notifications: [...claim.notifications, 'Decision recorded: REJECT; Claim remains open']
  };
}
