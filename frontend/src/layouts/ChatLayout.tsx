import { useState } from "react";
import { Outlet } from "react-router-dom"
import ChatPanel from "../components/chat/ChatPanel";
import { Navbar } from "../components/navbar/Navbar";
import { chatlists } from "../data/chatlists";

export function ChatLayout () {

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
            
            < Outlet />
    
            <div
            className={`
                fixed
                right-5
                top-25
                bottom-5
                left-5
                z-50
                rounded-xl
                overflow-hidden
                border
                border-border-default
                bg-surface-base
    
                md:left-auto
                md:right-5
                md:bottom-5
                md:w-2/5
                lg:w-100
    
                ${isChatOpen ? "block" : "hidden"}
    
                lg:block
                `}
            >
                <ChatPanel friends={chatlists} />
            </div>
        </div>
    );
}