import { Outlet } from "react-router-dom"
import { Navbar } from "../components/navbar/Navbar";

export function MainLayout () {
    return (
        <div>
            <Navbar
                username="Agent Nova"
                avatar="/avatars/avatar1.png"
                status="online"
                showChat={false}
            />
            
            < Outlet />
    
        </div>
    );
}