import { apiGet, apiPost } from './client';
import type {
  Assignment,
  Claim,
  ClaimEvent,
  DocumentItem,
  PolicySummary,
  Recommendation,
  SlaClock,
  Stage,
  SurveyInfo,
} from '../types';

export interface BackendDocument {
  id: number;
  document_type: string;
  file_name: string;
  file_url: string;
  status: string;
  extracted_data?: string;
}

export interface BackendEvent {
  id: number;
  actor_name: string;
  actor_type: string;
  action: string;
  payload: string;
  created_at: string;
}

export interface BackendAssignment {
  officer_id: number;
  officer?: { name: string };
  workload_score: number;
  skill_score: number;
  total_score: number;
}

export interface BackendSurvey {
  surveyor_id?: number | null;
  survey_status?: string | null;
  survey_sla_due_at?: string | null;
  survey_completed_at?: string | null;
  survey_report_url?: string | null;
  survey_photos?: string[];
}

export interface BackendRecommendation {
  outcome: string;
  confidence: number;
  reasons: string;
}

export interface BackendPolicy {
  policy_number: string;
  policy_holder_name: string;
  vehicle_plate: string;
  vehicle_model: string;
  coverage_type: string;
  max_coverage_amount: number;
  deductible_amount: number;
  status: string;
}

export interface BackendClaim {
  id: number;
  claim_number: string;
  policy_id: number;
  stage: string;
  document_completeness: string;
  survey_required: boolean;
  surveyor_id?: number | null;
  surveyor?: { id: number; name: string; role?: string; specialty?: string } | null;
  survey_status?: string | null;
  survey_sla_due_at?: string | null;
  survey_completed_at?: string | null;
  survey_report_url?: string | null;
  claim_type: string;
  status: string;
  documents?: BackendDocument[];
  policy?: BackendPolicy | null;
  events?: BackendEvent[];
  assignment?: BackendAssignment | null;
  recommendation?: BackendRecommendation | null;
  claim_sla_due_at?: string | null;
  stage_sla_due_at?: string | null;
  fraud_signal?: string | null;
}

export interface TriageItem {
  id: number;
  claimNumber: string;
  line: string;
  stage: string;
  atRisk: boolean;
  fraud: boolean;
}

export function resetDemo(): Promise<BackendClaim> {
  return apiPost<BackendClaim>('/v1/demo/reset');
}

export function uploadDocument(id: number): Promise<BackendClaim> {
  return apiPost<BackendClaim>(`/v1/claim/${id}/documents`);
}

export function runAssessment(id: number): Promise<BackendClaim> {
  return apiPost<BackendClaim>(`/v1/claim/${id}/assessment`);
}

export function submitApproval(
  id: number,
  action: 'APPROVE' | 'REJECT' | 'REJECT_FRAUD',
  officerId: number,
  notes = '',
): Promise<BackendClaim> {
  return apiPost<BackendClaim>(`/v1/claim/${id}/approval`, { officer_id: officerId, action, notes });
}

export function listClaims(): Promise<TriageItem[]> {
  return apiGet<BackendClaim[]>(`/v1/claim`).then((items) => items.map(toTriageItem));
}

export function getClaim(id: number): Promise<BackendClaim> {
  return apiGet<BackendClaim>(`/v1/claim/${id}`);
}

function computeAtRisk(due?: string | null): boolean {
  if (!due) return false;
  return Math.round((new Date(due).getTime() - Date.now()) / 60000) < 30;
}

const STAGE_MAP: Record<string, Stage> = {
  INTAKE: 'Intake',
  DOCUMENT_VERIFICATION: 'Intake',
  ASSIGNMENT: 'Assignment',
  SURVEY: 'Survey',
  ASSESSMENT: 'Assessment',
  DECISION: 'Decision',
  CLOSED: 'Closed',
};

function toStage(stage: string): Stage {
  return STAGE_MAP[stage] ?? 'Intake';
}

function splitReasons(text: string): string[] {
  return text
    .split(/\.\s+|\n+/)
    .map((part) => part.trim())
    .filter(Boolean)
    .map((part) => (part.endsWith('.') ? part : `${part}.`));
}

