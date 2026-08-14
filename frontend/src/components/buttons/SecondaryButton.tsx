import type { BaseButtonProps } from "./button.types";

interface SecondaryButtonProps extends BaseButtonProps {
	active?: boolean;
	icon?: React.ReactNode;
}


export default function SecondaryButton({ children, onClick, icon, active }: SecondaryButtonProps) {
  return (
	<button 
		onClick={onClick}
		className={`
			inline-flex w-fit h-[60px] items-center justify-center
			
			px-[24px] py-[16px] gap-3 m-8
			
			rounded-lg
			border border-primary
			outline outline-1 outline-offset-[-1px] outline-primary
			
			bg-surface-elevated
			
			text-heading-medium
			
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
			
			${active
        ? `
			text-text-primary
            border-border-active
            outline-border-active
            shadow-[0px_0px_16px_0px_rgba(40,215,255,0.5)]
        `
        : `
            border-border-default
            outline-border-default
            text-text-secondary
            shadow-[0px_0px_16px_0px_rgba(30,77,120,0.15)]
        `}
			`}
	>
		{icon}
		{children}
	</button>
	);
}
