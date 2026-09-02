import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import GamePage from './pages/game/GamePage'
import ComponentsPage from './pages/ComponentsPage';
import LandingPage from './pages/LandingPage';
import HomePage from './pages/HomePage';
import MultiplayerPage from './pages/MultiplayerPage';
import LeaderboardPage from './pages/LeaderboardPage';
import ProfilePage from './pages/ProfilePage';
import { ChatLayout } from './layouts/ChatLayout';

function App() {
    return (
        <Router>
          <Routes>
            <Route path="/" element={<LandingPage/>} />
            <Route element={<ChatLayout />}>
                <Route path="/home" element={<HomePage />} />
                <Route path="/multiplayer" element={<MultiplayerPage />} />
            </Route>
                <Route path="/leaderboard" element={<LeaderboardPage />} />
                <Route path="/profile" element={<ProfilePage />} />
                <Route path="/game" element={<GamePage/>} />
                <Route path="/components" element={<ComponentsPage/>} />
            </Routes>
        </Router>
    )
}

export default App;