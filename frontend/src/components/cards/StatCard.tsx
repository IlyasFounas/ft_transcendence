interface StatCardProps {
    title: string;
    value: string;
    subtitle: string;
    icon: React.ReactNode;
}

export default function StatCard({
    title,
    value,
    subtitle,
    icon,
}: StatCardProps) {
    return (
        <div
            className="
                min-h-28
                bg-[linear-gradient(120deg,#071220_0%,#0B243A_57%,#07192E_100%)]
                border
                border-border-default
                rounded-md
                px-5
                py-2
                lg:p-4
                flex
                items-center
                lg:justify-center
                gap-10
            "
        >
            <div className="text-primary shrink-0">
                {icon}
            </div>

            <div className="uppercase">
                <p className="text-body-primary lg:text-heading-small text-text-primary">
                    {title}
                </p>
                <p className="text-heading-small lg:text-hero text-text-primary">
                    {value}
                </p>
                <p className="text-label-small lg:text-body-primary text-text-secondary">
                    {subtitle}
                </p>
            </div>

        </div>
    );
}