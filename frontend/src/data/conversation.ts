import type { Message } from "../types/message.type";

export const conversations: Record<string, Message[]> = {
  "Agent Elle": [
    {
      id: 1,
      variant: "received",
      message: "Hey Nova, are you free?",
      time: "03:00 PM",
    },
    {
      id: 2,
      variant: "sent",
      message: "Yeah, what's up?",
      time: "03:01 PM",
    },
    {
      id: 3,
      variant: "received",
      message: "Great do you wanna join the game?",
      time: "03:02 PM",
    },
  ],

  "Agent Louis": [
    {
      id: 4,
      variant: "received",
      message: "Are you heading out?",
      time: "03:00 PM",
    },
    {
      id: 5,
      variant: "sent",
      message: "Yeah, I'll be there soon.",
      time: "03:01 PM",
    },
    {
      id: 6,
      variant: "received",
      message: "Bye",
      time: "03:01 PM",
    },
  ],

  "Agent Mike": [
    {
      id: 7,
      variant: "sent",
      message: "Are you joining us?",
      time: "03:00 PM",
    },
    {
      id: 8,
      variant: "received",
      message: "Yeah, I'm getting ready now.",
      time: "03:01 PM",
    },
    {
      id: 9,
      variant: "received",
      message: "Be there in 5",
      time: "03:01 PM",
    },
  ],

  "Agent Dean": [
    {
      id: 10,
      variant: "received",
      message: "We're starting the mission soon.",
      time: "02:52 PM",
    },
    {
      id: 11,
      variant: "sent",
      message: "Give me a couple of minutes.",
      time: "02:53 PM",
    },
    {
      id: 12,
      variant: "received",
      message: "Now? I'm ready to join",
      time: "02:54 PM",
    },
  ],

  "Agent Mark": [
    {
      id: 13,
      variant: "sent",
      message: "Hey Mark, did you get my message?",
      time: "02:08 PM",
    },
    {
      id: 14,
      variant: "received",
      message: "Yeah, I just saw it.",
      time: "02:09 PM",
    },
    {
      id: 15,
      variant: "received",
      message: "What is it?",
      time: "02:10 PM",
    },
  ],
};