export interface Message {
    id: number;
    variant: "received" | "sent";
    message: string;
    time: string;
}