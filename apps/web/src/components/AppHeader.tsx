import { User } from '@phosphor-icons/react/User';
import { Bell } from '@phosphor-icons/react/Bell';
import { ArrowClockwise } from '@phosphor-icons/react/ArrowClockwise';
import { Claim } from '../types';
import { NotificationsPopover } from './NotificationsPopover';

interface AppHeaderProps {
  claim: Claim | null;
  notificationsOpen: boolean;
  onToggleNotifications: () => void;
  onSeed: () => void;
  onReset: () => void;
}

export function AppHeader({
  claim,
  notificationsOpen,
  onToggleNotifications,
  onSeed,
  onReset
}: AppHeaderProps) {
  return (
    <header className="bg-white border-b border-zinc-200 px-6 py-3 flex items-center justify-between shadow-xs">
      <div className="flex items-center gap-3">
        <span className="font-semibold text-lg tracking-tight">Taskmaster Claims Ops</span>
        <span className="text-xs bg-teal-50 text-teal-700 border border-teal-200 px-2 py-0.5 rounded font-medium">
          Claims Officer Console
        </span>
      </div>

      <div className="flex items-center gap-4">
        <div className="text-sm text-zinc-600 flex items-center gap-2">
          <User size={16} weight="bold" />
          <span>A. Rahman</span>
        </div>

        <div className="relative">
          <button
            onClick={onToggleNotifications}
            className="relative p-2 rounded-md hover:bg-zinc-100 text-zinc-700 transition focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2"
            aria-label="Notifications"
            aria-expanded={notificationsOpen}
            aria-controls="notifications-popover"
          >
            <Bell size={20} weight="regular" />
            {claim && claim.notifications.length > 0 && (
              <span className="absolute top-1 right-1 w-4 h-4 bg-red-600 text-white rounded-full text-[10px] flex items-center justify-center font-bold">
                {claim.notifications.length}
              </span>
            )}
          </button>

          <NotificationsPopover
            open={notificationsOpen}
            notifications={claim ? claim.notifications : []}
          />
        </div>

        <div className="flex items-center gap-2">
          {!claim ? (
            <button
              onClick={onSeed}
              className="bg-teal-700 hover:bg-teal-800 focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 text-white px-4 py-2 rounded-md text-sm font-medium transition cursor-pointer"
            >
              Seed demo Claim
            </button>
          ) : (
            <button
              onClick={onReset}
              className="border border-zinc-300 hover:bg-zinc-50 focus-visible:ring-2 focus-visible:ring-teal-700 focus-visible:ring-offset-2 px-3 py-2 rounded-md text-sm font-medium text-zinc-700 transition flex items-center gap-1.5 cursor-pointer"
            >
              <ArrowClockwise size={16} />
              Reset
            </button>
          )}
        </div>
      </div>
    </header>
  );
}
