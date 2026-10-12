import { createContext, useContext, useEffect, useState, useRef} from "react"
import type { ReactNode } from "react"
import { useAuth } from "./AuthContext"
import type { RoundResult, PlayerRoundResult } from "../types/game";

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

interface RoundResultBroadcast {
	round_nb: number;
	base: { lat: number; long: number; country_code: string }; 
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
	leaveGame: () => void;
	kickPlayer: (targetId: number) => void;
	sendGuess: (guess: { lat: number; lng: number }) => void;
	startGame: () => void;
	startNextRound: () => void;
	connect: () => void;
	resetGame: () => void;
	round: RoundInfo | null;
	roundHistory: RoundResult[];
	lastRoundResult: {
		isLastRound: boolean;
		nextRoundAt: string | null;
		actualLocation: { latitude: number; longitude: number };
		allResults: PlayerRoundResult[];
		fullHistory?: any[]; 
	} | null;
	gameStartAt: string | null;
}

const GameSocketContext = createContext<GameSocketContextType>({
	players: [],
	lobbyInfos: null,
	isConnected: false,
	isClosed: false,
	isCurrentUserHost: false,
	leaveLobby: () => {},
	leaveGame: () => {},
	kickPlayer: () => {},
	sendGuess: () => {},
	startGame: () => {},
	startNextRound: () => {},
	connect: () => {},
	resetGame: () => {},
	round: null,
	roundHistory: [],
	lastRoundResult: null,
	gameStartAt: null,
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
		allResults: PlayerRoundResult[];
		fullHistory?: any[];
	} | null>(null);
	const [gameStartAt, setGameStartAt] = useState<string | null>(null);

	const isCurrentUserHost = players.find((p) => p.id == user?.ID)?.isHost ?? false;
	function connect(){
		if (isAuthLoading || !user) return;
		if (socketRef.current) return;

		// const socket = new WebSocket(
		// 	// `${window.location.protocol}://${window.location.host}/api/ws/lobby/join?user_id=${user.ID}`
		// 	`ws://campus.ocefunding.com:5130/api/ws/lobby/join?user_id=${user.ID}`
		// );
		const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
		const socket = new WebSocket(
			`${proto}//${window.location.host}/api/ws/lobby/join?user_id=${user.ID}`
		);
		socketRef.current = socket;
		setIsClosed(false);

		socket.onopen = () => setIsConnected(true);

		socket.onmessage = (event) => {
			const message = JSON.parse(event.data);
			if (message.action === "lobby_players") {
				const mappedPlayers: Player[] = message.payload.map((p: RawPlayer) => ({
					id: p.id,
					username: p.username,
					avatar: p.avatar ?? DEFAULT_AVATAR,
					isHost: p.host,
				}));
				setPlayers(mappedPlayers);
			}
			else if (message.action === "lobby_infos") {
				setLobbyInfos(message.payload);
			}
			else if (message.action === "game_started") {
				setGameStartAt(message.payload.next_round_at);
			}
			else if (message.action === "round_started") {
				const data = message.payload as { round_nb: number; mid: number; round_ends_at: string };
				setRound({
					round_nb: data.round_nb,
					mid: data.mid,
					round_ends_at: data.round_ends_at,
				});
				setLastRoundResult(null);
				setGameStartAt(null);
			}
			else if (message.action === "round_results") {
				const data = message.payload as RoundResultBroadcast;
				const own = data.results.find((r) => r.user_id === user.ID);
				
				if (own && own.has_guessed && own.round_infos) {
					setRoundHistory((prev) => [
						...prev,
						{
							round: data.round_nb,
							points: own.round_infos.score,
							distance: own.round_infos.distance,
							guessLocation: {
								latitude: own.round_infos.guess.lat,
								longitude: own.round_infos.guess.long,
							},
						},
					]);
				} else {
					setRoundHistory((prev) => [
						...prev,
						{
							round: data.round_nb,
							points: 0,
							distance: 0,
							guessLocation: { latitude: 0, longitude: 0 }, 
						},
					]);
				}
				setLastRoundResult({
					isLastRound: data.is_last_round,
					nextRoundAt: data.next_round_at,
					actualLocation: {
						latitude: data.base.lat,
						longitude: data.base.long,
					},
					allResults: data.results,
				});
			}
			else if (message.action === "error") {
				if (message.payload?.message === "the game has already started") {
					socket.send(JSON.stringify({ action: "start_round", payload: "" }));
				}
			}
			else if (message.action === "game_results") {
				const rounds = message.payload.rounds;
				
				const lastRoundData = rounds.reduce((prev, current) => 
					(prev.round_nb > current.round_nb) ? prev : current
				);
				const finalRoundHistory = rounds.map((r: any) => {
					const own = r.results.find((res: any) => res.user_id === user.ID);
					return {
						round: r.round_nb,
						points: own?.round_infos?.score ?? 0,
						distance: own?.round_infos?.distance ?? 0,
						guessLocation: own?.has_guessed && own.round_infos ? {
							latitude: own.round_infos.guess.lat,
							longitude: own.round_infos.guess.long,
						} : { latitude: 0, longitude: 0 }
					};
				});
				setRoundHistory(finalRoundHistory);
			   
				setLastRoundResult({
					isLastRound: true,
					nextRoundAt: null,
					actualLocation: {
						latitude: lastRoundData.base.lat,
						longitude: lastRoundData.base.long,
					},
					allResults: lastRoundData.results,
					fullHistory: rounds,
				}
			);
		}	
	}

		socket.onerror = (err) => console.log("WS error:", err);

		socket.onclose = () => {
			setIsConnected(false);
			setIsClosed(true);
			socketRef.current = null;
		};
	}

	useEffect(() => {
		return () => {
			if (socketRef.current) {
				socketRef.current.onclose = null;
				socketRef.current.close();
				socketRef.current = null;
			}
		};
	}, [user]);

	function leaveLobby() {
		socketRef.current?.send(JSON.stringify({ action: "leave", payload: "" }));
	}

	function leaveGame() {
		socketRef.current?.send(JSON.stringify({ action: "leave", payload: "" }));
}   

	function sendGuess(guess: { lat: number; lng: number }) {
		
		if (!round) {
			return;
		}

		if (typeof guess.lat !== 'number' || typeof guess.lng !== 'number' || isNaN(guess.lat) || isNaN(guess.lng)) {
			return;
		}

		socketRef.current?.send(JSON.stringify({
			action: "guess",
			payload: {
				latitude: guess.lat,
				longitude: guess.lng,
				round_nb: round.round_nb,
			},
		}));
	}

	function kickPlayer(targetId: number) {
		socketRef.current?.send(JSON.stringify({
			action: "kick_player",
			payload: { target_id: targetId },
		}));
	}

	function startGame() {
		socketRef.current?.send(JSON.stringify({
			action: "start_game",
			payload: "",
		}));
	}

	function startNextRound() {
		setRound(null);
		setLastRoundResult(null);

		socketRef.current?.send(JSON.stringify({
			action: "start_round",
			payload: "",
		}));
	}

	function resetGame() {
		if (socketRef.current) {
			socketRef.current.onclose = null;
			socketRef.current.close();
			socketRef.current = null;
		}
		setIsConnected(false);
		setIsClosed(false);
		setPlayers([]);
		setLobbyInfos(null);
		setRound(null);
		setRoundHistory([]);
		setLastRoundResult(null);
		setGameStartAt(null);
	}

	return (
		<GameSocketContext.Provider
		value={{ players, lobbyInfos, isConnected, isClosed, isCurrentUserHost, leaveLobby, leaveGame, kickPlayer, startGame, startNextRound, round, roundHistory, lastRoundResult, sendGuess, connect, resetGame, gameStartAt}}
		>
		{children}
		</GameSocketContext.Provider>
	);
}

export function useLobbySocket() {
	return useContext(GameSocketContext);
}