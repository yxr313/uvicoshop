import type { Metadata } from "next";
import { Instrument_Sans, Bricolage_Grotesque } from "next/font/google";
import { Providers } from "./providers";
import "./globals.css"; // <-- ESTA LÍNEA ES VITAL

const instrumentSans = Instrument_Sans({
  subsets: ["latin"],
  variable: "--font-instrument-sans",
  display: "swap",
});

const bricolageGrotesque = Bricolage_Grotesque({
  subsets: ["latin"],
  variable: "--font-bricolage-grotesque",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Uvico Studio | Print on Demand",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="es"
      className={`${instrumentSans.variable} ${bricolageGrotesque.variable}`}
      suppressHydrationWarning
    >
      <body className="antialiased h-[100dvh] w-screen overflow-hidden bg-neutral-50 text-foreground font-sans">
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}