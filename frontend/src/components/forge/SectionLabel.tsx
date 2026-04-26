interface SectionLabelProps {
  children: React.ReactNode;
  className?: string;
}

export const SectionLabel = ({ children, className = "" }: SectionLabelProps) => (
  <div className={`flex items-center gap-3 ${className}`}>
    <span className="h-px w-8 bg-primary" />
    <span className="terminal-label">{children}</span>
  </div>
);
