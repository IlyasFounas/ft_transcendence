import { useEffect, useState, useRef } from "react";
import { useAuth } from "../context/AuthContext";
import type { StatusBadgeProps } from "../components/primitives/StatusBadge";

const DEFAULT_AVATAR = "avatars/avatar1.png";

interface Player {
  id: number;
  username: string;
  role: string;
  status: StatusBadgeProps["status"];
  avatar: string;
  isHost?: boolean;
}

interface RawPlayer {
  id: number;
  username: string;
  avatar: string | null;
  isHost: boolean;
}

export function useLobbySocket() {
	
	const { user, isAuthLoading } = useAuth();
	const [players, setPlayers] = useState<Player[]>([]);
	const [isConnected, setIsConnected] = useState(false);
	const socketRef = useRef<WebSocket | null>(null);

	useEffect(() => {
		if (isAuthLoading || !user) return;

		const socket = new WebSocket(
		`ws://localhost:5130/api/ws/lobby/join?user_id=${user.ID}`
		);
		socketRef.current = socket;

		socket.onopen = () => {
		setIsConnected(true);
		};

		socket.onmessage = (event) => {
		const rawPlayers: RawPlayer[] = JSON.parse(event.data);

		const mappedPlayers: Player[] = rawPlayers.map((p) => ({
			id: p.id,
			username: p.username,
			role: "",
			status: "waiting",
			avatar: p.avatar ?? DEFAULT_AVATAR,
		}));

		setPlayers(mappedPlayers);
		};

		socket.onerror = (err) => {
		console.log("WS error:", err);
		};

		socket.onclose = () => {
		setIsConnected(false);
		};

		return () => {
		socket.close();
		socketRef.current = null;
		};
	}, [user, isAuthLoading]);

	function leaveLobby() {
		socketRef.current?.send(JSON.stringify({ Action: "leave", Payload: "" }));
	}

	return { players, isConnected, leaveLobby };
}