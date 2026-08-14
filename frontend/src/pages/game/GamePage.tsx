import '../../styles/GamePage.css'

import LoadStreetView from './StreetView/LoadStreetView'
import Load2DMap from './2DMap/Load2DMap.tsx'

export default function GamePage() {
  return (
    <div className="game-container">
      <div className="game-map">
        <h3>Mapillary Street View</h3>
        <LoadStreetView id="271688394642360" moving={false} zooming={true} panning={true} />
      </div>
      <div className="game-sidebar">
        <h3>Game Info</h3>
        <Load2DMap latitude={45.782429} longitude={4.745860} />
      </div>
    </div>
  )
}