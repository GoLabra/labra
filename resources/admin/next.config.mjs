
/**
 * RUNTIME_ENV enables loading runtime overrides from public/env.runtime.json
 * (via the injected script) instead of relying only on build-time env.
 */
/** @type {import('next').NextConfig} */
const nextConfig = (() => {

	switch(process.env.BUILDTYPE ){
		case 'static':
			return {
				output: 'export',
				trailingSlash: true,
				images: { unoptimized: true },
				distDir: 'out.static'
			}
		case 'labra':
			return {
				output: 'export',
				trailingSlash: true,
				basePath: '/labradmin',
				assetPrefix: '/labradmin',
				images: { unoptimized: true },
				distDir: 'out.static',
				env: { 
					RUNTIME_ENV: 'true',
					BASE_PATH: '/labradmin'
				}
			}
	};

	return {}

})();

export default nextConfig;
