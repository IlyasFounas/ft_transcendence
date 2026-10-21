import 'mapillary-js/dist/mapillary.css';

import { useEffect, useRef } from 'react';
import { Viewer, CancelMapillaryError } from 'mapillary-js';

type Props = {
    id: string;
    interactive?: boolean;
    reloadKey?: number; // changes when the server asks to load the image again (round_started received again)
    onLoadError?: () => void; // the image can't be loaded (black screen)
};

// an image still not displayed after this time counts as failed: a request can also hang without ever failing
const LOAD_TIMEOUT_MS = 15000;

export default function LoadStreetView({
    id,
    interactive = true,
    reloadKey = 0,
    onLoadError,
}: Props) {
    const containerRef = useRef<HTMLDivElement>(null);
    const viewerRef = useRef<Viewer | null>(null);
    // the latest callback, without restarting the loading each time the parent renders again
    const onLoadErrorRef = useRef(onLoadError);

    useEffect(() => {
        onLoadErrorRef.current = onLoadError;
    });

    useEffect(() => {
		const viewer = new Viewer({
			container: 'mapillary-viewer',
			accessToken: import.meta.env.MAPILLARY_ACCESS_TOKEN,

			component: {
				cover: false,

				zoom: true,

				pointer: {
					scrollZoom: true,
					touchZoom: true,
					doubleClickZoom: true,
					touchPan: true,
					dragPan: true,
					earthControls: true,
				},

				sequence: true,

				keyboard: {
					keySequenceNavigation: true,
					keySpatialNavigation: true,
					keyZoom: true,
				},

				direction: true,
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

		// becomes true once this attempt is over: loaded, failed, or abandoned (new image, reload, unmount,
		// or effects run twice by StrictMode in development). An abandoned attempt never reports anything
		let settled = false;
		const fail = (reason: unknown) => {
			if (settled) return;
			settled = true;
			console.error('Failed to load Mapillary image:', id, reason);
			onLoadErrorRef.current?.();
		};
		const timeout = setTimeout(() => fail('timeout'), LOAD_TIMEOUT_MS);

		viewer.moveTo(id)
			.then(() => {
				settled = true;
				clearTimeout(timeout);
			})
			.catch((error) => {
				clearTimeout(timeout);
				// another moveTo started before this one ended (new image, reload): it is not a failure
				if (error instanceof CancelMapillaryError) return;
				fail(error);
			});

		return () => {
			settled = true;
			clearTimeout(timeout);
		};
	}, [id, reloadKey]);

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