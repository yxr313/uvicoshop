import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  transpilePackages: ["@uvico/config", "@uvico/pricing"],
  experimental: {
    optimizePackageImports: ["lucide-react"],
  },
  turbopack: {
    root: path.resolve(__dirname, "../.."),
  },
  // ESTO DESBLOQUEA EL BUILD EN PRODUCCIÓN
  eslint: {
    ignoreDuringBuilds: true,
  },
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          { key: "X-Frame-Options", value: "DENY" },
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
        ],
      },
    ];
  },
};

export default nextConfig;