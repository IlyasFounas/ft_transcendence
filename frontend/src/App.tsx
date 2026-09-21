import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
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
import ProtectedRoute from './components/routes/ProtectedRoute';
import PublicRoute from './components/routes/PublicRoute';

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
                        <Route path="/multiplayer" element={<MultiplayerPage />} />
                    </Route>
                    <Route element={<MainLayout />}>
                        <Route path="/leaderboard" element={<LeaderboardPage />} />
                        <Route path="/profile" element={<ProfilePage />} />
                    </Route>
                </Route>
            </Routes>
        </Router>
    )
}

export default App;
