import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import GamePage from './pages/game/GamePage'
import Home from './pages/HomePage' 
import ComponentsPage from './pages/ComponentsPage';
import LobbyPage from './pages/LobbyPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home/>} />
        <Route path="/game" element={<GamePage/>} />
        <Route path="/components" element={<ComponentsPage/>} />
        <Route path="/lobby" element={<LobbyPage/>} />
      </Routes>
    </Router>
  )
}

export default App;