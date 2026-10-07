import NavbarGameplay from '../components/navbar/NavbarGameplay'
import { Outlet } from "react-router-dom"


export default function GameLayout()
{
	return (
		 <div>
            <NavbarGameplay/>
            
            < Outlet />
    
        </div>
	)
}