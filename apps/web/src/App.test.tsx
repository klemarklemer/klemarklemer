import { render, screen, fireEvent } from '@testing-library/react';
import App from './App';

describe('Claims Ops Console', () => {
  it('renders initial state and allows seeding demo claim', () => {
    render(<App />);
    expect(screen.getByText('Taskmaster Claims Ops')).toBeInTheDocument();
    expect(screen.getByText('No Claim loaded')).toBeInTheDocument();

    const seedButton = screen.getAllByRole('button', { name: /seed demo claim/i })[0];
    fireEvent.click(seedButton);

    expect(screen.getByText('CLM-2026-0817')).toBeInTheDocument();
    expect(screen.getByText('Incomplete')).toBeInTheDocument();
  });

  it('allows uploading police report to complete documents and view assignment and recommendation', () => {
    render(<App />);
    fireEvent.click(screen.getAllByRole('button', { name: /seed demo claim/i })[0]);

    expect(screen.getByText('Incomplete')).toBeInTheDocument();

    const uploadLabel = screen.getByLabelText('Upload police report');
    fireEvent.change(uploadLabel, { target: { files: [new File(['dummy'], 'police.pdf', { type: 'application/pdf' })] } });

    expect(screen.getByText('Complete')).toBeInTheDocument();
    expect(screen.getByText('J. Lim')).toBeInTheDocument();
    expect(screen.getByText('APPROVE')).toBeInTheDocument();
  });

  it('requires confirmation to approve and closes claim', () => {
    render(<App />);
    fireEvent.click(screen.getAllByRole('button', { name: /seed demo claim/i })[0]);
    fireEvent.change(screen.getByLabelText('Upload police report'), { target: { files: [new File(['dummy'], 'police.pdf')] } });

    const approveButton = screen.getByRole('button', { name: /^approve$/i });
    expect(approveButton).toBeDisabled();

    const checkbox = screen.getByLabelText(/I confirm this creates a Decision/i);
    fireEvent.click(checkbox);

    expect(approveButton).toBeEnabled();

    fireEvent.click(approveButton);

    expect(screen.getByText('Assessment PDF Stored')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /^approve$/i })).not.toBeInTheDocument();
  });
});
