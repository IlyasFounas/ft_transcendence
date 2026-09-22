import 'mapillary-js/dist/mapillary.css';

import { useEffect, useRef } from 'react';
import { Viewer } from 'mapillary-js';

type Props = {
    id: string;
    moving: boolean;
    zooming: boolean;
    panning: boolean;
    interactive?: boolean;
};

export default function LoadStreetView({
    id,
    moving,
    zooming,
    panning,
    interactive = true,
}: Props) {
    const containerRef = useRef<HTMLDivElement>(null);
    const viewerRef = useRef<Viewer | null>(null);

    useEffect(() => {
		const viewer = new Viewer({
			container: 'mapillary-viewer',
			accessToken: import.meta.env.VITE_MAPILLARY_ACCESS_TOKEN,

			component: {
				cover: false,

				zoom: zooming,

				pointer: {
					scrollZoom: zooming,
					touchZoom: zooming,
					doubleClickZoom: zooming,
					touchPan: panning,
					dragPan: panning,
					earthControls: panning,
				},

				sequence: moving,

				keyboard: {
					keySequenceNavigation: moving,
					keySpatialNavigation: moving,
					keyZoom: zooming,
				},

				direction: moving,
			},

			imageTiling: true,
			trackResize: true,
		});
        viewerRef.current = viewer;

        const resizeObserver = new ResizeObserver(() => {
            viewer.resize();
        });

        if (containerRef.current) {
            resizeObserver.observe(containerRef.current);
        }

        return () => {
            resizeObserver.disconnect();
            viewer.remove();
            viewerRef.current = null;
        };
    }, []);

	useEffect(() => {
		const viewer = viewerRef.current;

		if (!viewer) return;

		viewer.moveTo(id).catch((error) => {
			console.error('Failed to move to Mapillary image:', error);
		});
	}, [id]);

    // Enable / disable interaction
    useEffect(() => {
        const viewer = viewerRef.current;

        if (!viewer) return;

        const components = [
            'pointer',
            'keyboard',
            'zoom',
            'sequence',
            'direction',
        ] as const;

        components.forEach((name) => {
            if (interactive) {
                viewer.activateComponent(name);
            } else {
                viewer.deactivateComponent(name);
            }
        });
    }, [interactive]);

    return (
        <div
            ref={containerRef}
            id="mapillary-viewer"
            style={{
                width: '100%',
                height: '100%',
            }}
        />
    );
}