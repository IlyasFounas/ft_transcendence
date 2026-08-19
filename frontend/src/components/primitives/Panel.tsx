type PanelProps = {
  title: string;
  children: React.ReactNode;
  className?: string;
};

export default function Panel({ title, children, className = "" }: PanelProps) {
  return (
    <div
      className={`rounded-[12px] border border-border-default/30
                  shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden ${className}`}
    >
      	<div className="flex items-center justify-between border-b border-border-default/30 
                      bg-surface-title px-4 py-3">
        <span className="text-heading-small text-text-primary/70 uppercase">
          {title}
        </span>
      </div>

      <div className="p-4 bg-[#020A18]">{children}</div>
    </div>
  );
}