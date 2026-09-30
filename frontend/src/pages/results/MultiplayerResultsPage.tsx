import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";
import { useLobbySocket } from "../../context/GameSocketContext";

import NextRoundTimer from "../../components/results/NextRoundTimer";
import ResultInfoRow from "../../components/results/ResultInfoRow";
import RankResultsList from "../../components/results/RankResultsList";
import ResultMap from "../game/2DMap/ResultMap";
import { useCurrentRoundStats } from "../../utils/useCurrentRoundStats";
import { useRoundGuesses } from "../../utils/useRoundGuesses";
import type { GameResult, PlayerRoundResult } from "../../types/game";
import PrimaryButton from "../../components/buttons/PrimaryButton";

export default function MultiplayerResultsPage() {
    const location = useLocation();
    const navigate = useNavigate();
    const { user } = useAuth();
    const { round, players, leaveLobby } = useLobbySocket();
    const [hoveredPlayerId, setHoveredPlayerId] = useState<string | null>(null);

    const { 
        currentRound, 
        totalRounds, 
        allResults, 
        actualLocation, 
        nextRoundAt, 
        isLastRound 
    } = location.state as {
        currentRound: number;
        totalRounds: number;
        allResults: PlayerRoundResult[];
        actualLocation: { latitude: number; longitude: number };
        nextRoundAt: string | null;
        isLastRound: boolean;
    };

    useEffect(() => {
        if (round && round.round_nb > currentRound) {
            navigate("/game/multi");
        }
    }, [round, currentRound, navigate]);

	const finishGame = () => {
		leaveLobby();
		navigate("/home");
	}
	
    const myResult = allResults.find(r => r.user_id === user?.ID);
    
    const formattedPlayers = allResults.map(res => {
        const playerInfo = players.find(p => p.id === res.user_id);
        return {
            player: {
                username: res.username,
                role: playerInfo?.isHost ? "Host" : "Player",
                avatar: playerInfo?.avatar ?? "/avatars/avatar1.png",
                id: res.user_id.toString()
            },
            rounds: [{
                round: res.round_nb,
                points: res.round_infos?.score ?? 0,
                distance: res.round_infos?.distance ?? 0,
                guessLocation: res.round_infos ? {
                    latitude: res.round_infos.guess.lat,
                    longitude: res.round_infos.guess.long
                } : { latitude: 0, longitude: 0 }
            }]
        };
    });

    const gameResult: GameResult = {
        currentRound,
        totalRounds,
        currentUsername: user?.Username ?? "",
        distance: myResult?.round_infos?.distance ?? 0,
        players: formattedPlayers,
    };

    const {
        animatedScore,
        scoreJustFinished,
        animatedDistance,
        distanceJustFinished,
    } = useCurrentRoundStats(gameResult);

    const guesses = useRoundGuesses(gameResult);

    return (
        <div className="grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_394px] gap-4 p-4 items-start">
            <div className="flex flex-col gap-4 order-1 lg:col-start-1">
                <ResultInfoRow
                    title="Mission Debrief"
                    infos={`Round ${currentRound} of ${totalRounds}`}
                    variant="primary"
                />

                <div className="w-full h-60 lg:h-[650px] outline outline-1 outline-sky-950 rounded-[10px] overflow-hidden">
                    <ResultMap
                        actualLocation={actualLocation}
                        guesses={guesses}
                        currentUsername={gameResult.currentUsername}
                        hoveredPlayerId={hoveredPlayerId}
                    />
                </div>

                <div className="flex flex-row gap-2">
                    <ResultInfoRow
                        title="Score"
                        infos={String(animatedScore)}
                        variant="secondary"
                        suffix="pts"
                        animate={scoreJustFinished}
                    />
                    <ResultInfoRow
                        title="Distance"
                        infos={String(animatedDistance)}
                        variant="secondary"
                        suffix="km"
                        animate={distanceJustFinished}
                    />
                </div>
            </div>

            <div className="flex flex-col items-center gap-4 order-2 lg:col-start-2">
                <ResultInfoRow title="Round Ranking" variant="primary" />
                <RankResultsList
                    players={gameResult.players}
                    currentUsername={gameResult.currentUsername}
                    onHoverPlayer={setHoveredPlayerId}
                    scrollable
                /> 
                {!isLastRound ? (
                    <NextRoundTimer targetTime={nextRoundAt}>Next round starts in</NextRoundTimer> 
                ) : (
					<PrimaryButton size="medium" onClick={finishGame} >Finish Mission</PrimaryButton>
                )}
            </div>
        </div>
    );
}