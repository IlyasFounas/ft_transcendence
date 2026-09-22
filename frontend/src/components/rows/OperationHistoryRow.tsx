interface OperationHistoryRowProps {
    operationId: string;
    dateTime: string;
    points: number;
}

export default function OperationHistoryRow({
    operationId,
    dateTime,
    points,
}: OperationHistoryRowProps) {
    return (
        <div className="
            grid
            grid-cols-3
            h-16
            px-4
            bg-surface-base
            border
            border-border-default
            items-center
            text-text-secondary
            text-label-small
            lg:text-heading-small
            "
        >
            <span className="text-left">{operationId}</span>
            <span className="text-center">{dateTime}</span>
            <span className="text-right">
                {points.toLocaleString()} <span>PTS</span>
            </span>
        </div>
    );
}