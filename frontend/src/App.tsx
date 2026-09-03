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
import AuthenticationPage from './pages/AuthenticationPage';
import ProtectedRoute from './components/authentication/ProtectedRoute';

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
                  <Route path="/leaderboard" element={<LeaderboardPage />} />
                  <Route path="/profile" element={<ProfilePage />} />
              </Route>
            </Route>
            <Route element={<MainLayout />}>
              <Route path="/components" element={<ComponentsPage/>} />
            </Route>
          </Routes>
        </Router>
    )
}

export default App;
