import type { BaseButtonProps } from "./button.types";

interface SecondaryButtonProps extends BaseButtonProps {
	active?: boolean;
	icon?: React.ReactNode;
	className?: string;
	variant?: "classic" | "danger"
}

export default function SecondaryButton({ children, onClick, icon, active, className, variant="classic", }: SecondaryButtonProps) {
  return (
	<button 
		onClick={onClick}
		className={`
			inline-flex w-fit min-h-15 items-center justify-center
			
			px-6 py-4 gap-3
			
			rounded-lg
			border border-primary
			outline -outline-offset-1 outline-primary
			
			bg-surface-elevated
			
			lg:text-heading-medium
			text-body-primary

			shadow-[0_0_6px_rgba(53,223,255,0.5)]
			
			uppercase
			
			cursor-pointer
			transition-transform
			duration-200
			ease-out
			
			hover:scale-[1.05]
			hover:brightness-110
			hover:drop-shadow-[0_0_10px_rgba(40,215,255,0.5)]

			active:scale-[0.98]
			active:drop-shadow-none
			active:opacity-75
			
			${variant === "danger"
  ? `
            hover:border-red-500
            hover:outline-red-500
            hover:text-red-500
            hover:drop-shadow-[0_0_10px_rgba(255,59,59,0.6)]
        `
  : `
            hover:brightness-110
            hover:drop-shadow-[0_0_10px_rgba(40,215,255,0.5)]
        `}
			${active
		? `
			text-text-primary
			border-border-active
			outline-primary
			shadow-[0px_0px_16px_0px_rgba(40,215,255,0.5)]
		`
		: `
			border-surface-elevated
			outline-border-default
			text-text-secondary
			shadow-[0px_0px_16px_0px_rgba(30,77,120,0.15)]
		`}
			${className}
			`}
	>
		{icon}
		{children}
	</button>
	);
}
