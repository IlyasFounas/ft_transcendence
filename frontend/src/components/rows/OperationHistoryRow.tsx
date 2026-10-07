import { formatOperationId } from "../../utils/formatOperationId";

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
            xl:py-6
            xl:px-10
            xl:h-20
            bg-surface-base
            border
            border-border-default
            items-center
            text-text-secondary
            text-label-small
            md:text-heading-small
            xl:text-heading-medium
            "
        >
            <span className="text-left">{formatOperationId(operationId)}</span>
            <span className="text-center">{dateTime}</span>
            <span className="text-right">
                {points.toLocaleString()} <span>PTS</span>
            </span>
        </div>
    );
}