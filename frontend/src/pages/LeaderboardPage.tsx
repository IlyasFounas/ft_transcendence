import { Navbar } from "../components/navbar/Navbar";

export default function LeaderboardPage() {
  return (
    <div>
        <Navbar
            username="Agent Nova"
            avatar="/avatars/avatar1.png"
            status="online"
        />
        <h1>Leaderboard page</h1> 
    </div>
  );
}