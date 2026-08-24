import { Navbar } from "../components/navbar/Navbar"

export default function ProfilePage() {
  return (
    <div>
        <Navbar
            username="Agent Nova"
            avatar="/avatars/avatar1.png"
            status="online"
        />
        <h1>Profile page</h1> 
    </div>
  );
}