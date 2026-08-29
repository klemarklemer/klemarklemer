import { Check } from '@phosphor-icons/react/Check';

interface ToastProps {
  message: string | null;
}

export function Toast({ message }: ToastProps) {
  if (!message) return null;
  return (
    <div className="fixed bottom-4 right-4 bg-zinc-900 text-white px-4 py-3 rounded-lg shadow-xl text-sm z-50 flex items-center gap-2" role="status">
      <Check size={18} className="text-teal-400" />
      <span>{message}</span>
    </div>
  );
}
