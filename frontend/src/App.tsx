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
import MultiplayerResultsPage from './pages/game/MultiplayerResultsPage';
import SoloResultsPage from './pages/game/SoloResultsPage';
import GameLayout from './layouts/GameLayout';
import GamePlayPage from './pages/game/GamePlayPage';
import LobbyPage from './pages/LobbyPage';

function App() {
    return (
        <Router>
          <Routes>
            <Route path="/" element={<LandingPage/>} />
            <Route element={<ChatLayout />}>
                <Route path="/home" element={<HomePage />} />
                <Route path="/multiplayer" element={<MultiplayerPage />} />
                <Route path="/profile" element={<ProfilePage />} />
            </Route>
            <Route element={<MainLayout />}>
                <Route path="/leaderboard" element={<LeaderboardPage />} />
                <Route path="/game" element={<GamePage/>} />
                <Route path="/components" element={<ComponentsPage/>} />
				<Route path="/lobby" element={< LobbyPage/>}/>
              </Route>
			<Route element={<GameLayout />}>
				<Route path="/multiresults" element={<MultiplayerResultsPage/>} />
				<Route path="/soloresults" element={<SoloResultsPage/>} />
			</Route>
			<Route path="/gameplay" element={<GamePlayPage />} />
            </Routes>
        </Router>
    )
}

export default App;