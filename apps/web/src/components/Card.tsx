import { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
}

export function Card({ children, className = '' }: CardProps) {
  return <div className={`bg-white border border-zinc-200 rounded-lg p-5 shadow-xs ${className}`}>{children}</div>;
}
