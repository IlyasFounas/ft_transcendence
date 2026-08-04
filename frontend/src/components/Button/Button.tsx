type ButtonProps = {
	children: React.ReactNode;
	icon?: React.ReactNode;
	onClick?: () => void;
	disabled?: boolean;
}

export default function Button({ children, onClick, icon, disabled }: ButtonProps) {
  return (
	<button 
	 	disabled={disabled}
		onClick={onClick}
		className="
			inline-flex items-center justify-center
			
			px-24 py-16 gap-8 m-8
			
			rounded-lg
			border border-primary
			outline outline-1 outline-offset-[-1px] outline-primary
			
			bg-surface-elevated
			
			text-heading-medium

			shadow-[0_0_6px_rgba(53,223,255,0.5)]

			uppercase

			cursor-pointer

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
