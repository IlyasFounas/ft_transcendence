import { useState } from "react";
import { Outlet } from "react-router-dom"
import { Navbar } from "../components/navbar/Navbar";
import ChatPanel from "../components/chat/ChatPanel";
import Footer from "../components/routes/Footer";

export function MainLayout () {

    const [isChatOpen, setIsChatOpen] = useState(false);

    return (
        <>
        { isChatOpen && (
            <div
                className="
                    fixed
                    inset-0
                    z-30
                    bg-black/40
                    xl:hidden
                "
                onClick={() => setIsChatOpen(false)}
            />
        )}
        <div className="min-h-dvh flex flex-col">
            <div className="sticky top-0 z-60 w-full shrink-0">
                <Navbar
                    status="online"
                    isChatOpen={isChatOpen}
                    onToggleChat={() => setIsChatOpen((current) => !current)}
                    showChat={true}
                />
            </div>

			<Outlet />
			<Footer />

            <div
                className={`
                fixed
                top-28
                bottom-5
                left-5
                right-5
                z-50
                rounded-xl
                overflow-hidden
                border
                border-border-default
                bg-surface-base
                md:left-auto
                md:w-[min(26rem,calc(100vw-2.5rem))]
                ${isChatOpen ? "block" : "hidden"}
                `}
            >
                <ChatPanel />
            </div>
        </div>
    </>
);
}