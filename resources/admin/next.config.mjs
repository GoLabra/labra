const forLabra = process.env.LABRA_BUILD === 'true';


/** @type {import('next').NextConfig} */
const nextConfig = {
  ...(forLabra && {
    output: 'export',
    trailingSlash: true,
    basePath: '/labradmin',
    assetPrefix: '/labradmin',
    images: { unoptimized: true },
	distDir: 'out.labradmin',
  }),
};

export default nextConfig;
