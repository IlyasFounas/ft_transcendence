export function formatOperationId(id: string): string {
    return `OP-${id.padStart(6, "0")}`;
}