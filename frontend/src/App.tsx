import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import GamePage from './pages/game/GamePage'
import ComponentsPage from './pages/ComponentsPage';
import LandingPage from './pages/LandingPage';
import HomePage from './pages/HomePage';
import MultiplayerPage from './pages/MultiplayerPage';
import LeaderboardPage from './pages/LeaderboardPage';
import ProfilePage from './pages/ProfilePage';
import { ChatLayout } from './layouts/ChatLayout';
import { MainLayout } from './layouts/MainLayout';
import LobbyPage from './pages/lobby/LobbyPage';
import LobbyRoomPage from './pages/lobby/LobbyRoomPage';
import AuthenticationPage from './pages/AuthenticationPage';
import ProtectedRoute from './components/authentication/ProtectedRoute';
import MultiplayerResultsPage from './pages/results/MultiplayerResultsPage';
import SoloResultsPage from './pages/results/SoloResultsPage';
import GameLayout from './layouts/GameLayout';
import GamePlayPage from './pages/game/GamePlayPage';
import SummaryResultsPage from './pages/results/SummaryResultsPage';

function App() {
    return (
        <Router>
          <Routes>

            <Route path="/" element={<LandingPage/>} />
            <Route path="/login" element={<AuthenticationPage/>} />
            <Route path="/game" element={<GamePage/>} />

            <Route element={<ProtectedRoute />}>
              <Route element={<ChatLayout />}>
                  <Route path="/home" element={<HomePage />} />
                  <Route path="/multiplayer" element={<MultiplayerPage />} />
                  <Route path="/profile" element={<ProfilePage />} />
                  <Route path="/lobby" element={<LobbyPage/>} />
                  <Route path="/lobbyroom" element={<LobbyRoomPage/>} />
              </Route>
			  <Route element={<GameLayout />}>
                <Route path="/soloresults" element={<SoloResultsPage/>} />
                <Route path="/multiresults" element={<MultiplayerResultsPage/>} />
                <Route path="/summaryresults" element={<SummaryResultsPage/>} />
			  </Route>
            </Route>
            <Route element={<MainLayout />}>
                <Route path="/game" element={<GamePage/>} />
                <Route path="/components" element={<ComponentsPage/>} />
				<Route path="/leaderboard" element={<LeaderboardPage/>} />
              </Route>
            </Routes>
        </Router>
    )
}

export default App;
