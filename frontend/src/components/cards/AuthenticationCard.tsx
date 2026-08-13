interface AuthenticationCardProps {
    children: React.ReactNode;
}

export default function AuthenticationCard ({
    children,
}: AuthenticationCardProps) {
    return (
        <div className="
            w-[620px]
            h-[750px]
            rounded-4xl
            bg-surface-base
            border
            border-border-default
            shadow-[inset_0_0_15px_0_#3586FF]
            px-[80px]
            py-[60px]
            "
        >
            <div className="flex flex-col items-center text-center gap-[10px]">
                <div className="text-hero text-text-primary uppercase">Welcome, Agent</div>
                <div className="text-heading-small text-text-secondary uppercase">Sign in to begin your next mission</div>
            </div>
            {children}
        </div>
    );
}