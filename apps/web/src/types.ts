export type Stage = 'Intake' | 'Assignment' | 'Survey' | 'Assessment' | 'Decision' | 'Closed';

export interface PolicySummary {
  policyNumber: string;
  policyHolder: string;
  vehicle: string;
  coverage: string;
  maxCoverage: number;
  deductible: number;
  status: string;
}

export interface DocumentItem {
  id: string;
  name: string;
  required: boolean;
  present: boolean;
  extractedData?: string;
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

export interface SurveyInfo {
  required: boolean;
  status?: string;
  surveyorId?: number;
  surveyorName?: string;
  completedAt?: string;
  reportUrl?: string;
}

export interface Claim {
  id: string;
  policyId: string;
  line: string;
  stage: Stage;
  documentCompleteness: string;
  documents: DocumentItem[];
  policy?: PolicySummary;
  assignment: Assignment | null;
  survey?: SurveyInfo | null;
  recommendation: Recommendation | null;
  fraudSignal?: string;
  decision: 'APPROVE' | 'REJECT' | null;
  closed: boolean;
  events: ClaimEvent[];
  notifications: string[];
  sla: SlaClock;
  reportLabel?: string;
}
