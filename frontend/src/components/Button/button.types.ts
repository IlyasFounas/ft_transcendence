export interface BaseButtonProps {
	children: React.ReactNode;
	onClick?: () => void;
	disabled?: boolean;
	state?: "default" | "hover" | "clicked" | "disabled";
}