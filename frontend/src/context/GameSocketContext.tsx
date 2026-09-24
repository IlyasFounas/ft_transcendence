import { createContext, useContext, useEffect, useState, useRef} from "react"
import type { ReactNode } from "react"
import { useAuth } from "./AuthContext"
import type { RoundResult } from "../types/game";

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
	join_code: string;
	round_total: number;
	timer: number;
}

interface RoundInfo {
	round_nb: number;
	mid: number;
	round_ends_at: string;
}

interface LocationInfo {
	lat: number;
	long: number;
	country_code: string;
}

interface RoundInfos {
	distance: number;
	score: number;
	base_score: number;
	country_found: boolean;
	base: LocationInfo;
	guess: LocationInfo;
	streak: number;
	time_left_seconds: number;
}

interface PlayerRoundResult {
	user_id: number;
	username: string;
	round_nb: number;
	has_guessed: boolean;
	round_infos: RoundInfos | null;
	total_score: number;
}

interface RoundResultBroadcast {
	results: PlayerRoundResult[];
	is_last_round: boolean;
	next_round_at: string | null;
}

interface GameSocketContextType {
	players: Player[];
	lobbyInfos: LobbyInfos | null;
	isConnected: boolean;
	isClosed: boolean;
	isCurrentUserHost: boolean;
	leaveLobby: () => void;
	kickPlayer: (targetId: number) => void;
	sendGuess: (guess: { lat: number; lng: number }) => void;
	startGame: () => void;
	round: RoundInfo | null;
	roundHistory: RoundResult[];
	lastRoundResult: {
		isLastRound: boolean;
		nextRoundAt: string | null;
		actualLocation: { latitude: number; longitude: number };
	} | null;
}

const GameSocketContext = createContext<GameSocketContextType>({
	players: [],
	lobbyInfos: null,
	isConnected: false,
	isClosed: false,
	isCurrentUserHost: false,
	leaveLobby: () => {},
	kickPlayer: () => {},
	sendGuess: () => {},
	startGame: () => {},
	round: null,
	roundHistory: [],
	lastRoundResult: null,
});

export function GameSocketProvider({ children }: { children: ReactNode }) {
	const { user, isAuthLoading } = useAuth();
	const [players, setPlayers] = useState<Player[]>([]);
	const [lobbyInfos, setLobbyInfos] = useState<LobbyInfos | null>(null);
	const [isConnected, setIsConnected] = useState(false);
	const socketRef = useRef<WebSocket | null>(null);
	const [isClosed, setIsClosed] = useState(false);
	const [round, setRound] = useState<RoundInfo | null>(null);
	const [roundHistory, setRoundHistory] = useState<RoundResult[]>([]);
	const [lastRoundResult, setLastRoundResult] = useState<{
		isLastRound: boolean;
		nextRoundAt: string | null;
		actualLocation: { latitude: number; longitude: number };
	} | null>(null);

	const isCurrentUserHost = players.find((p) => p.id == user?.ID)?.isHost ?? false;

	useEffect(() => {
		if (isAuthLoading || !user) return;

		const socket = new WebSocket(
		`ws://localhost:5130/api/ws/lobby/join?user_id=${user.ID}`
		);
		socketRef.current = socket;

		socket.onopen = () => setIsConnected(true);

		socket.onmessage = (event) => {
		const message = JSON.parse(event.data);

		if (message.type === "lobby_players") {
			const mappedPlayers: Player[] = message.data.map((p: RawPlayer) => ({
			id: p.id,
			username: p.username,
			avatar: p.avatar ?? DEFAULT_AVATAR,
			isHost: p.host,
			}));
			setPlayers(mappedPlayers);
		}
		else if (message.type === "round_started") {
			const data = message.data as { round_nb: number; mid: number; round_ends_at: string };
			setRound({
			round_nb: data.round_nb,
			mid: data.mid,
			round_ends_at: data.round_ends_at,
			});
		}
		else if (message.type === "round_results") {
			const data = message.data as RoundResultBroadcast;
			const own = data.results.find((r) => r.user_id === user.ID);
			
			if (own && own.round_infos) {
			setRoundHistory((prev) => [
				...prev,
				{
				round: own.round_nb,
				points: own.round_infos!.score,
				distance: own.round_infos!.distance,
				guessLocation: {
					latitude: own.round_infos!.guess.lat,
					longitude: own.round_infos!.guess.long,
				},
				},
			]);

			setLastRoundResult({
				isLastRound: data.is_last_round,
				nextRoundAt: data.next_round_at,
				actualLocation: {
				latitude: own.round_infos.base.lat,
				longitude: own.round_infos.base.long,
				},
			});
			}
		}
		else if (message.join_code) {
			setLobbyInfos(message);
		}
		};

		socket.onerror = (err) => console.log("WS error:", err);

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

	function sendGuess(guess: { lat: number; lng: number }) {
		if (!round) return;
		socketRef.current?.send(JSON.stringify({
		Action: "guess",
		Payload: {
			latitude: guess.lat,
			longitude: guess.lng,
			round_nb: round.round_nb,
		},
		}));
	}

	function kickPlayer(targetId: number) {
		socketRef.current?.send(JSON.stringify({
		Action: "kick_player",
		Payload: { target_id: targetId },
		}));
	}

	function startGame() {
		socketRef.current?.send(JSON.stringify({
		Action: "start_round",
		Payload: "",
		}));
	}

	return (
		<GameSocketContext.Provider
		value={{ players, lobbyInfos, isConnected, isClosed, isCurrentUserHost, leaveLobby, kickPlayer, startGame, round, roundHistory, lastRoundResult, sendGuess }}
		>
		{children}
		</GameSocketContext.Provider>
	);
}

export function useLobbySocket() {
  	return useContext(GameSocketContext);
}