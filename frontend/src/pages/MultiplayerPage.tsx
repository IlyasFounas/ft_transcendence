import { Navbar } from "../components/navbar/Navbar"

export default function MultiplayerPage() {
  return (
    <div>
        <Navbar
            username="Agent Nova"
            avatar="/avatars/avatar1.png"
            status="online"
        />
        <h1>Multiplayer page</h1> 
    </div>
  );
}