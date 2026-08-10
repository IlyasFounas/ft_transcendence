import type { BaseButtonProps } from "./button.types";
import { SmallButtonSvg, MediumButtonSvg, LargeButtonSvg } from "./ButtonShape";
interface PrimaryButtonProps extends BaseButtonProps {
	subtitle?: string;
	size?: "small" | "medium" | "large";
}

const sizeStyles = {
  small: {
    container: "w-[254px] h-[98px]",
    title: "text-heading-small",
    subtitle: "text-caption-small",
  },
  medium: {
    container: "w-[394px] h-[124px]",
    title: {
      default: "text-heading-large",
      withSubtitle: "text-heading-medium",
    },
    subtitle: "text-caption-medium",
  },
  large: {
    container: "w-[534px] h-[159px]",
    title: "text-hero",
    subtitle: "text-caption-large",
  },
};

const buttonShapes = {
  large: LargeButtonSvg,
  medium: MediumButtonSvg,
  small: SmallButtonSvg,
};

export default function PrimaryButton({ children, onClick, subtitle, disabled, size = "large"}: PrimaryButtonProps) {
	
	const styles = sizeStyles[size];

	const Shape = buttonShapes[size];

	const titleClass = typeof styles.title === "object" ? subtitle ? styles.title.withSubtitle : styles.title.default : styles.title;

  	return (
		<button
		disabled = {disabled}
		onClick = {onClick}
		className={`
			group
			relative

			cursor-pointer
			transition-transform
			duration-200
			ease-out
			enabled:hover:scale-[1.02]
			enabled:active:scale-[0.98]

			enabled:active:[--button-gradient-start:#07192E]
			enabled:active:[--button-gradient-middle:#010101]
			enabled:active:[--button-gradient-end:#07192E]

			disabled:cursor-default
    		disabled:opacity-40
			${styles.container}
		`}
		>
  		<Shape className="
			absolute inset-0 
			w-full h-full
			transition-all
			duration-200
			ease-out
			group-enabled:group-hover:brightness-125
			group-enabled:group-hover:drop-shadow-[0_0_20px_rgba(40,215,255,0.5)]
			group-enabled:group-active:brightness-100
			group-enabled:group-active:drop-shadow-none

			group-disabled:brightness-100
			group-disabled:drop-shadow-none"
  		/>
		<div className=" relative flex h-full flex-col items-center justify-center">
			<span className={`text-primary uppercase ${titleClass}`}>
			{children}
			</span>
			{subtitle && (
			<span className={`text-text-secondary uppercase ${styles.subtitle}`}>
				{subtitle}
			</span>
			)}
		</div>
		</button>
	);
}
