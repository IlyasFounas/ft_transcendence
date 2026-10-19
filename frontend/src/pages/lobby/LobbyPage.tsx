import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import SecondaryButton from "../../components/buttons/SecondaryButton";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import CreateOperationCard from "../../components/lobby/CreateOperationCard";
import JoinOperationCard from "../../components/lobby/JoinOperationCard";
import { useAuth } from "../../context/AuthContext";
import { checkRoomCode } from "../../services/lobby";
import { useLobbySocket } from "../../context/GameSocketContext";

export default function LobbyPage()
{
	const [mode, setMode] = useState<"create" | "join">("join");
	const [rounds, setRounds] = useState(1);
  	const [timeLimit, setTimeLimit] = useState(30);
	const navigate = useNavigate();
	const { authFetch } = useAuth();
	const [isChecking, setIsChecking] = useState(false);
	const [isShaking, setIsShaking] = useState(false);
	const [errorMessage, setErrorMessage] = useState<string | null>(null);
	const [code, setCode] = useState("");
	const { lobbyInfos, round, connect } = useLobbySocket();
	
	async function handleCreateOperation() {

		const body = new URLSearchParams();
		body.append("timer", timeLimit.toString());
		body.append("round_total", rounds.toString());

		setIsChecking(true);
			const result = await authFetch("/api/lobby/create", {
				method: "POST",
				headers: { "Content-Type": "application/x-www-form-urlencoded" },
				body: body.toString(),
				});
				setIsChecking(false);

				if (result === null) {
					setErrorMessage("Network error, please try again");
					return;
				}

				if (result.status !== 200) {
		setErrorMessage(result.data.message);
		return;
		}

		const cur = await authFetch("/api/lobby/current");
		console.log("current après create :", cur?.data);
		const gameId = cur?.data?.game_id;
		if (gameId) {
			await authFetch(`/api/game/${gameId}/mode`, {
				method: "PATCH",
				headers: { "Content-Type": "application/x-www-form-urlencoded" },
				body: new URLSearchParams({ mode: "multi" }).toString(),
			});
		}

		connect();
		navigate("/lobby");
	}

	useEffect(() => {
		if (!errorMessage) return;

		const timer = setTimeout(() => setErrorMessage(null), 4000);

		return () => clearTimeout(timer);
	}, [errorMessage]);

	function triggerInvalidShake(message: string) {
		setErrorMessage(message);
		setIsShaking(true);
		setTimeout(() => {
			setCode("");
			setIsShaking(false);
		}, 400);
	}

	async function handleJoinOperation() {
		if (code.length !== 7 || isChecking) return;

		setIsChecking(true);
		const result = await checkRoomCode(code, authFetch);
		setIsChecking(false);

		if (result.status === "ok") {
			connect();
			navigate("/lobby");
			return;
		}

		if (result.status === "full") {
			triggerInvalidShake("Room is full");
			return;
		}

		triggerInvalidShake("Invalid code");
		}

	useEffect(() => {
		if (lobbyInfos && !round) {
			navigate("/lobby");
		}
	}, [lobbyInfos, round, navigate]);

	return (
		<div className="min-h-[calc(100vh-5rem)] w-full bg-cover bg-no-repeat bg-center bg-[url(/backgrounds/earth_background.png)] py-24 xl:pr-[27rem] flex flex-col justify-center">
			<div className="flex min-h-[50px] items-center justify-center gap-4 xl:gap-12">
				<SecondaryButton
					active={mode === "create"}
					className="px-[12px] text-heading-small xl:text-heading-large"
					onClick={() => setMode("create")}
				>
					Create Operation
				</SecondaryButton>

				<SecondaryButton
					active={mode === "join"}
					className="text-heading-small xl:text-heading-large"
					onClick={() => setMode("join")}
				>
					Join Operation
				</SecondaryButton>
			</div>
			<div className="
						flex
						h-50 md:h-[300px]
						w-[calc(100%-2rem)]
						max-w-[500px]
						mx-auto m-8
						rounded-lg
						border
						border-surface-elevated
						bg-[#020B16]
						opacity-90">
					{mode === "create" ? (
						<CreateOperationCard onRoundsChange={setRounds} onTimeChange={setTimeLimit} errorMessage={errorMessage}/>
					) : (
						<JoinOperationCard onChange={(setCode)} code={code} onSubmit={handleJoinOperation} isChecking={isChecking} isShaking={isShaking} errorMessage={errorMessage}/>
					)}
			</div>
			<div className="flex justify-center">
					{mode === "create" ? (
						<PrimaryButton size="small" className="w-300" onClick={handleCreateOperation} disabled={isChecking}>Create Operation</PrimaryButton>
					) : (
						<PrimaryButton size="small" disabled={isChecking} onClick={handleJoinOperation}>Join Operation</PrimaryButton>
					)}
			</div>
	</div>
	)
}