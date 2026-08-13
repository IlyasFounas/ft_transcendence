import 'mapillary-js/dist/mapillary.css';

import { useEffect } from 'react';
import { Viewer } from 'mapillary-js';

type Props = {
  id: string,
  moving: boolean,
  zooming: boolean,
  panning: boolean,
}

export default function LoadStreetView({id , moving, zooming, panning}: Props) {
  useEffect(() => {
      const viewer = new Viewer({
      container: 'mapillary-viewer',
      accessToken: import.meta.env.VITE_MAPILLARY_ACCESS_TOKEN, //TODO: if this env variable is not set, docker should not start the frontend container
      imageId: id,
      component: {
        cover: false, // starts directly in street view mode, without needing to click
        zoom: zooming ? true : false, // enables zooming in/out of the image
        pointer: {
          scrollZoom: zooming ? true : false, // enables zooming in/out of the image by scrolling (on desktop)
          touchZoom: zooming ? true : false, // enables zooming in/out of the image by pinching (on mobile)
          doubleClickZoom: zooming ? true : false, // enables zooming in/out of the image by double clicking (double tap on mobile)
          touchPan: panning ? true : false, // enables panning around the image by dragging (on mobile)
          dragPan: panning ? true : false, // enables panning around the image by clicking and dragging (swiping on mobile)
          earthControls: panning ? true : false, // enables panning around the image by clicking and dragging (swiping on mobile)
        },
        sequence: moving ? true : false, // enables moving to the next/previous image in the sequence
        keyboard: {
          keySequenceNavigation: moving ? true : false, // enables moving to the next/previous image in the sequence by using the arrow keys
          keySpatialNavigation: moving ? true : false, // enables moving to the next/previous image in the sequence by using the arrow keys
          keyZoom: zooming ? true : false, // enables zooming in/out of the image by using the + and - keys
        },
        direction: moving ? true : false,  // removes the directional navigation arrows
      },
      imageTiling: true, // enables image tiling when zooming (zooming in on the image will load higher resolution tiles)
      trackResize: true, //adjust the viewer size with the window size
    });

    return () => {
      viewer.remove();
    };
  }, []);
  return (
    <div id="mapillary-viewer" style={{ width: '100%', height: '100%' }}></div>
  );
}