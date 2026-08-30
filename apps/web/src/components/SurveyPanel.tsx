import { MapPin, Clock, CheckCircle, Warning, FileText } from '@phosphor-icons/react';
import { Claim } from '../types';
import { Card } from './Card';

interface SurveyPanelProps {
  claim: Claim;
}

export function SurveyPanel({ claim }: SurveyPanelProps) {
  if (!claim.survey?.required) return null;

  const survey = claim.survey;
  const isCompleted = survey.status === 'COMPLETED';

  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'COMPLETED':
        return 'bg-teal-50 text-teal-700 border-teal-200';
      case 'IN_PROGRESS':
        return 'bg-amber-50 text-amber-700 border-amber-200';
      case 'ASSIGNED':
        return 'bg-blue-50 text-blue-700 border-blue-200';
      case 'OVERDUE':
        return 'bg-red-50 text-red-700 border-red-200';
      default:
        return 'bg-zinc-50 text-zinc-700 border-zinc-200';
    }
  };

  const getStatusIcon = (status?: string) => {
    switch (status) {
      case 'COMPLETED':
        return <CheckCircle size={14} className="text-teal-600" />;
      case 'IN_PROGRESS':
        return <Clock size={14} className="text-amber-600" />;
      case 'ASSIGNED':
        return <MapPin size={14} className="text-blue-600" />;
      case 'OVERDUE':
        return <Warning size={14} className="text-red-600" />;
      default:
        return <Warning size={14} className="text-zinc-600" />;
    }
  };

  return (
    <Card>
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-semibold text-sm text-zinc-800 uppercase tracking-wide">Survey</h3>
        <span
          className={`text-xs px-2 py-0.5 rounded font-medium ${getStatusColor(survey.status)}`}
        >
          {getStatusIcon(survey.status)}
          {survey.status ? survey.status.replace('_', ' ') : 'Pending'}
        </span>
      </div>

      <dl className="space-y-2 text-sm">
        {survey.surveyorName && (
          <div className="flex items-center justify-between py-1.5 border-b border-zinc-100">
            <dt className="text-zinc-500">Assigned Surveyor</dt>
            <dd className="text-zinc-900 font-mono font-medium">{survey.surveyorName}</dd>
          </div>
        )}
        {survey.completedAt && (
          <div className="flex items-center justify-between py-1.5 border-b border-zinc-100">
            <dt className="text-zinc-500">Completed</dt>
            <dd className="text-zinc-900 font-mono">{new Date(survey.completedAt).toLocaleString()}</dd>
          </div>
        )}
        {survey.reportUrl && (
          <div className="flex items-center justify-between py-1.5 border-b border-zinc-100">
            <dt className="text-zinc-500">Report</dt>
            <dd className="flex items-center gap-2">
              <FileText size={14} className="text-teal-600" />
              <a
                href={survey.reportUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="text-teal-700 hover:underline text-xs font-mono"
              >
                View Survey Report
              </a>
            </dd>
          </div>
        )}
      </dl>

      {!isCompleted && (
        <div className="pt-3 border-t border-zinc-100">
          <p className="text-xs text-zinc-500 mb-2">
            Survey in progress. Surveyor will upload report and photos upon completion.
          </p>
        </div>
      )}

      {isCompleted && (
        <div className="pt-3 border-t border-zinc-100">
          <p className="text-xs text-teal-700 font-medium">
            Survey completed. Claim has progressed to Assessment stage.
          </p>
        </div>
      )}
    </Card>
  );
}