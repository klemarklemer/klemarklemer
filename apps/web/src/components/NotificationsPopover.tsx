interface NotificationsPopoverProps {
  open: boolean;
  notifications: string[];
}

export function NotificationsPopover({ open, notifications }: NotificationsPopoverProps) {
  if (!open) return null;
  return (
    <div
      id="notifications-popover"
      role="dialog"
      aria-label="Notifications"
      className="absolute right-0 mt-2 w-80 bg-white border border-zinc-200 rounded-lg shadow-lg p-4 z-50"
    >
      <h3 className="font-semibold text-sm mb-2 text-zinc-800">Notifications</h3>
      <ul className="space-y-2 text-sm text-zinc-600">
        {notifications.map((notification, index) => (
          <li key={index} className="border-b border-zinc-100 pb-2 last:border-0 last:pb-0">
            {notification}
          </li>
        ))}
      </ul>
    </div>
  );
}
