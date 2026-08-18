import type { BaseButtonProps } from "./button.types";

type SmallButtonVariant = 
    | "primary" 
    | "danger";

interface SmallButtonProps extends BaseButtonProps {
    variant?: SmallButtonVariant;
}

export default function SmallButton({
    children,
    variant = "primary",
    disabled,
    onClick,
}: SmallButtonProps) {
    const variantClasses = {
        primary: `
            border-primary
            text-primary
            bg-[linear-gradient(to_right,#07192E_0%,#0D3A61_57%,#07192E_100%)]
            hover:brightness-125
            active:brightness-75
            disabled:hover:brightness-100
        `,

        danger: `
            border-status-red
            text-status-red
            bg-[#18070D]
            hover:bg-[#35101A]
            active:bg-[#100409]
            disabled:hover:bg-[#18070D]
        `,
    };

    return (
        <button
            onClick={onClick}
            disabled={disabled}
            className={`
            px-2
            py-1
            rounded-md
            border-2
            uppercase
            text-body-primary
            transition
            cursor-pointer
            duration-150
            disabled:opacity-40
            disabled:cursor-default
            ${variantClasses[variant]}
        `}
        >
            {children}
        </button>
    );
}