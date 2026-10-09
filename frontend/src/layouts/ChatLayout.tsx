import { useState } from "react";
import { Outlet } from "react-router-dom"
import ChatPanel from "../components/chat/ChatPanel";
import { Navbar } from "../components/navbar/Navbar";
import { chatlists } from "../data/chatlists";
import Footer from "../components/routes/Footer";

export function ChatLayout () {

    const [isChatOpen, setIsChatOpen] = useState(false);
    
        return (
            <>
            { isChatOpen && (
                <div
                    className="
                        fixed
                        inset-0
                        z-30
                        bg-black/20
                        lg:hidden
                    "
                    onClick={() => setIsChatOpen(false)}
                />
            )}
            <div className="h-dvh flex flex-col overflow-hidden">
                <Navbar
                    status="online"
                    isChatOpen={isChatOpen}
                    onToggleChat={() => setIsChatOpen((current) => !current)}
                    showChat={true}
                />
                
                < Outlet />
				<Footer />
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
                    md:bottom-10
                    md:w-2/5
                    lg:w-100
        
                    ${isChatOpen ? "block" : "hidden"}
        
                    lg:block
                    `}
                >
                    <ChatPanel friends={chatlists} />
                </div>
            </div>
        </>
    );
}