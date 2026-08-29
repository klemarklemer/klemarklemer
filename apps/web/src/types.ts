export type Stage = 'Intake' | 'Assignment' | 'Assessment' | 'Decision' | 'Closed';

export interface DocumentItem {
  id: string;
  name: string;
  required: boolean;
  present: boolean;
}

export interface AssignmentScore {
  officerId: string;
  officerName: string;
  skill: number;
  workload: number;
  score: number;
}

export interface Assignment {
  scores: AssignmentScore[];
  winnerId: string;
  winnerName: string;
}

export interface Recommendation {
  outcome: 'APPROVE' | 'REJECT' | 'MANUAL_REVIEW';
  reasons: string[];
  confidence: number;
}

export interface ClaimEvent {
  id: string;
  at: string;
  actor: string;
  action: string;
  detail?: string;
}

export interface SlaClock {
  deadline: string;
  remainingMinutes: number;
  atRisk: boolean;
}

export interface Claim {
  id: string;
  policyId: string;
  line: string;
  stage: Stage;
  documents: DocumentItem[];
  assignment: Assignment | null;
  recommendation: Recommendation | null;
  decision: 'APPROVE' | 'REJECT' | null;
  closed: boolean;
  events: ClaimEvent[];
  notifications: string[];
  sla: SlaClock;
}
