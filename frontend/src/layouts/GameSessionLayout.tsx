import { Outlet } from "react-router-dom";
import GameGuard from "../components/routes/GameGuard";

export const GameSessionLayout = () => {
	return (
		<>
			<GameGuard />
			<Outlet />
		</>
	);
};