import OperationHistoryRow from "../rows/OperationHistoryRow";
import type { Operation } from "../../types/operation.types";

interface OperationHistoryProps {
    operations: Operation[];
}
function TableHeader() {
    return (
        <div className="
            grid
            grid-cols-3
            bg-surface-elevated
            text-label-small
            md:text-heading-small
            xl:text-heading-medium
            items-center
            text-text-primary
            uppercase
            px-4
            py-4
            xl:py-6
            xl:px-10
            "
        >
            <span className="text-left">operation id</span>
            <span className="text-center">date & time</span>
            <span className="text-right">xp earned</span>
        </div>
    );
}

export default function OperationHistory({
    operations,
}: OperationHistoryProps) {
    return (
        <div
            className="
                bg-surface-base
                p-4
                xl:p-10
                border
                border-border-default
                rounded-lg

                xl:h-full
                xl:min-h-0
                xl:flex
                xl:flex-col
            "
        >
            <h2 className="text-heading-small lg:text-heading-medium xl:text-heading-large text-center uppercase mb-4">
                Operation History
            </h2>
            <TableHeader />
            <div
                className="
                    max-h-64
                    xl:max-h-none
                    xl:flex-1
                    xl:min-h-0
                    overflow-y-auto
                    custom-scrollbar
                "
            >
                {operations.map((operation) => (
                    <OperationHistoryRow
                        key={operation.operationId}
                        operationId={operation.operationId}
                        dateTime={operation.dateTime}
                        points={operation.points}
                    />
                ))}
            </div>
        </div>
    );
}