import type { Message } from "../types/message.type"

export const conversations: Record<string, Message[]> = {
    "Agent Elle": [
        {
            id: 1,
            variant: "received",
            message: "Hello! How are you doing?",
            time: "03:01 PM",
        },
        {
            id: 2,
            variant: "sent",
            message: "I'm good, how are you?",
            time: "03:02 PM",
        },
        {
            id: 3,
            variant: "received",
            message: "Great! Do you wanna join game now?",
            time: "03:02 PM",
        },
    ],

    "Agent Louis": [
        {
            id: 4,
            variant: "received",
            message: "Hey, are you online?",
            time: "02:30 PM",
        },
        {
            id: 5,
            variant: "sent",
            message: "Yeah, what's up?",
            time: "02:31 PM",
        },
    ],
};