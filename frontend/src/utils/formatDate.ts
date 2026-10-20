export function formatOperationDate(date: string): string {
    return new Date(date)
        .toLocaleString("en-GB", {
            day: "2-digit",
            month: "short",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        })
        .replace(",", "")
        .replace(":", ".")
        .toUpperCase();
}