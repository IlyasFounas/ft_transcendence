import { Navbar } from "../components/navbar/Navbar";

export default function HomePage() {
  return (
    <div>
      <Navbar
          username="Agent Nova"
          avatar="/avatars/avatar1.png"
          status="online"
      />
      <h1>Home page</h1> 
    </div>
  );
}