export function toClaimView(b: BackendClaim): Claim {
  const documents: DocumentItem[] = (b.documents ?? []).map((doc) => ({
    id: String(doc.id),
    name: doc.file_name,
    required: true,
    present: doc.status === 'VERIFIED',
    extractedData: doc.extracted_data || undefined,
  }));

  const policy: PolicySummary | undefined = b.policy
    ? {
        policyNumber: b.policy.policy_number,
        policyHolder: b.policy.policy_holder_name,
        vehicle: `${b.policy.vehicle_model} (${b.policy.vehicle_plate})`,
        coverage: b.policy.coverage_type,
        maxCoverage: b.policy.max_coverage_amount,
        deductible: b.policy.deductible_amount,
        status: b.policy.status,
      }
    : undefined;

  const isComplete = b.document_completeness === 'COMPLETE';

  const assignment: Assignment | null = b.assignment
    ? {
        winnerId: String(b.assignment.officer_id),
        winnerName: b.assignment.officer?.name ?? `Officer ${b.assignment.officer_id}`,
        scores: [
          {
            officerId: String(b.assignment.officer_id),
            officerName: b.assignment.officer?.name ?? `Officer ${b.assignment.officer_id}`,
            skill: b.assignment.skill_score,
            workload: b.assignment.workload_score,
            score: b.assignment.total_score,
          },
        ],
      }
    : null;

  const survey: SurveyInfo | null = b.survey_required
    ? {
        required: true,
        status: b.survey_status || undefined,
        surveyorId: b.surveyor_id || undefined,
        surveyorName: b.surveyor?.name || undefined,
        completedAt: b.survey_completed_at || undefined,
        reportUrl: b.survey_report_url || undefined,
      }
    : { required: false };

  const recommendation: Recommendation | null = b.recommendation
    ? {
        outcome: b.recommendation.outcome as Recommendation['outcome'],
        reasons: splitReasons(b.recommendation.reasons),
        confidence: b.recommendation.confidence,
      }
    : null;

  const closed = b.stage === 'CLOSED';

  // Determine decision and report label from the last DECISION_ISSUED event
  let decision: 'APPROVE' | 'REJECT' | null = null;
  let reportLabel: string | undefined = undefined;
  if (closed) {
    const decisionEvent = (b.events ?? []).slice().reverse().find(
      (e) => e.action === 'DECISION_ISSUED' && e.payload,
    );
    if (decisionEvent) {
      try {
        const payload = JSON.parse(decisionEvent.payload);
        decision = payload.binding_outcome === 'REJECT' ? 'REJECT' : 'APPROVE';
        if (payload.generated_report) {
          reportLabel = payload.generated_report;
        }
      } catch {
        decision = 'APPROVE';
      }
    } else {
      decision = 'APPROVE';
    }
  }

  const events: ClaimEvent[] = (b.events ?? []).map((event) => ({
    id: String(event.id),
    at: event.created_at,
    actor: event.actor_name,
    action: event.action,
    detail: event.payload || undefined,
  }));

  const notifications: string[] = [];
  if (!isComplete) notifications.push('Document completeness incomplete: police report missing');
  if (assignment) notifications.push(`Assigned to ${assignment.winnerName}`);
  if (closed) notifications.push('Claim closed by human approval');

  const due = b.stage_sla_due_at ?? b.claim_sla_due_at ?? null;
  const sla: SlaClock = { deadline: due ?? '', remainingMinutes: 0, atRisk: false };
  if (due) {
    sla.remainingMinutes = Math.round((new Date(due).getTime() - Date.now()) / 60000);
    sla.atRisk = computeAtRisk(due);
  }

  return {
    id: b.claim_number,
    policyId: String(b.policy_id),
    line: b.claim_type,
    stage: toStage(b.stage),
    documentCompleteness: b.document_completeness,
    documents,
    policy,
    assignment,
    survey,
    recommendation,
    fraudSignal: b.fraud_signal || undefined,
    decision,
    closed,
    events,
    notifications,
    sla,
    reportLabel,
  };
}

export function toTriageItem(b: BackendClaim): TriageItem {
  return {
    id: b.id,
    claimNumber: b.claim_number,
    line: b.claim_type,
    stage: toStage(b.stage),
    atRisk: computeAtRisk(b.stage_sla_due_at ?? b.claim_sla_due_at),
    fraud: Boolean(b.fraud_signal),
  };
}
