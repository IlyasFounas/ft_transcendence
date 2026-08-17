import LobbyList from "./LobbyList";

const mockPlayers = [
  { id: "1", username: "Agent Tommy", role: "Field Analyst", status: "ready" as const },
  { id: "2", username: "Agent Nova", role: "Recon Specialist", status: "waiting" as const },
  { id: "3", username: "Agent Kaine", role: "Tech Support", status: "ready" as const },
  { id: "4", username: "Agent Vex", role: "Strategist", status: "ready" as const },
  { id: "5", username: "Agent Lyra", role: "Field Analyst", status: "waiting" as const },
  { id: "6", username: "Agent Lyra", role: "Field Analyst", status: "waiting" as const },
];

export default function LobbyTestPage() {
  return (
        <LobbyList numberOfPlayers={mockPlayers.length} listOfPlayers={mockPlayers} />
  );
}