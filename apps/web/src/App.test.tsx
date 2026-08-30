import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import App from './App';
import { apiPost } from './api/client';
import type { BackendClaim } from './api/claims';

vi.mock('./api/client', () => ({ apiPost: vi.fn(), apiGet: vi.fn() }));

const resetClaim: BackendClaim = {
  id: 1,
  claim_number: 'CLM-2026-0817',
  policy_id: 1,
  stage: 'DOCUMENT_VERIFICATION',
  document_completeness: 'INCOMPLETE',
  survey_required: false,
  claim_type: 'MOTOR',
  status: 'OPEN',
  documents: [
    { id: 1, document_type: 'DAMAGE_PHOTO', file_name: 'damage_front_bumper.jpg', file_url: 'x', status: 'VERIFIED' },
  ],
  events: [
    {
      id: 1,
      actor_name: 'Supervisor',
      actor_type: 'AGENT',
      action: 'CLAIM_INTAKE_INITIALIZED',
      payload: '{}',
      created_at: '2026-08-29T09:10:00Z',
    },
  ],
  assignment: null,
  recommendation: null,
  stage_sla_due_at: '2026-08-29T09:35:00Z',
};

const assessedClaim: BackendClaim = {
  ...resetClaim,
  document_completeness: 'COMPLETE',
  stage: 'DECISION',
  documents: [
    ...resetClaim.documents!,
    { id: 2, document_type: 'POLICE_REPORT', file_name: 'police_report.pdf', file_url: 'x', status: 'VERIFIED' },
  ],
  assignment: {
    officer_id: 2,
    officer: { name: 'J. Lim' },
    workload_score: 8,
    skill_score: 4,
    total_score: 6.5,
  },
  recommendation: {
    outcome: 'APPROVE',
    confidence: 0.94,
    reasons: 'Damage matches barrier collision. Policy in-force. Loss substantiated.',
  },
};

const closedClaim: BackendClaim = { ...assessedClaim, stage: 'CLOSED', status: 'CLOSED' };

const mockedPost = apiPost as unknown as ReturnType<typeof vi.fn>;

describe('Claims Ops Console', () => {
  beforeEach(() => {
    mockedPost.mockReset();
  });

  it('renders initial state and seeds a demo Claim from the backend', async () => {
    mockedPost.mockResolvedValueOnce(resetClaim);

    render(<App />);
    expect(screen.getByText('Taskmaster Claims Ops')).toBeInTheDocument();
    expect(screen.getByText('No Claim loaded')).toBeInTheDocument();

    fireEvent.click(screen.getAllByRole('button', { name: /seed demo claim/i })[0]);
    await waitFor(() => expect(screen.getByText('CLM-2026-0817')).toBeInTheDocument());
    expect(screen.getByText('Incomplete')).toBeInTheDocument();
  });

  it('uploads police report and renders assignment and recommendation', async () => {
    mockedPost.mockResolvedValueOnce(resetClaim);

    render(<App />);
    fireEvent.click(screen.getAllByRole('button', { name: /seed demo claim/i })[0]);
    await waitFor(() => expect(screen.getByText('CLM-2026-0817')).toBeInTheDocument());

    mockedPost
      .mockResolvedValueOnce({ ...resetClaim, document_completeness: 'COMPLETE', assignment: assessedClaim.assignment })
      .mockResolvedValueOnce(assessedClaim);

    fireEvent.change(screen.getByLabelText('Upload police report'), {
      target: { files: [new File(['dummy'], 'police.pdf', { type: 'application/pdf' })] },
    });

    await waitFor(() => expect(screen.getByText('Complete')).toBeInTheDocument());
    expect(screen.getByText('J. Lim')).toBeInTheDocument();
    expect(screen.getByText('APPROVE')).toBeInTheDocument();
  });

  it('requires confirmation to approve and closes the Claim', async () => {
    mockedPost.mockResolvedValueOnce(resetClaim);

    render(<App />);
    fireEvent.click(screen.getAllByRole('button', { name: /seed demo claim/i })[0]);
    await waitFor(() => expect(screen.getByText('CLM-2026-0817')).toBeInTheDocument());

    mockedPost
      .mockResolvedValueOnce({ ...resetClaim, document_completeness: 'COMPLETE', assignment: assessedClaim.assignment })
      .mockResolvedValueOnce(assessedClaim)
      .mockResolvedValueOnce(closedClaim);

    fireEvent.change(screen.getByLabelText('Upload police report'), {
      target: { files: [new File(['dummy'], 'police.pdf')] },
    });
    await waitFor(() => expect(screen.getByText('APPROVE')).toBeInTheDocument());

    const approveButton = screen.getByRole('button', { name: /^approve$/i });
    expect(approveButton).toBeDisabled();

    fireEvent.click(screen.getByLabelText(/I confirm this creates a Decision/i));
    expect(approveButton).toBeEnabled();

    fireEvent.click(approveButton);
    await waitFor(() =>
      expect(screen.getByText(/Assessment PDF Stored|claim_assessment_report/)).toBeInTheDocument(),
    );
    expect(screen.queryByRole('button', { name: /^approve$/i })).not.toBeInTheDocument();
  });
});
