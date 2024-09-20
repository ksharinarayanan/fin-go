/** @type {import('next').NextConfig} */
const nextConfig = {
	async rewrites() {
		return [
			{
				source: "/backend/:path*",
				destination: "http://localhost:8080/:path*", // Proxy to backend
			},
		];
	},
};

export default nextConfig;
