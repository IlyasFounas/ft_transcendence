import { useEffect, useState, useRef } from "react";
import { useAuth } from "../context/AuthContext";

const DEFAULT_AVATAR = "avatars/avatar1.png";

interface Player {
  id: number;
  username: string;
  avatar: string;
  isHost?: boolean;
}

interface RawPlayer {
  id: number;
  username: string;
  avatar: string | null;
  host: boolean;
}

interface LobbyInfos {
	join_code : string;
	round_nb: number;
	timer: number;
}

export function useLobbySocket() {
	
	const { user, isAuthLoading } = useAuth();
	const [players, setPlayers] = useState<Player[]>([]);
	const [lobbyInfos, setLobbyInfos] = useState<LobbyInfos | null >(null);
	const [isConnected, setIsConnected] = useState(false);
	const socketRef = useRef<WebSocket | null>(null);
	const [isClosed, setIsClosed] = useState(false);
	const isCurrentUserHost = players.find((p) => p.id === user?.ID)?.isHost ?? false;

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
			const data = JSON.parse(event.data);

			if (Array.isArray(data))
			{
				const mappedPlayers: Player[] = data.map((p: RawPlayer) => ({
					id: p.id,
					username: p.username,
					avatar: p.avatar ?? DEFAULT_AVATAR,
					isHost: p.host,
				}));
				
				setPlayers(mappedPlayers);
			}
			else
			{
				setLobbyInfos(data);
			}
		}

		socket.onerror = (err) => {
		console.log("WS error:", err);
		};

		socket.onclose = () => {
			setIsConnected(false);
			setIsClosed(true);
		};

		return () => {
			socket.onclose = null;
			socket.close();
			socketRef.current = null;
		};
	}, [user, isAuthLoading]);

	function leaveLobby() {
		socketRef.current?.send(JSON.stringify({ Action: "leave", Payload: "" }));
	}

	function kickPlayer(targetId: number) {
	socketRef.current?.send(JSON.stringify({
		Action: "kick_player",
		Payload: { target_id: targetId },
	}));
	}

	return { players, isConnected, isClosed, isCurrentUserHost, leaveLobby, kickPlayer, lobbyInfos };
}