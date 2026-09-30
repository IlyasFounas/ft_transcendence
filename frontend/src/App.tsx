import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom'
import ComponentsPage from './pages/ComponentsPage';
import LandingPage from './pages/LandingPage';
import HomePage from './pages/HomePage';
import LeaderboardPage from './pages/LeaderboardPage';
import ProfilePage from './pages/ProfilePage';
import { ChatLayout } from './layouts/ChatLayout';
import { MainLayout } from './layouts/MainLayout';
import LobbyPage from './pages/lobby/LobbyPage';
import LobbyRoomPage from './pages/lobby/LobbyRoomPage';
import AuthenticationPage from './pages/AuthenticationPage';
import ProtectedRoute from './components/routes/ProtectedRoute';
import PublicRoute from './components/routes/PublicRoute';
import GamePlayPage from './pages/game/GamePlayPage';
import { GameSocketProvider } from './context/GameSocketContext';
import SoloResultPage from './pages/results/SoloResultsPage';
import MultiplayerResultsPage from './pages/results/MultiplayerResultsPage';

function App() {
    return (
        <Router>
            <Routes>
                <Route element={<PublicRoute />}>
                    <Route path="/" element={<LandingPage/>} />
                    <Route path="/login" element={<AuthenticationPage/>} />
                    <Route path="/components" element={<ComponentsPage/>} />
                </Route>
                <Route element={<ProtectedRoute />}>
                    <Route element={<ChatLayout />}>
                        <Route path="/home" element={<HomePage />} />
                        <Route path="/multiplayer" element={<LobbyPage />} />
                    </Route>
                    <Route element={<MainLayout />}>
                        <Route path="/leaderboard" element={<LeaderboardPage />} />
                        <Route path="/profile" element={<ProfilePage />} />
                    </Route>
					<Route element={<GameSocketProvider><Outlet /></GameSocketProvider>}>
						<Route element={<ChatLayout />}>
							<Route path="/lobbyroom" element={<LobbyRoomPage />} />
						</Route>
						<Route path="/game/:gameId" element={<GamePlayPage />} />
						<Route path="/soloresults" element={<SoloResultPage />} />
						<Route path="/multiresults" element={<MultiplayerResultsPage />} />
					</Route>
                </Route>
            </Routes>
        </Router>
    )
}

export default App;
