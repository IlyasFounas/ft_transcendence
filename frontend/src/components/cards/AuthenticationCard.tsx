interface AuthenticationCardProps {
    children: React.ReactNode;
}

export default function AuthenticationCard ({
    children,
}: AuthenticationCardProps) {
    return (
        <div className="
            relative
            flex
            flex-col
            min-w-xs
            min-h-140
            rounded-4xl
            bg-surface-base
            border
            border-border-default
            shadow-[inset_0_0_15px_0_#3586FF]
            
            px-7
            py-10
    
            sm:w-136
            sm:min-h-125

            md:pt-15
            "
        >
            <div className="
                flex
                flex-col
                items-center
                text-center
                gap-2
                "
            >
                <div className="text-heading-medium text-text-primary uppercase md:text-hero">
                    Welcome, Agent
                </div>
                <div className="text-body-primary text-text-secondary uppercase md:text-heading-small">
                    Sign in to begin your next mission
                </div>
            </div>
            <div className="
                flex
                flex-1
                items-center
                justify-center
                "
            >
                {children}
            </div>
        </div>
    );
}