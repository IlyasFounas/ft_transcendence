import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import GamePage from './pages/game/GamePage'
import Home from './pages/HomePage' 
import ComponentsPage from './pages/ComponentsPage';
import LobbyPage from './pages/LobbyPage';
import GamePlayPage from './pages/game/GamePlayPage';
import SoloResultPage from './pages/game/SoloResultsPage';
import MultiplayerResultsPagev1 from './pages/game/MultiplayerResultsPageV1';
import MultiplayerResultsPagev2 from './pages/game/MultiplayerResultsPageV2';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home/>} />
        <Route path="/game" element={<GamePage/>} />
        <Route path="/components" element={<ComponentsPage/>} />
        <Route path="/lobby" element={<LobbyPage/>} />
        <Route path="/gameplay" element={<GamePlayPage/>} />
        <Route path="/soloresults" element={<SoloResultPage/>} />
		<Route path="/multiresultsv1" element={<MultiplayerResultsPagev1/>} />
		<Route path="/multiresultsv2" element={<MultiplayerResultsPagev2/>} />
      </Routes>
    </Router>
  )
}

export default App;