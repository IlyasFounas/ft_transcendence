import type { InputHTMLAttributes } from "react";

interface InputProps {
    label: string;
    type: InputHTMLAttributes<HTMLInputElement>["type"];
    value: string;
    placeholder?: string;
    onChange: (value: string) => void;
}

export default function Input ({
    label,
    type,
    value,
    placeholder,
    onChange,
}: InputProps) {

    return (
        <div className="flex flex-col gap-1 sm:gap-3 w-full">
            <label
                className="
                    text-body-primary
                    text-text-secondary
                    uppercase
                    sm:text-heading-small
                    md:text-heading-medium
                "
            >
                {label}
            </label>
            <input 
                type={type}
                value={value}
                placeholder={placeholder}
                onChange={(event) => onChange(event.target.value)}
                className="
                    w-full
                    h-11
                    md:h-15
                    px-4
                    sm:px-6
                    rounded-xl
                    bg-[linear-gradient(90deg,#07192E_0%,rgba(13,58,97,0.6)_38%,#07192E_100%)]
                    
                    text-text-primary
                    text-body-primary
                    sm:text-heading-small
                    md:text-heading-medium

                    border
                    border-border-default
                    transition
                    duration-150
                    focus:outline-none
                    focus:border-primary
                    focus:shadow-[0_0_8px_rgba(53,223,255,0.35)]
                    "
            />
        </div>
    );
}