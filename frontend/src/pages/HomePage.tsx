import { Navbar } from "../components/navbar/Navbar";
import { useState } from "react";
import ChatPanel from "../components/chat/ChatPanel";
import type { Friend } from "../types/user.types";

const friends: Friend[] = [
    {
        username: "Agent Elle",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "online",
        lastMessage: "Great do you wanna join the game?",
        timestamp: "03:02 PM",
    },
    {
        username: "Agent Louis",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "online",
        lastMessage: "Bye",
        timestamp: "03:01 PM",
    },
    {
        username: "Agent Mike",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "online",
        lastMessage: "Be there in 5",
        timestamp: "03:01 PM",
    },
    {
        username: "Agent Dean",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "in-mission",
        lastMessage: "Now? I'm ready to join",
        timestamp: "02:54 PM",
    },
    {
        username: "Agent Mark",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "offline",
        lastMessage: "What is it?",
        timestamp: "02:10 PM",
    },
    {
        username: "Agent Louis",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "online",
        lastMessage: "Bye",
        timestamp: "03:01 PM",
    },
    {
        username: "Agent Mike",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "online",
        lastMessage: "Be there in 5",
        timestamp: "03:01 PM",
    },
    {
        username: "Agent Dean",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "in-mission",
        lastMessage: "Now? I'm ready to join",
        timestamp: "02:54 PM",
    },
    {
        username: "Agent Mark",
        avatar: "/avatars/avatar1.png",
        role: "Field Analyst",
        status: "offline",
        lastMessage: "What is it?",
        timestamp: "02:10 PM",
    },
];

export default function HomePage() {

  const [isChatOpen, setIsChatOpen] = useState(false);

  return (
    <div>
      <Navbar
          username="Agent Nova"
          avatar="/avatars/avatar1.png"
          status="online"
          isChatOpen={isChatOpen}
          onToggleChat={() => setIsChatOpen((current) => !current)}
      />

      <div
        className={`
            group
            fixed
            right-5
            top-25
            bottom-5
            z-50
            w-28
            rounded-xl
            overflow-hidden
            border
            border-border-default
            bg-surface-base
            transition-[width]
            duration-200
            ease-in-out
            hover:w-80

            ${isChatOpen ? "block" : "hidden"}

            lg:block
            `}
        >
            <ChatPanel friends={friends} />
        </div>
    </div>
  );
}
