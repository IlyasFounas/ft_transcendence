import { BrowserRouter as Router, Routes, Route} from 'react-router-dom'
import { ChatLayout } from './layouts/ChatLayout';
import { MainLayout } from './layouts/MainLayout';
import { GameSessionLayout } from './layouts/GameSessionLayout';
import { GameSocketProvider } from './context/GameSocketContext';
import LandingPage from './pages/LandingPage';
import HomePage from './pages/HomePage';
import ProfilePage from './pages/ProfilePage';
import LobbyPage from './pages/lobby/LobbyPage';
import LobbyRoomPage from './pages/lobby/LobbyRoomPage';
import AuthenticationPage from './pages/AuthenticationPage';
import ProtectedRoute from './components/routes/ProtectedRoute';
import PublicRoute from './components/routes/PublicRoute';
import NotFoundPage from './pages/errors/404';
import GamePlayPage from './pages/game/GamePlayPage';
import SoloResultPage from './pages/results/SoloResultsPage';
import MultiplayerResultsPage from './pages/results/MultiplayerResultsPage';
import SummaryResultsPage from './pages/results/SummaryResultsPage';
import LeaderboardPage from './pages/LeaderboardPage';

function App() {
	return (
		<Router>
			<GameSocketProvider>
				<Routes>
					<Route element={<PublicRoute />}>
						<Route path="/" element={<LandingPage/>} />
                    	<Route path="*" element={<NotFoundPage />} />
						<Route path="/login" element={<AuthenticationPage/>} />
					</Route>
					<Route element={<ProtectedRoute />}>
						<Route element={<GameSessionLayout />}>
							<Route element={<ChatLayout />}>
								<Route path="/home" element={<HomePage />} />
								<Route path="/multiplayer" element={<LobbyPage />} />
								<Route path="/lobby" element={<LobbyRoomPage />} />
							</Route>
							<Route element={<MainLayout />}>
								<Route path="/profile" element={<ProfilePage />} />
								<Route path="/leaderboard" element={<LeaderboardPage />} />
							</Route>
							<Route path="/game/:gameId" element={<GamePlayPage />} />
							<Route path="/results/solo" element={<SoloResultPage />} />
							<Route path="/results/multiplayers" element={<MultiplayerResultsPage />} />
							<Route path="/results/summary" element={<SummaryResultsPage />} />
						</Route>
					</Route>
				</Routes>
			</GameSocketProvider>
		</Router>
	)
}

export default App;