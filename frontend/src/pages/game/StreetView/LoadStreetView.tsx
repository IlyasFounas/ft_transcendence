import { useEffect } from 'react';
import { Viewer } from 'mapillary-js';

export default function LoadStreetView() {
  useEffect(() => {
      const viewer = new Viewer({
      container: 'mapillary-viewer',
      accessToken: import.meta.env.VITE_MAPILLARY_ACCESS_TOKEN, //TODO: if this env variable is not set, docker should not start the frontend container
      imageId: '271688394642360' //TODO: get random location id for the round
    });

    return () => {
      viewer.remove();
    };
  }, []);
  return (
    <div id="mapillary-viewer" style={{ width: '100%', height: '100%' }}></div>
  );
}