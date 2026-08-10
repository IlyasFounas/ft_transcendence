import '../../styles/GamePage.css'

import LoadStreetView from './StreetView/LoadStreetView'

export default function GamePage() {
  return (
    <div className="game-container">
      <div className="game-map">
        <LoadStreetView/><h2>Mapillary Street View</h2>
      </div>
      <div className="game-sidebar">
        <h3>Game Info</h3>
      </div>
    </div>
  )
}