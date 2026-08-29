import { Claim } from '../types';
import { Card } from './Card';

interface AssignmentPanelProps {
  claim: Claim;
}

export function AssignmentPanel({ claim }: AssignmentPanelProps) {
  return (
    <Card>
      <h3 className="font-semibold text-sm text-zinc-800 uppercase tracking-wide mb-3">Assignment</h3>
      {!claim.assignment ? (
        <p className="text-sm text-zinc-500 italic">Pending classification and assignment score computation.</p>
      ) : (
        <div>
          <div className="mb-3 text-sm text-zinc-700">
            Assigned owner: <strong className="text-zinc-900">{claim.assignment.winnerName}</strong> (Highest score)
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-zinc-50 text-zinc-600 uppercase text-xs border-b border-zinc-200">
                <tr>
                  <th className="py-2 px-3">Claims Officer</th>
                  <th className="py-2 px-3">Skill</th>
                  <th className="py-2 px-3">Workload</th>
                  <th className="py-2 px-3">Score</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 font-mono text-xs">
                {claim.assignment.scores.map((score) => (
                  <tr
                    key={score.officerId}
                    className={
                      score.officerId === claim.assignment?.winnerId
                        ? 'bg-teal-50/50 font-semibold text-teal-900'
                        : 'text-zinc-700'
                    }
                  >
                    <td className="py-2 px-3 font-sans font-medium">
                      {score.officerName} {score.officerId === claim.assignment?.winnerId && '(winner)'}
                    </td>
                    <td className="py-2 px-3">{score.skill}</td>
                    <td className="py-2 px-3">{score.workload}</td>
                    <td className="py-2 px-3 font-bold">{score.score}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </Card>
  );
}
