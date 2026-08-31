import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite'
// import svgr from "vite-plugin-svgr";

export default defineConfig({
  plugins: [react(), 
		tailwindcss(),
//  	 svgr({
//   svgoConfig: {
//     plugins: [
//       {
//         name: "preset-default",
//       },
//       {
//         name: "prefixIds",
//       },
//     ],
//   },
// }),
	],
  server: {
    port: 3000,
    host: true,
    watch: {
      usePolling: true,
    },
  },
});