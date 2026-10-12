import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useLobbySocket } from "../../context/GameSocketContext";

export default function GameGuard() {
    const { round, gameStartAt, lobbyInfos, resetGame } = useLobbySocket();
    const location = useLocation();
    const navigate = useNavigate();

    useEffect(() => {
       	const isPlaying = location.pathname.includes("/game/") || location.pathname.includes("/results/") || location.pathname.includes("summary");
        const isLobby = location.pathname === "/lobby";
        const isMultiMenu = location.pathname === "/multiplayer";

        if (round && !isPlaying) {
            navigate("/game"); 
        }
        else if (gameStartAt && !isLobby && !isPlaying) {
            navigate("/lobby");
        }
        else if (lobbyInfos && isMultiMenu && !round) {
            navigate("/lobby");
        }
    }, [round, gameStartAt, lobbyInfos, location.pathname, navigate]);

	useEffect(() => {
		const onGamePage = location.pathname.startsWith("/game/");
		const onlobby = location.pathname === "/lobby";
		const missingState = (onGamePage && !round && !gameStartAt) || (onlobby && !lobbyInfos);
		if (!missingState) return;

		const timeout = setTimeout(() => {
			resetGame();
			navigate("/home");
		}, 3000);
		return () => clearTimeout(timeout);
	}, [location.pathname, round, gameStartAt, lobbyInfos, resetGame, navigate]);

    return null;
}