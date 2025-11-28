/** @type {import('next').NextConfig} */
const nextConfig = {
    async rewrites() {
        return [
            {
                source: '/api/:path*',
                destination: 'http://127.0.0.1:58080/api/v1/:path*',
            },
        ]
    },
    turbopack: {
        root: process.cwd(),
    }
};

module.exports = nextConfig;
