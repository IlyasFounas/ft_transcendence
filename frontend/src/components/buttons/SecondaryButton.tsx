import type { BaseButtonProps } from "./button.types";

interface SecondaryButtonProps extends BaseButtonProps {
	icon?: React.ReactNode;
}


export default function SecondaryButton({ children, onClick, icon, disabled }: SecondaryButtonProps) {
  return (
	<button 
	 	disabled={disabled}
		onClick={onClick}
		className="
			inline-flex w-fit items-center justify-center
			
			px-[24px] py-[16px] gap-8 m-8
			
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

			enabled:hover:scale-[1.05]
			enabled:hover:brightness-110
 			enabled:hover:drop-shadow-[0_0_10px_rgba(40,215,255,0.5)]

			enabled:active:scale-[0.98]
			enabled:active:drop-shadow-none
			enabled:active:opacity-75

			disabled:border-border-default
			disabled:outline-border-default
    		disabled:text-text-secondary
    		disabled:shadow-[0px_0px_16px_0px_rgba(30,77,120,0.15)]
			disabled:cursor-default
		"
	>
		{icon}
		{children}
	</button>
	);
}
