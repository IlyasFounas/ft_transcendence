import { Outlet } from "react-router-dom"
import { Navbar } from "../components/navbar/Navbar";

export function MainLayout () {
    return (
        <div>
            <Navbar
                status="online"
                showChat={false}
            />
            
            < Outlet />
    
        </div>
    );
}