import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import GamePage from './pages/game/GamePage'
import HomePage from './pages/HomePage' 
import ComponentsPage from './pages/ComponentsPage';
import LandingPage from './pages/LandingPage';

import { useAuth } from "./context/AuthContext";

function App() {
  const { isAuthenticated } = useAuth();
  
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage/>} />
        <Route path="/home" element={<HomePage />} />
        <Route path="/game" element={<GamePage/>} />
        <Route path="/components" element={<ComponentsPage/>} />
      </Routes>
    </Router>
  )
}

export default App;