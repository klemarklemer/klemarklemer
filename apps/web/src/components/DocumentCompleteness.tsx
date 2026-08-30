import { FileArrowUp } from '@phosphor-icons/react';
import { Claim } from '../types';
import { Card } from './Card';

interface DocumentCompletenessProps {
  claim: Claim;
  isComplete: boolean;
  onUpload: () => void;
  busy?: boolean;
}

export function DocumentCompleteness({ claim, isComplete, onUpload, busy }: DocumentCompletenessProps) {
  return (
    <Card>
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-semibold text-sm text-zinc-800 uppercase tracking-wide">Document completeness</h3>
        <span
          className={`text-xs px-2 py-0.5 rounded font-medium ${
            isComplete
              ? 'bg-teal-50 text-teal-700 border border-teal-200'
              : 'bg-red-50 text-red-700 border border-red-200'
          }`}
        >
          {isComplete ? 'Complete' : 'Incomplete'}
        </span>
      </div>

      <ul className="space-y-2 mb-4">
        {claim.documents.map((document) => (
          <li
            key={document.id}
            className="flex items-center justify-between text-sm py-1.5 border-b border-zinc-100 last:border-0"
          >
            <span className="text-zinc-700 flex items-center gap-2">
              <span className={`w-2 h-2 rounded-full ${document.present ? 'bg-teal-600' : 'bg-red-500'}`} />
              {document.name}
            </span>
            <span className={`text-xs font-medium ${document.present ? 'text-teal-700' : 'text-red-600'}`}>
              {document.present ? 'Present' : 'Missing'}
            </span>
          </li>
        ))}
      </ul>

      {!isComplete && !claim.closed && (
        <div className="pt-2">
          <label
            htmlFor="police-upload-input"
            className={`inline-flex items-center gap-2 bg-teal-700 hover:bg-teal-800 focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 text-white px-4 py-2 rounded-md text-sm font-medium transition cursor-pointer ${busy ? 'opacity-60 cursor-not-allowed' : ''}`}
          >
            <FileArrowUp size={16} />
            Upload police report
          </label>
          <input
            id="police-upload-input"
            type="file"
            className="sr-only"
            onChange={onUpload}
            disabled={busy}
            aria-label="Upload police report"
          />
          <p className="text-xs text-zinc-500 mt-1.5">
            Uploading marks the police report as present and refreshes the Claim.
          </p>
        </div>
      )}
    </Card>
  );
}
