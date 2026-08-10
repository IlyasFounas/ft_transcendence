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
        <div className="grid grid-cols-3 h-[64px] px-4 bg-surface-base border border-border-default items-center">
            <span className="text-secondary text-heading-small text-left">{operationId}</span>
            <span className="text-secondary text-heading-small text-center">{dateTime}</span>
            <span className="text-secondary text-heading-small text-right">{points}<span> PTS</span></span>
        </div>
    );
}